package model

type Task struct {
	ID GUID `gorm:"primary_key"`

	Metadata
	Content
}

type Metadata struct {
	Created   Timestamp
	Modified  Timestamp  `gorm:"index"`
	Completed *Timestamp `gorm:"index"`
	Archived  *Timestamp `gorm:"index"`
}

type Content struct {
	Description string
	Due         Timestamp `gorm:"index"`
}
