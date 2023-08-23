package models

type Post struct {
	Id string `json:"id"`
	PostContent string `json:"post_content"`
	UserId string `json:"user_id"`
}