package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"thordrive.ai/devcontainers/pkg/registry"
	"thordrive.ai/devcontainers/pkg/spec"
)

type Args struct {
	Verbose bool
}

func main() {
	args := Args{}

	flag.BoolVar(&args.Verbose, "v", false, "print logs")
	flag.Parse()

	files, err := ioutil.ReadDir("containers")
	if err != nil {
		log.Fatal(err)
	}

	var build_tree spec.BuildTree
	if err := spec.Walk(files, spec.ResolveBuildTree(&build_tree)); err != nil {
		log.Fatal(err)
	}

	get_date := func(ref string) (time.Time, error) {
		img_manifest, err := registry.DefaultClient.GetImageManifest(ref)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to get image manifest: %w", err)
		}

		if len(img_manifest.History) == 0 {
			return time.Time{}, fmt.Errorf("history is empty: %v", img_manifest)
		}

		history_v1 := &registry.ImageManifestHistoryV1{}
		if err := json.Unmarshal([]byte(img_manifest.History[0].V1Compatibility), history_v1); err != nil {
			return time.Time{}, fmt.Errorf("failed to parse v1 history: %v", img_manifest)
		}

		return history_v1.Created, nil
	}

	root_entries := build_tree.RootEntries()
	for _, root_entry := range root_entries {
		if args.Verbose {
			log.Printf("iterate root %s\n", root_entry.Ref)
		}

		for _, child_entry := range root_entry.Childs {
			if !child_entry.IsOrigin() {
				continue
			}

			if args.Verbose {
				log.Printf("iterate child %s\n", child_entry.Ref)
			}

			root_entry_date, err := get_date(root_entry.Ref)
			if err != nil {
				log.Fatalf("failed to get date for %s: %s", root_entry.Ref, err)
			}

			if args.Verbose {
				log.Printf("date %s %s", root_entry_date, root_entry.Ref)
			}

			child_entry_date, err := get_date(child_entry.Ref)
			if err != nil {
				if !(errors.Is(err, registry.ErrNotFound) || errors.Is(err, registry.ErrUnauthorized)) {
					log.Fatalf("failed to get date for %s: %s", child_entry.Ref, err)
				}

				child_entry_date = time.Time{}
			}

			if args.Verbose {
				log.Printf("date %s %s", child_entry_date, child_entry.Ref)
			}

			if root_entry_date.Before(child_entry_date) {
				continue
			}

			fmt.Println(child_entry.Ref)
		}
	}
}
