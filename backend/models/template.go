package models

type Template struct {
	Id      string `json:"id" bson:"_id,omitempty"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
