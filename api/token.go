package api

const (
	GetTokenEndpoint = "/mobile/get_token.json"
)

type GetTokenRequest struct {
	Pass string `json:"pass"`
}

func NewGetTokenRequest(pass string) GetTokenRequest {
	return GetTokenRequest{
		Pass: pass,
	}
}

type GetTokenResponse struct {
	Ok        bool   `json:"ok"`
	ChatToken string `json:"chat_token"`
}
