package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique_index"`
	Email    string `gorm:"not null;unique_index"`
	Name     string
	Password string `gorm:"not null;unique"`
	PhotoURL string
	Posts    []Post `gorm:"many2many:user_posts;"`
}

type Post struct {
	ID        uint64 `gorm:"primary_key"`
	Title     string `gorm:"not null;unique_index"`
	Author    User   `gorm:"many2many:user_posts;"`
	URL       string `gorm:"not null;unique_index"`
	Summary   string
	FullText  string `gorm:"type:LONGTEXT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
