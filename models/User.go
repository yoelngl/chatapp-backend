package models

import "time"

type User struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username  string    `gorm:"not null;unique" json:"username"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `gorm:"default:CURENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURENT_TIMESTAMP()" json:"updated_at"`
}
