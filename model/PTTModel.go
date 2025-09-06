package model

type PttVoice struct {
	ID       string `json:"id"`
	Channel  string `json:"channel"`
	FileName string `json:"fileName"`
}
type PttUser struct {
	ID      string `json:"id"`
	Token   string `json:"token,omitempty"`
	Channel string `json:"channel,omitempty"`
}
