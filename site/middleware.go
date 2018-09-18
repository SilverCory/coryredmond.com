package site

import (
	"regexp"
	"time"

	"github.com/SilverCory/gin-redisgo-cooldowns"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	redstore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"github.com/utrack/gin-csrf"
)

var allowedInRegisterRegex = regexp.MustCompile(`(?i)^(/(logout|tos|privacy|((vendor|js|css|img|auth)/*.)))|/$`)

const CSP = `
default-src 'self';
img-src 'self' https://cdnjs.cloudflare.com/ https://placekitten.com/ https://cdn.discordapp.com/;
script-src 'self' https://cdnjs.cloudflare.com/ajax/libs/cookieconsent2/3.0.3/cookieconsent.min.js 'sha256-SplWdsqEBp8LjzZSKYaEfDXhXSi0/oXXxAnQSYREAuI=';
style-src 'self' https://cdnjs.cloudflare.com/ajax/libs/cookieconsent2/3.0.3/cookieconsent.min.css https://fonts.googleapis.com 'unsafe-inline';
font-src 'self' https://fonts.gstatic.com;
`

type Middleware struct {
	blog *Blog
}

var m *Middleware

func (b *Blog) AddPreMiddleware() (err error) {
	m = &Middleware{b}

	if err = m.setupSessions(); err != nil {
		return
	}

	return
}

func (b *Blog) AddPostMiddleware() (err error) {
	m.setupCors()
	m.setupCsrf()
	m.setupSecurity()
	m.setupIPCooldowns()

	return
}

func (m *Middleware) setupCsrf() {
	m.blog.Gin.Use(csrf.Middleware(csrf.Options{
		Secret: m.blog.Config.Web.CSRFSecret,
		ErrorFunc: func(c *gin.Context) {

			if c.Request.URL.Path == "/CSPReport" {
				return
			}

			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
}

func (m *Middleware) setupSessions() (err error) {
	conf := m.blog.Config.Redis
	if !(conf.Enabled) {
		store := cookie.NewStore([]byte("dankest_cory_blog_ever"))
		m.blog.Gin.Use(sessions.Sessions("coryredmond_sessions", store))
		return
	}

	store, err := redstore.NewStoreWithPool(m.blog.Redis, []byte("dankest_cory_blog_ever"))
	if err != nil {
		return
	}

	if redstore.SetKeyPrefix(store, "coryredmondblog.sessions.sesh:"); err != nil {
		panic(err)
	}

	store.Options(sessions.Options{
		Secure:   true,
		MaxAge:   int(((24 + time.Hour) * 7).Seconds()),
		HttpOnly: true,
		Domain:   "coryredmond.com",
	})

	m.blog.Gin.Use(sessions.Sessions("coryredmond_sessions", store))
	return nil
}

func (m *Middleware) setupIPCooldowns() {
	if m.blog.Config.Redis.Enabled {
		m.blog.Gin.Use(gin_redisgo_cooldowns.NewRateLimit(m.blog.Redis, "coryredmond.cooldown.general.ip:", 100, time.Second*5, nil))
	}
}

func (m *Middleware) setupCors() {
	if gin.IsDebugging() {
		return
	}

	m.blog.Gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://cdnjs.cloudflare.com", "https://fonts.googleapis.com", "https://coryredmond.com"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
}

func (m *Middleware) setupSecurity() {
	sec := secure.New(secure.Options{
		AllowedHosts:            []string{"coryredmond.com"},
		SSLRedirect:             false,
		SSLTemporaryRedirect:    false,
		SSLHost:                 "coryredmond.com",
		STSSeconds:              86400,
		STSIncludeSubdomains:    true,
		STSPreload:              true,
		ForceSTSHeader:          true,
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
		ContentTypeNosniff:      true,
		BrowserXssFilter:        true,
		ContentSecurityPolicy:   CSP,
		HostsProxyHeaders:       []string{},

		IsDevelopment: gin.IsDebugging(),
	})

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := sec.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	m.blog.Gin.Use(secureFunc)
}
