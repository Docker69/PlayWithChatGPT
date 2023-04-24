package models

type ChatRecord struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}
type Human struct {
	Id       string       `json:"id" bson:"_id,omitempty"`
	Name     string       `json:"name"`
	NickName string       `json:"nickName"`
	ChatIds  []ChatRecord `json:"chatIds"`
}
