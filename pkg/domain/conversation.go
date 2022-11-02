package domain

type Conversation struct {
	Id       string `json:"id,omitempty" db:"id"`
	FromId   string `json:"fromId" db:"from_id"`
	ToId     string `json:"toId" db:"to_id"`
	ChatRows int    `json:"chatRows" db:"chat_rows"`
	DayPass  int    `json:"dayPass" db:"day_pass"`
}
