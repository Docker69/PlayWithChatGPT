package models

//model for the auto ai, includes human id, chat id, ai name, ai role and ai goals

type AutoAI struct {
	Id      string   `json:"id" bson:"_id,omitempty"`
	ChatId  string   `json:"chat_id"`
	HumanId string   `json:"human_id"`
	Name    string   `json:"name"`
	Role    string   `json:"role"`
	Goals   []string `json:"goals"`
}
