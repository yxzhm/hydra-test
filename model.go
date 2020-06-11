package main

type Config struct {
	TokenUrl     string `json:"token_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	ThreadNum    int    `json:"thread_num"`
	LoopNum      int    `json:"loop_num"`
}
