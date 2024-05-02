package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"io"
	"net/http"
	"time"
)

const (
	DefaultHost = "https://www.hackmud.com"
)

type chatApi struct {
	client   *http.Client
	viewport viewport.Model
	pass     string
	token    *string
}

type ChatApi interface {
	FetchChatsForUser(user string, opts ...GetChatsRequestOption) ([]Chat, error)
}

func New(c *http.Client, chatPass string) ChatApi {
	return &chatApi{
		client: c,
		pass:   chatPass,
	}
}

func (a *chatApi) fetchToken() (string, error) {
	body, _ := json.Marshal(NewGetTokenRequest(a.pass))
	req, _ := http.NewRequest(http.MethodPost, DefaultHost+GetTokenEndpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting token: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	r := GetTokenResponse{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", fmt.Errorf("error umarshalling token response: %w", err)
	}

	if r.Ok != true {
		return "", errors.New(fmt.Sprintf("get token response was not ok: %s", string(data)))
	}
	return r.ChatToken, nil
}

func (a *chatApi) FetchChatsForUser(user string, opts ...GetChatsRequestOption) ([]Chat, error) {
	if len(opts) != 1 {
		return nil, errors.New("must specify before OR after")
	}
	var token string
	if a.token == nil {
		var err error
		token, err = a.fetchToken()
		if err != nil {
			return nil, fmt.Errorf("error fetching token for chats: %w", err)
		}
		a.token = &token
	}

	body, _ := json.Marshal(newGetChatsForUserRequest(*a.token, user, opts...))
	resp, err := a.client.Post(DefaultHost+GetChatsEndpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error fetching chats: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading ")
	}

	r := GetChatsResponse{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, fmt.Errorf("error umarshalling chats response: %w", err)
	}

	if r.Chats != nil {
		c, ok := r.Chats[user]
		if !ok {
			return nil, nil
		}
		return c, nil
	}

	return nil, nil
}

func toRubyTime(t time.Time) int64 {
	return t.Unix()
}
