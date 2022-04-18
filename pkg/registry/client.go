package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	DefaultClient *Client = NewClient()
)

type Client struct {
	host       string
	auth_store *AuthStore
}

func NewClient() *Client {
	return &Client{
		host:       "registry-1.docker.io",
		auth_store: NewAuthStore(),
	}
}

func (c *Client) newRequest(name string, resource string) (*http.Request, error) {
	token, err := c.auth_store.Get(name)
	if err != nil {
		log.Fatalf("failed to get token: %s", err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/v2/%s/%s", c.host, name, resource), nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
	}

	return req, nil
}

func (c *Client) request(name string, resource string) (*http.Response, error) {
	req, err := c.newRequest(name, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to new request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)

	if res.StatusCode == http.StatusUnauthorized {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
		return nil, fmt.Errorf("%w: %s", ErrUnauthorized, string(body))
	}

	return res, err
}

func (c *Client) GetImageManifest(ref string) (*ImageManifest, error) {
	name_tag := strings.SplitN(ref, ":", 2)

	res, err := c.request(name_tag[0], "manifests/"+name_tag[1])
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	image_manifest := &ImageManifest{}
	if err := json.Unmarshal(body, &image_manifest); err != nil {
		return nil, fmt.Errorf("failed to parse image manifest: %w", err)
	}

	return image_manifest, nil
}

func (c *Client) GetTags(name string) ([]string, error) {
	res, err := c.request(name, "tags/list")
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	tag_list := &TagList{}
	if err := json.Unmarshal(body, &tag_list); err != nil {
		return nil, fmt.Errorf("failed to parse tag list: %w", err)
	}

	return tag_list.Tags, nil
}
