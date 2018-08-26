package main

import (
	"github.com/SilverCory/coryredmond.com/data"
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

	conf := new(config.Config)
	e(conf.Load())

	blog, err := site.New(conf)
	e(err)

	IndexHandler.RegisterHandlers(blog)
	PostHandler.RegisterHandlers(blog)

	if conf.MySQL.Enabled {
		blog.Data, err = data.NewHandler(conf.MySQL)
		e(err)
	}

	if *address != "" {
		conf.Web.ListenAddress = *address
	}

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
