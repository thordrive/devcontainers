package spec

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"thordrive.ai/devcontainers/pkg/registry"
)

var (
	slotPattern = regexp.MustCompile(`\{(\d+)\}`)
)

type BuildTreeNode struct {
	Parent *BuildTreeNode
	Childs map[string]*BuildTreeNode

	Ref     string
	AliasOf string
}

func makeNode(ref string) *BuildTreeNode {
	return &BuildTreeNode{
		Childs: make(map[string]*BuildTreeNode),
		Ref:    ref,
	}
}

func (n *BuildTreeNode) addChild(child *BuildTreeNode) {
	child.Parent = n
	n.Childs[child.Ref] = child
}

func (n *BuildTreeNode) Walk(fn func(depth int, n *BuildTreeNode) bool) {
	n.walk(0, fn)
}

func (n *BuildTreeNode) IsOrigin() bool {
	return len(n.AliasOf) == 0
}

func (n *BuildTreeNode) walk(depth int, fn func(depth int, n *BuildTreeNode) bool) {
	if !fn(depth, n) {
		return
	}

	for _, child := range n.Childs {
		child.walk(depth+1, fn)
	}
}

type BuildTree struct {
	Entries map[string]*BuildTreeNode
}

func (bt *BuildTree) Pack() {
	for _, entry := range bt.Entries {
		if entry.IsOrigin() {
			continue
		}

		origin, ok := bt.Entries[entry.AliasOf]
		if !ok {
			panic("origin reference does not exists")
		}

		for k, v := range entry.Childs {
			origin.Childs[k] = v
		}
	}
}

func (bt *BuildTree) RootEntries() []*BuildTreeNode {
	var rst []*BuildTreeNode
	for _, entry := range bt.Entries {
		if entry.Parent == nil {
			rst = append(rst, entry)
		}
	}

	return rst
}

func ExpandBySemver(image Image, pattern *regexp.Regexp) ([]Image, error) {
	name := strings.SplitN(image.From, ":", 2)[0]

	remote_tags, err := registry.DefaultClient.GetTags(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag from registry: %w", err)
	}

	expanded_images := []Image{}

	// Expanding.
	for _, remote_tag := range remote_tags {
		match := pattern.FindStringSubmatch(remote_tag)
		if len(match) == 0 {
			continue
		}

		expanded_image := image
		expanded_image.From = name + ":" + match[0]
		expanded_image.Tags = append([]string{}, image.Tags...)
		for i, template := range image.Tags {
			var replace_err error = nil
			expanded_image.Tags[i] = slotPattern.ReplaceAllStringFunc(template, func(slot string) string {
				index, err := strconv.ParseUint(slot[1:len(slot)-1], 10, 64)
				if err != nil {
					replace_err = fmt.Errorf("failed to parse uint from %s: %w", template, err)
					return ""
				}

				// `index` is 0 based and the submatch item is 1 based since `match[0]` is the matched string.
				if index+1 >= uint64(len(match)) {
					replace_err = errors.New("index out of range: " + template)
					return ""
				}

				return match[index+1]
			})

			if replace_err != nil {
				return nil, replace_err
			}
		}

		expanded_images = append(expanded_images, expanded_image)
	}

	get_matched_ints := func(tag string) []uint64 {
		match := pattern.FindStringSubmatch(tag)
		if len(match) < 2 {
			return []uint64{}
		}

		rst := make([]uint64, 0, len(match)-1)
		for _, word := range match[1:] {
			d, _ := strconv.ParseUint(word, 10, 64) // It will not fail we already checked that while expanding.
			rst = append(rst, d)
		}

		return rst
	}

	// Merge duplicate tags.
	for _, expanded_image_lhs := range expanded_images {
		match_lhs := get_matched_ints(expanded_image_lhs.From[len(name)+1:])
		for index_rhs, expanded_image_rhs := range expanded_images {
			if expanded_image_lhs.From == expanded_image_rhs.From {
				// Self
				continue
			}

			match_rhs := get_matched_ints(expanded_image_rhs.From[len(name)+1:])

			if len(match_lhs) != len(match_rhs) {
				// Is this can be resolved?
				return nil, errors.New("maybe alias is matched?")
			}

			take_lhs := false
			for i := range match_lhs {
				if match_lhs[i] == match_rhs[i] {
					continue
				}

				take_lhs = match_lhs[i] > match_rhs[i]
				break
			}

			if !take_lhs {
				continue
			}

			for _, tag_lhs := range expanded_image_lhs.Tags {
				index_new := 0
				for _, tag_rhs := range expanded_image_rhs.Tags {
					if tag_lhs == tag_rhs {
						// Abandon duplicated one on the rhs.
						continue
					}

					expanded_image_rhs.Tags[index_new] = tag_rhs
					index_new++
				}

				expanded_images[index_rhs].Tags = expanded_image_rhs.Tags[:index_new]
			}
		}
	}

	return expanded_images, nil
}

func ExpandImage(image Image) ([]Image, error) {
	name, tag := func() (string, string) {
		ref := strings.SplitN(image.From, ":", 2)
		return ref[0], ref[1]
	}()

	if len(name) == 0 {
		return nil, errors.New(`"from" must have a name`)
	}

	if len(tag) == 0 {
		return nil, errors.New(`"from" must have a tag`)
	}

	if !((len(tag) > 2) && (tag[0] == '/') && (tag[len(tag)-1] == '/')) {
		return []Image{image}, nil
	}

	pattern, err := regexp.Compile(tag[1 : len(tag)-2])
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	return ExpandBySemver(image, pattern)
}

func ResolveBuildTree(build_tree *BuildTree) Walker {
	build_tree.Entries = make(map[string]*BuildTreeNode)
	entries := build_tree.Entries
	return func(manifest_path string, manifest *Manifest) error {
		for _, image := range manifest.Images {
			expanded_images, err := ExpandImage(image)
			if err != nil {
				return fmt.Errorf("failed to expand image (from: %s): %w", image.From, err)
			}

			for _, img := range expanded_images {
				parent, ok := entries[img.From]
				if !ok {
					// Parent not evaluated yet.
					parent = makeNode(img.From)
					entries[img.From] = parent
				}

				for _, ref := range manifest.RefsOf(img) {
					entry, ok := entries[ref]
					if !ok {
						entry = makeNode(ref)
						if origin_ref := manifest.RefOf(img); origin_ref != ref {
							entry.AliasOf = origin_ref
						}

						entries[ref] = entry
					}

					parent.addChild(entry)
				}
			}
		}

		return nil
	}
}
