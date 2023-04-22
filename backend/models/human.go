package models

type Human struct {
	Id       string   `json:"id" bson:"_id,omitempty"`
	Name     string   `json:"name"`
	NickName string   `json:"nickName"`
	ChatIds  []string `json:"chatIds"`
}
