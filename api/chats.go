package api

import "time"

const (
	GetChatsEndpoint = "/mobile/chats.json"
)

type GetChatsRequest struct {
	Token     string   `json:"chat_token"`
	Usernames []string `json:"usernames"`
	Before    int64    `json:"before,omitempty"`
	After     int64    `json:"after,omitempty"`
}

type GetChatsRequestOption func(request *GetChatsRequest)

func WithBefore(before time.Time) GetChatsRequestOption {
	return func(request *GetChatsRequest) {
		request.Before = toRubyTime(before)
	}
}

func WithAfter(after time.Time) GetChatsRequestOption {
	return func(request *GetChatsRequest) {
		request.After = toRubyTime(after)
	}
}

func newGetChatsForUserRequest(token string, user string, opts ...GetChatsRequestOption) GetChatsRequest {
	unames := make([]string, 1)
	unames[0] = user
	r := GetChatsRequest{
		Token:     token,
		Usernames: unames,
	}

	for _, o := range opts {
		o(&r)
	}

	return r
}

type GetChatsResponse struct {
	Ok    bool              `json:"ok"`
	Chats map[string][]Chat `json:"chats"`
}

type Chat struct {
	Id      string  `json:"id"`
	T       float64 `json:"t"`
	User    string  `json:"from_user"`
	Msg     string  `json:"msg"`
	IsJoin  bool    `json:"is_join,omitempty"`
	IsLeave bool    `json:"is_leave,omitempty"`
	Channel string  `json:"channel"`
}
