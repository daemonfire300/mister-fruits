package model

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Token struct {
	Value string `json:"value,omitempty"`
}
