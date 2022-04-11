package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"thordrive.ai/devcontainers/pkg/spec"
)

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

	files, err := ioutil.ReadDir("containers")
	if err != nil {
		log.Fatal(err)
	}

	var build_tree spec.BuildTree

	if err := spec.Walk(files, spec.ResolveBuildTree(&build_tree)); err != nil {
		log.Fatal(err)
	}

	walker := func(level int, entry *spec.BuildTreeNode) bool {
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
		build_tree.Pack()
	}

	root_entries := make([]*spec.BuildTreeNode, 0)
	if len(args.Root) == 0 {
		root_entries = build_tree.RootEntries()
	} else {
		origin_entry, ok := build_tree.Entries[args.Root]
		if !ok {
			log.Fatalln("reference does not exists:", args.Root)
		}

		if !origin_entry.IsOrigin() {
			origin_entry, ok = build_tree.Entries[origin_entry.AliasOf]
			if !ok {
				panic("origin reference does not exists")
			}
		}

		root_entries = append(root_entries, origin_entry)
		for _, entry := range build_tree.Entries {
			if entry.AliasOf != origin_entry.Ref {
				continue
			}

			root_entries = append(root_entries, entry)
		}
	}
	for _, entry := range root_entries {
		entry.Walk(walker)
	}
}
