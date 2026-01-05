package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}

type Task struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Done     bool   `json:"done"`
	Category string `json:"category"`
	Priority int    `json:"priority"`
	Position int    `json:"position"`
	UserID   uint   `json:"user_id"`
}
