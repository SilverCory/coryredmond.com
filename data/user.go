package data

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique_index"`
	Email    string `gorm:"not null;unique_index"`
	Name     string
	Password string  `gorm:"not null;unique"`
	Posts    []*Post `gorm:"-"`
	PhotoURL string
}
