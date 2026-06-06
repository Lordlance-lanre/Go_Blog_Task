package models

type Blog struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"foreignKey:UserID"`
}
