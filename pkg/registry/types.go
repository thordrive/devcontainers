package registry

import "time"

type ImageManifestHistoryV1 struct {
	Created time.Time `json:"created"`
}

type ImageManifestHistroy struct {
	V1Compatibility string `json:"v1Compatibility"`
}

type ImageManifest struct {
	History []ImageManifestHistroy `json:"history"`
}

type TagList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
