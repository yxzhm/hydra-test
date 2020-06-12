package main

type Config struct {
	TokenUrl      string `json:"token_url"`
	IntrospectUrl string `json:"introspect_url"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	Scope         string `json:"scope"`
	ThreadNum     int    `json:"thread_num"`
	LoopNum       int    `json:"loop_num"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
