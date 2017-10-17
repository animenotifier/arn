package arn

// GroupMember ...
type GroupMember struct {
	UserID UserID  `json:"userId"`
	Role   string  `json:"role"`
	Joined UTCDate `json:"joined"`
}
