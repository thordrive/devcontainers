package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"thordrive.ai/devcontainers/pkg/spec"
)

type node struct {
	parent *node
	childs map[string]*node

	Ref     string
	AliasOf string
}

func makeNode(ref string) *node {
	return &node{
		childs: make(map[string]*node),
		Ref:    ref,
	}
}

func (n *node) AddChild(child *node) {
	child.parent = n
	n.childs[child.Ref] = child
}

func (n *node) Walk(fn func(depth int, n *node) bool) {
	n.walk(0, fn)
}

func (n *node) IsOrigin() bool {
	return len(n.AliasOf) == 0
}

func (n *node) walk(depth int, fn func(depth int, n *node) bool) {
	if !fn(depth, n) {
		return
	}

	for _, child := range n.childs {
		child.walk(depth+1, fn)
	}
}

type Args struct {
	Root  string
	Strip int
	Depth int
	Pack  bool
}

func main() {
	var args Args

	flag.IntVar(&args.Strip, "strip", 0, "start print from n-th child")
	flag.IntVar(&args.Depth, "depth", 0, "graph depth")
	flag.BoolVar(&args.Pack, "pack", false, "merge sub-branches due to aliases")
	flag.Parse()

	if flag.NArg() > 0 {
		args.Root = flag.Arg(0)
	}

	entries := make(map[string]*node)

	files, err := ioutil.ReadDir("containers")
	if err != nil {
		log.Fatal(err)
	}

	if err := spec.Walk(files, func(_ string, manifest *spec.Manifest) bool {
		for _, image := range manifest.Images {
			parent, ok := entries[image.From]
			if !ok {
				// Parent not evaluated yet.
				parent = makeNode(image.From)
				entries[image.From] = parent
			}

			for _, ref := range manifest.RefsOf(image) {
				entry, ok := entries[ref]
				if !ok {
					entry = makeNode(ref)
					if origin_ref := manifest.RefOf(image); origin_ref != ref {
						entry.AliasOf = origin_ref
					}

					entries[ref] = entry
				}

				parent.AddChild(entry)
			}
		}

		return true
	}); err != nil {
		log.Fatal(err)
	}

	walker := func(level int, entry *node) bool {
		if args.Pack && !entry.IsOrigin() {
			return false
		}

		if level < args.Strip {
			return true
		}

		fmt.Printf("%s%s\n", strings.Repeat("  ", level-args.Strip), entry.Ref)

		return level-args.Strip+1 != args.Depth
	}

	if args.Pack {
		for _, entry := range entries {
			if entry.IsOrigin() {
				continue
			}

			origin, ok := entries[entry.AliasOf]
			if !ok {
				panic("origin reference does not exists")
			}

			for k, v := range entry.childs {
				origin.childs[k] = v
			}
		}
	}

	rootEntries := make([]*node, 0)
	if len(args.Root) == 0 {
		for _, entry := range entries {
			if entry.parent == nil {
				rootEntries = append(rootEntries, entry)
			}
		}
	} else {
		origin_entry, ok := entries[args.Root]
		if !ok {
			log.Fatalln("reference does not exists:", args.Root)
		}

		if !origin_entry.IsOrigin() {
			origin_entry, ok = entries[origin_entry.AliasOf]
			if !ok {
				panic("origin reference does not exists")
			}
		}

		rootEntries = append(rootEntries, origin_entry)
		for _, entry := range entries {
			if entry.AliasOf != origin_entry.Ref {
				continue
			}

			rootEntries = append(rootEntries, entry)
		}
	}
	for _, entry := range rootEntries {
		entry.Walk(walker)
	}
}
