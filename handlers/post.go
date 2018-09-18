package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/SilverCory/coryredmond.com/data"

	"github.com/SilverCory/coryredmond.com/site"
	"github.com/SilverCory/coryredmond.com/site/viewdata"
	"github.com/SilverCory/coryredmond.com/util"
	"github.com/gin-contrib/cache"
	"github.com/gin-gonic/gin"
)

type Post struct {
	site.Handler
	blog *site.Blog
}

func (p *Post) RegisterHandlers(b *site.Blog) error {
	p.blog = b

	b.Gin.GET("/post_preview/*page", cache.CachePageWithoutHeader(b.CacheStore, time.Hour*3, p.handlePostPreview))
	//b.Gin.GET("/post_content/:post", cache.CachePage(p.blog.CacheStore, time.Hour*3, p.handlePostContent))

	// fixme https://github.com/julienschmidt/httprouter/issues/73
	b.Gin.GET("/p/:post", cache.CachePage(p.blog.CacheStore, time.Hour*3, p.handlePost))
	return nil
}

func (p *Post) handlePostPreview(ctx *gin.Context) {
	pageStr := ctx.Param("page")

	var page int
	var err error
	pageStr = pageStr[1:]
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			panic(err)
		}
	}

	if page < 1 || page > p.blog.TotalPages {
		site.Error404(ctx)
		ctx.Abort()
		return
	}

	v := viewdata.Default(ctx)

	nextPage := page + 1
	if nextPage <= p.blog.TotalPages {
		v.Set("NextPage", nextPage)
	} else {
		v.Set("NextPage", -1)
	}

	posts := p.blog.LoadPostsViaDate(page)
	v.Set("Posts", posts)
	v.HTML(200, "pages/post_preview.html")
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

	post := data.Post{
		PostURLID: id,
	}

	// Load blog
	if err := p.blog.Data.Engine.First(&post, &post).Error; err != nil {
		ctx.Error(err)
		site.Error500(ctx)
		fmt.Println(err)
		return
	}

	// Load author
	author, err := post.GetAuthor(p.blog.Data)
	if err != nil {
		ctx.Error(err)
		site.Error500(ctx)
		fmt.Println(err)
		return
	}

	// Load fulltext
	text, _, err := post.GetFullText(p.blog.Data)
	if err != nil {
		ctx.Error(err)
		site.Error500(ctx)
		fmt.Println(err)
		return
	}

	v := viewdata.Default(ctx)
	v.Set("Title", post.Title)

	OGInfo := map[string]string{
		"og:title":               post.Title,
		"og:type":                "article",
		"article:published_time": post.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		"article:author":         "https://coryredmond.com/",
		"article:section":        "Front Page",
	}

	if post.CreatedAt.Before(post.UpdatedAt) {
		OGInfo["article:modified_time"] = post.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z")
	}

	// TODO tags

	v.Set("OGInfo", OGInfo)
	v.Set("Post", post)
	v.Set("Author", author)
	v.Set("FullText", text)
	v.HTML(200, "pages/post.html")
}
