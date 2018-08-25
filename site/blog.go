package site

import (
	"html/template"
	"io"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gomodule/redigo/redis"

	"github.com/SilverCory/coryredmond.com/config"
	"github.com/gin-gonic/gin"
)

type Blog struct {
	Gin    *gin.Engine
	Config *config.Config
	Redis  *redis.Pool
	// TODO add data.
}

// Handler an interface for sections of the site that handle.
type Handler interface {
	RegisterHandlers() error
}

func New(conf *config.Config) (ret *Blog, err error) {
	ret = &Blog{
		Gin:    gin.New(),
		Config: conf,
	}

	// Setup the logging
	if err = ret.setUpLogs(); err != nil {
		return
	}

	ret.configureTemplates()

	if err = ret.registerPagesAndMiddleware(); err != nil {
		return
	}

	return

}

// Set up logging to files
func (b *Blog) setUpLogs() error {
	if f, err := os.Stat(b.Config.Web.LogDirectory); os.IsNotExist(err) || !f.IsDir() {
		if err := os.MkdirAll(b.Config.Web.LogDirectory, 0644); err != nil {
			return err
		}
	}

	f, err := os.Create(b.Config.Web.LogDirectory + "/gin.log")
	if err != nil {
		return err
	}

	errF, err := os.Create(b.Config.Web.LogDirectory + "/gin_err.log")
	if err != nil {
		return err
	}

	b.Gin.Use(gin.LoggerWithWriter(io.MultiWriter(f)))
	b.Gin.Use(gin.RecoveryWithWriter(io.MultiWriter(errF, os.Stderr)))
	return nil
}

func (b *Blog) configureTemplates() {
	//Load the HTML templates and templating functions.
	b.Gin.SetFuncMap(template.FuncMap{
		"comments": func(s string) template.HTML { return template.HTML(s) },
		"ASCII":    GetAscii,
	})
	b.Gin.LoadHTMLGlob(b.Config.Web.TemplateGlob)
}

// Setup pages and middleware
func (b *Blog) registerPagesAndMiddleware() error {

	// Static files to load
	b.Gin.Use(static.Serve("/", static.LocalFile(b.Config.Web.StaticFilePath, false)))

	// Routes and middleware are registered in order.
	// Any pages that require no middleware should go above middleware addition
	// Same goes for any pages that can have the pre middleware, and post middleware.

	if err := b.AddPreMiddleware(); err != nil {
		return err
	}

	if err := b.AddPostMiddleware(); err != nil {
		return err
	}
	//
	//oauth := Oauth{Web: ret}
	//oauth.RegisterHandlers()
	//
	//pages := &Pages{Web: ret}
	//pages.RegisterHandlers()
	return nil
}
