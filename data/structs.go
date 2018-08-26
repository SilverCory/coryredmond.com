package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null;unique"`
	PhotoURL string
	Posts    []Post `gorm:"many2many:user_posts;"`
}

type Post struct {
	ID        uint64 `gorm:"primary_key"`
	Title     string `gorm:"not null;unique"`
	Author    []User `gorm:"many2many:user_posts;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
