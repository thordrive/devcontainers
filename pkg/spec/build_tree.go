package spec

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

func ResolveBuildTree(build_tree *BuildTree) Walker {
	build_tree.Entries = make(map[string]*BuildTreeNode)
	entries := build_tree.Entries
	return func(manifest_path string, manifest *Manifest) bool {
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

				parent.addChild(entry)
			}
		}

		return true
	}
}
