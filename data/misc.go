package data

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type UserPosts struct {
	gorm.Model
	UserId uint
	PostId uint
}

func (up *UserPosts) Set(h *Handler) error {
	if up.UserId == 0 || up.PostId == 0 {
		return errors.New("bad ID")
	}

	return h.Engine.FirstOrCreate(up, up).Error
}

func (p *Post) GetAuthor(h *Handler) (User, error) {
	var user User
	userPosts := new(UserPosts)

	if err := h.Engine.First(userPosts, "post_id = ?", p.ID).Error; err != nil {
		return user, nil
	}

	err := h.Engine.First(&user, "id = ?", userPosts.UserId).Error
	p.Author = user
	return user, err
}

func (p *Post) SetAuthor(h *Handler) error {
	u := &UserPosts{UserId: p.Author.ID, PostId: p.ID}
	fmt.Printf("%#v\n", u)
	return u.Set(h)
}

func (u *User) GetPosts(h *Handler) ([]*Post, error) {
	var posts []*Post
	var userPosts []*UserPosts
	var err error

	if err := h.Engine.Model(&UserPosts{}).Where("user_id = ?", u.ID).Error; err != nil {
		return posts, err
	}

	for _, v := range userPosts {
		var post *Post
		err = h.Engine.First(post, "id = ?", v.PostId).Error
		if err == nil {
			posts = append(posts, post)
		}
	}

	u.Posts = posts
	return posts, err
}
