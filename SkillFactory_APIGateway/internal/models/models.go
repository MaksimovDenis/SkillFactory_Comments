package models

type Feeds struct {
	Id      int    `json:"id" db:"id" bson:"id, omitempty"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
	Link    string `json:"link" db:"link"`
	PubDate string `json:"pub_date" db:"pub_date"`
}

type Filter struct {
	Limit  int    `json:"limit" db:"limit"`
	Filter string `json:"filter" db:"filter"`
}

type Comments struct {
	ID              int32   `json:"ID"`
	NewsID          int32   `json:"NewsID"`
	ParentCommentID *int    `json:"ParentCommentID"`
	Content         string  `json:"Content"`
	CreatedAt       *string `json:"CreatedAt"`
}
