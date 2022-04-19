package spec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	ManifestFilename string = "manifest.yml"
)

type Image struct {
	Tags       []string
	From       string
	Dockerfile string            `yaml:",omitempty"`
	Args       map[string]string `yaml:",omitempty"`
}

type Manifest struct {
	Name       string
	Dockerfile string            `yaml:",omitempty"`
	Args       map[string]string `yaml:",omitempty"`
	Images     []Image
}

func ReadManifest(manifest_path string, manifest *Manifest) error {
	data, err := os.ReadFile(manifest_path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := yaml.Unmarshal(data, manifest); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	if len(manifest.Dockerfile) == 0 {
		manifest.Dockerfile = filepath.Join(filepath.Dir(manifest_path), "Dockerfile")
	}

	return nil
}

func (m Manifest) RefOf(image Image) string {
	return m.Name + ":" + image.Tags[0]
}

func (m Manifest) RefsOf(image Image) []string {
	refs := make([]string, len(image.Tags))
	for i, tag := range image.Tags {
		refs[i] = m.Name + ":" + tag
	}

	return refs
}

type Walker func(manifest_path string, manifest *Manifest) error

func Walk(files []fs.FileInfo, fn Walker) error {
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		context_path, err := filepath.Abs(filepath.Join("containers", file.Name()))
		if err != nil {
			return err
		}

		manifest_path := filepath.Join(context_path, ManifestFilename)

		if err := (func() error {
			if _, err := os.Stat(manifest_path); errors.Is(err, os.ErrNotExist) {
				return nil
			} else if err != nil {
				return fmt.Errorf("failed to stat manifest: %w", err)
			}

			var manifest Manifest
			if err := ReadManifest(manifest_path, &manifest); err != nil {
				return fmt.Errorf("failed to read manifest: %w", err)
			}

			if err := fn(manifest_path, &manifest); err != nil {
				return fmt.Errorf("failed to walk: %w", err)
			}

			return nil
		})(); err != nil {
			return fmt.Errorf("error at %s: %w", manifest_path, err)
		}
	}

	return nil
}
