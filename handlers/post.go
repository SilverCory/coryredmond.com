package handlers

import (
	"github.com/SilverCory/coryredmond.com/site"
	"github.com/SilverCory/coryredmond.com/site/viewdata"
	"github.com/SilverCory/coryredmond.com/util"
	"github.com/gin-gonic/gin"
)

type Post struct {
	site.Handler
}

func (p *Post) RegisterHandlers(b *site.Blog) error {
	// fixme https://github.com/julienschmidt/httprouter/issues/73
	b.Gin.GET("/p/:post", p.handlePost)
	return nil
}

func (p *Post) handlePost(ctx *gin.Context) {

	// TODO get post
	postUrl := ctx.Param("post")
	postUrl = util.GetPostIDFromURL(postUrl)
	if postUrl == "" {
		site.Error404(ctx)
		return
	}

	id, err := util.DecodeID(postUrl)
	if err != nil || id == 0 {
		site.Error404(ctx)
		return
	}

	v := viewdata.Default(ctx)
	v.Set("Title", "Home")
	v.Set("OGInfo", map[string]string{
		"og:title":           "Home",
		"og:type":            "profile",
		"profile:first_name": "Cory",
		"profile:last_name":  "Redmond",
		"profile:username":   "CoryOry",
		"profile:gender":     "Male",
	})
	v.HTML(200, "pages/post.html")
}
