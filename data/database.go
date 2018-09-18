package data

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/SilverCory/coryredmond.com/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Handler struct {
	Engine    *gorm.DB
	SQLConfig config.MySQL
}

func NewHandler(sqlConf config.MySQL) (*Handler, error) {

	if !(sqlConf.Enabled) {
		fmt.Println("No SQL conf enabled!")
		return nil, errors.New("no SQL conf enabled")
	}

	var err error

	handler := &Handler{}
	handler.SQLConfig = sqlConf

	handler.Engine, err = gorm.Open("mysql", sqlConf.URI)
	if err != nil {
		return nil, err
	}
	handler.Engine = handler.Engine.Debug()

	err = handler.Sync()
	if err != nil {
		return nil, err
	}

	return handler, nil

}

// $2a$15$CgUx3pr9phONckRzy6fX3eB8RlhNUDFWCJf7qVsndTXTJn7YoScA6
func (h *Handler) Sync() error {
	h.Engine.AutoMigrate(&UserPosts{}, &User{}, &Post{}, &PostText{})
	gob.Register(User{})
	gob.Register(Post{})
	gob.Register(UserPosts{})

	h.Engine.Model(&UserPosts{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	h.Engine.Model(&UserPosts{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")
	h.Engine.Model(&UserPosts{}).AddUniqueIndex("unique_key_index_user_post", "post_id", "user_id")
	return nil
}
