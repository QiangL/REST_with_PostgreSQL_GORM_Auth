package models

type Todo struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      uint   `json:"userID"`
	Done        bool   `json:"done"`
}
