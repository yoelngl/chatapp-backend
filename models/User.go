package models

type User struct {
	ID          string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username    string `gorm:"not null;unique" json:"username"`
	Email       string `gorm:"not null;unique" json:"email"`
	Password    string `gorm:"not null" json:"-"`
	FullName    string `json:"full_name"`
	DateOfBirth string `json:"date_of_birth"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
