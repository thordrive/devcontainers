package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
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
	host string
	keys map[string]Auth
}

func NewAuthStore() *AuthStore {
	return &AuthStore{
		host: "auth.docker.io",
		keys: make(map[string]Auth),
	}
}

func (s *AuthStore) Get(name string) (string, error) {
	auth, ok := s.keys[name]
	if ok && !auth.IsExpired() {
		return auth.Token, nil
	}

	res, err := http.Get(fmt.Sprintf("https://%s/token?service=registry.docker.io&scope=repository:%s:pull", s.host, name))
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
