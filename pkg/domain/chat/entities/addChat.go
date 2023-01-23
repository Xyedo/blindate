package chatEntity

type New struct {
	Message string  `json:"message" binding:"required,max=4096"`
	ReplyTo *string `json:"replyTo" binding:"omitempty,uuid"`
}
