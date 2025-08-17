package model

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `json:"password"`
}
