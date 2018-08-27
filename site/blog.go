package site

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/SilverCory/coryredmond.com/config"
	"github.com/SilverCory/coryredmond.com/data"
	"github.com/SilverCory/coryredmond.com/util"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"

	"github.com/gomodule/redigo/redis"
)

type Blog struct {
	Gin         *gin.Engine
	Config      *config.Config
	Redis       *redis.Pool
	StaticStore persistence.CacheStore
	Data        *data.Handler
}

// Handler an interface for sections of the site that handle.
type Handler interface {
	RegisterHandlers(b *Blog) error
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

	if err = ret.loadPrePagesAndMiddleware(); err != nil {
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
	b.Gin.Use(RecoveryWithWriter(io.MultiWriter(errF, os.Stderr)))
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
func (b *Blog) loadPrePagesAndMiddleware() error {

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

	if b.Redis != nil {
		b.StaticStore = persistence.NewRedisCacheWithPool(b.Redis, 3*time.Hour)
	} else {
		b.StaticStore = persistence.NewInMemoryStore(3 * time.Hour)
	}

	return nil
}

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(out io.Writer) gin.HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					stack := util.Stack(3)
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					logger.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s", time.Now().Format("2006/01/02 - 15:04:05"), string(httprequest), err, stack, string([]byte{27, 91, 48, 109}))
				}
				if XHR(c) {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
				} else {
					Error500(c)
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
