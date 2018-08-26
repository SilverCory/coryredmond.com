package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/SilverCory/coryredmond.com/config"
	"github.com/SilverCory/coryredmond.com/data"
	flag "github.com/ogier/pflag"
)

func main() {
	configLoc := flag.StringP("configLoc", "c", "./config.json", "The config location")
	username := flag.StringP("username", "u", "cory", "The user name of the account")
	email := flag.StringP("email", "e", "", "The email address of the account")
	flag.Parse()

	conf := new(config.Config)
	config.SaveLocation = *configLoc
	e(conf.Load())

	if !conf.MySQL.Enabled {
		fmt.Println("MySQL isn't enabled...")
		return
	}

	password := ""

	fmt.Print("Enter your password: ")
	fmt.Scan(&password)

	dataHandler, err := data.NewHandler(conf.MySQL)
	e(err)

	user := &data.User{Username: *username}
	if err := dataHandler.Engine.Find(user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			e(err)
		}
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	e(err)

	user.Password = string(pw)
	user.Email = *email

	if user.CreatedAt.IsZero() {
		err = dataHandler.Engine.Create(user).Error
	} else {
		err = dataHandler.Engine.Save(user).Error
	}
	e(err)

}

func e(e error) {
	if e != nil {
		panic(e)
	}
}
