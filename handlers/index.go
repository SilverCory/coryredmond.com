package handlers

import (
	"github.com/SilverCory/coryredmond.com/data"
	"github.com/SilverCory/coryredmond.com/site"
	"github.com/SilverCory/coryredmond.com/site/viewdata"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Index struct {
	site.Handler
	blog *site.Blog
}

func (i *Index) RegisterHandlers(b *site.Blog) error {
	i.blog = b
	b.Gin.GET("/", func(ctx *gin.Context) {
		v := viewdata.Default(ctx)
		v.Set("Title", "Home")
		v.Set("OGInfo", map[string]string{
			"og:title":           "Home",
			"og:type":            "profile",
			"og:image":           "https://coryredmond.com/img/home-bg.jpg",
			"profile:first_name": "Cory",
			"profile:last_name":  "Redmond",
			"profile:username":   "CoryOry",
			"profile:gender":     "Male",
		})
		v.HTML(200, "pages/index.html")
	})

	b.Gin.GET("/logout", func(ctx *gin.Context) {
		sesh := sessions.Default(ctx)
		sesh.Clear()
		if err := sesh.Save(); err != nil {
			// TODO proper error pages
			ctx.JSON(500, err)
		}
		ctx.Redirect(302, "/")
	})

	b.Gin.GET("/login", func(ctx *gin.Context) {
		v := viewdata.Default(ctx)
		v.Set("Title", "Login")
		v.Set("OGInfo", map[string]string{
			"og:title": "Login",
			"og:type":  "website",
			"og:image": "https://coryredmond.com/img/home-bg.jpg",
		})
		v.HTML(200, "pages/login.html")
	})

	b.Gin.POST("/login", i.handleLoginPost)

	b.Gin.NoRoute(site.Error404)
	b.Gin.GET("/500", site.Error500)
	b.Gin.GET("/500viapanic", func(ctx *gin.Context) {
		panic("ahhh")
	})

	return nil
}

func (i *Index) handleLoginPost(ctx *gin.Context) {
	username, exists := ctx.GetPostForm("username")
	if !exists {
		ctx.String(400, "You can't log in because it's not implemented!")
		return
	}

	password, exists := ctx.GetPostForm("password")
	if !exists {
		ctx.String(400, "You can't log in because it's not implemented!")
		return
	}

	user := &data.User{Username: username}
	if err := i.blog.Data.Engine.Find(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.String(400, "Invalid username and/or password..")
		} else {
			ctx.String(500, "Error loading from database: %q", err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.String(400, "Invalid username and/or password..")
	}

	sesh := sessions.Default(ctx)
	sesh.Set("user", *user)
	if err := sesh.Save(); err != nil {
		ctx.String(500, "Unable to log you in: %q", err)
	}
}
