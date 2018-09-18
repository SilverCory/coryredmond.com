package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/SilverCory/coryredmond.com/data"
	"github.com/gin-contrib/cache/persistence"

	"github.com/SilverCory/coryredmond.com/util"

	"github.com/gomodule/redigo/redis"

	"github.com/SilverCory/coryredmond.com/handlers"
	"github.com/SilverCory/coryredmond.com/site"
	flag "github.com/ogier/pflag"

	"github.com/SilverCory/coryredmond.com/config"
)

var (
	IndexHandler = new(handlers.Index)
	PostHandler  = new(handlers.Post)
)

func main() {

	address := flag.StringP("address", "a", "", "The address to listen on (overrides)")
	disableSSL := flag.BoolP("disablessl", "s", false, "Disable SSL")
	flag.Parse()

	// Config load
	conf := new(config.Config)
	e(conf.Load())

	// Create the blog instance.
	blog, err := site.New(conf)
	e(err)

	// Setup redis
	if conf.Redis.Enabled {
		blog.Redis = &redis.Pool{
			MaxIdle:     10,
			IdleTimeout: 240 * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Dial: func() (redis.Conn, error) {
				return util.DialWithDB(conf.Redis.Network, conf.Redis.Address, conf.Redis.Password, conf.Redis.Database)
			},
		}
	}

	// Along with redis add the cache store.
	if blog.Redis != nil {
		blog.CacheStore = persistence.NewRedisCacheWithPool(blog.Redis, 3*time.Hour)
	} else {
		blog.CacheStore = persistence.NewInMemoryStore(3 * time.Hour)
	}

	// Setup Database
	if conf.MySQL.Enabled {
		blog.Data, err = data.NewHandler(conf.MySQL)
		e(err)
	}

	user := new(data.User)
	user.Username = "cory"
	user.Email = "me@cory.red"

	enigne := blog.Data.Engine.Debug()
	e(enigne.First(user, user).Error)

	fmt.Printf("%#v\n", user)

	post0 := data.Post{
		Author:     *user,
		URL:        "URL1",
		FullTextID: 0,
		PostURLID:  12312311123,
		Summary:    "A short summary",
		Title:      "Hello World",
	}

	post1 := data.Post{
		Author:     *user,
		URL:        "URL2",
		FullTextID: 5,
		PostURLID:  11235345123123,
		Summary:    "A slightly longer summary for you",
		Title:      "Goodbye moon.",
	}

	e(enigne.FirstOrCreate(&post0, post0).Error)
	fmt.Printf("%#v\n", post0)
	e(post0.SetAuthor(blog.Data))
	e(enigne.FirstOrCreate(&post1, post1).Error)
	fmt.Printf("%#v\n", post0)
	e(post1.SetAuthor(blog.Data))

	// Load and configure middleware
	if err = blog.LoadPrePagesAndMiddleware(); err != nil {
		e(err)
	}

	// Register handlers
	PostHandler.RegisterHandlers(blog) // Before middleware
	IndexHandler.RegisterHandlers(blog)

	go func() {
		blog.ReloadAllPosts()
		fmt.Println("Loaded all posts, total pages: " + strconv.Itoa(blog.TotalPages))
		for {
			time.Sleep(2 * time.Hour)
			blog.ReloadAllPosts()
			fmt.Println("Reloaded all posts, total pages: " + strconv.Itoa(blog.TotalPages))
		}
	}()

	// Configure override address
	if *address != "" {
		conf.Web.ListenAddress = *address
	}

	// Run with ssl or not
	if *disableSSL {
		e(blog.Gin.Run(conf.Web.ListenAddress))
	} else {
		e(blog.RunAutoTLS())
	}

}

func e(e error) {
	if e != nil {
		panic(e)
	}
}
