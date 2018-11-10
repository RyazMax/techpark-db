package models

type Post struct {
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	Id       int    `json:"id"`
	IsEdited bool   `json:"isEdited"`
	Message  string `json:"message"`
	Parent   int    `json:"parent"`
	Thread   int    `json:"thread"`
}

type PostFull struct {
	Author User
	Forum  Forum
	Post   Post
	Thread Thread
}

type PostUpdate struct {
	Message string `json:"message"`
}

type Posts []Post
