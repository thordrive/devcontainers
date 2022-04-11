package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"thordrive.ai/devcontainers/pkg/spec"
)

type Auth struct {
	Token     string    `json:"token"`
	ExpiresIn uint      `json:"expires_in"`
	IssuedAt  time.Time `json:"issued_at"`
}

func (a *Auth) IsExpired() bool {
	// TODO
	return false
}

type AuthStore struct {
	keys map[string]Auth
}

func NewAuthStore() *AuthStore {
	return &AuthStore{
		keys: make(map[string]Auth),
	}
}

func (s *AuthStore) Get(name string) (string, error) {
	auth, ok := s.keys[name]
	if ok && !auth.IsExpired() {
		return auth.Token, nil
	}

	res, err := http.Get(fmt.Sprintf("https://%s/token?service=registry.docker.io&scope=repository:%s:pull", "auth.docker.io", name))
	if err != nil {
		return "", fmt.Errorf("failed to request: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	if err := json.Unmarshal(body, &auth); err != nil {
		return "", fmt.Errorf("failed to parse auth: %w", err)
	}

	s.keys[name] = auth
	return auth.Token, nil
}

type ImageManifestHistoryV1 struct {
	Created time.Time `json:"created"`
}

type ImageManifestHistroy struct {
	V1Compatibility string `json:"v1Compatibility"`
}

type ImageManifest struct {
	History []ImageManifestHistroy `json:"history"`
}

type RegistryClient struct {
	auth_store *AuthStore
}

func NewRegistryClient() *RegistryClient {
	return &RegistryClient{
		auth_store: NewAuthStore(),
	}
}

func (c *RegistryClient) GetImageManifest(ref string) (*ImageManifest, error) {
	name_tag := strings.SplitN(ref, ":", 2)

	token, err := c.auth_store.Get(name_tag[0])
	if err != nil {
		log.Fatalf("failed to get token for %s: %s", ref, err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", name_tag[0], name_tag[1]), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to new request: %w", err)
	}

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("%s: %s", res.Status, string(body))
	}

	image_manifest := &ImageManifest{}
	if err := json.Unmarshal(body, &image_manifest); err != nil {
		return nil, fmt.Errorf("failed to parse auth: %w", err)
	}

	return image_manifest, nil
}

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

	registry_client := NewRegistryClient()
	get_date := func(ref string) (time.Time, error) {
		img_manifest, err := registry_client.GetImageManifest(ref)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to get image manifest: %w", err)
		}

		if len(img_manifest.History) == 0 {
			return time.Time{}, fmt.Errorf("history is empty: %v", img_manifest)
		}

		history_v1 := &ImageManifestHistoryV1{}
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
				log.Fatalf("failed to get date for %s: %s", child_entry.Ref, err)
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
