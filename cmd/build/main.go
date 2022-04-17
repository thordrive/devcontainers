package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"thordrive.ai/devcontainers/pkg/spec"
)

func printUsage() {
	fmt.Printf(`Build docker image with given reference.
Usage: %s REFERENCE
  REFERENCE
      Reference of the image to build
`, os.Args[0])
}

type BuildContext struct {
	path     string
	manifest spec.Manifest
	image    spec.Image
}

type Args struct {
	Reference string
	DryRun    bool
}

func main() {
	args := Args{
		DryRun: false,
	}

	flag.BoolVar(&args.DryRun, "dry-run", false, "stop before build")
	flag.Parse()

	if flag.NArg() != 1 {
		printUsage()
		os.Exit(1)
	} else {
		args.Reference = flag.Arg(0)
	}

	docker_binary, err := exec.LookPath("docker")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Fatal("docker not found")
		} else {
			log.Fatal(err)
		}
	}

	files, err := ioutil.ReadDir("containers")
	if err != nil {
		log.Fatal(err)
	}

	build_ctx := BuildContext{
		path: "",
	}

	if err := spec.Walk(files, func(manifest_path string, manifest *spec.Manifest) error {
		for _, image := range manifest.Images {
			for _, ref := range manifest.RefsOf(image) {
				if args.Reference != ref {
					continue
				}

				build_ctx = BuildContext{
					path:     filepath.Dir(manifest_path),
					manifest: *manifest,
					image:    image,
				}

				return io.EOF
			}
		}

		return nil
	}); err != nil {
		if !errors.Is(err, io.EOF) {
			log.Fatal(err)
		}
	}

	if len(build_ctx.path) == 0 {
		log.Fatalln("failed to find reference", args.Reference)
	}

	// fmt.Printf("build_ctx: %v\n", build_ctx)

	cmd := &exec.Cmd{
		Path:   docker_binary,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	cmd.Args = func() []string {
		args := []string{cmd.Path, "build"}
		args = append(args, "--file", build_ctx.manifest.Dockerfile)

		for _, tag := range build_ctx.image.Tags {
			args = append(args, "--tag", fmt.Sprintf("%s:%s", build_ctx.manifest.Name, tag))
		}

		build_args := make(map[string]string)
		for k, v := range build_ctx.manifest.Args {
			build_args[k] = v
		}
		for k, v := range build_ctx.image.Args {
			build_args[k] = v
		}

		build_args["FROM"] = build_ctx.image.From

		for k, v := range build_args {
			args = append(args, "--build-arg", fmt.Sprintf("%s=%s", k, v))
		}

		args = append(args, build_ctx.path)

		return args
	}()

	fmt.Printf("run: %v\n", cmd.Args)

	if args.DryRun {
		return
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
