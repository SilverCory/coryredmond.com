package site

import (
	"net/http"

	"github.com/SilverCory/coryredmond.com/site/viewdata"
	"github.com/gin-gonic/gin"
)

type ErrorData struct {
	Code    int
	Name    string
	Message string
}

var Error404 = func(ctx *gin.Context) {
	Error(ctx, ErrorData{http.StatusNotFound, "Page Not Found", "The page at this URL is no longer available or never existed."})
}

var Error500 = func(ctx *gin.Context) {
	Error(ctx, ErrorData{http.StatusInternalServerError, "Internal Server Error", "An error occurred processing your request, sorry."})
}

func Error(ctx *gin.Context, data ErrorData) {
	v := viewdata.Default(ctx)
	v.Clear()
	v.Set("Title", data.Name)
	v.Set("Error", data)
	if data.Code == 404 {
		v.Set("Image", "404-bg.jpg")
	} else if data.Code >= 500 && data.Code <= 599 {
		v.Set("Image", "50x-bg.jpg")
	}
	v.Set("OGInfo", map[string]string{
		"og:title": v.GetStringDefault("Title", "Error"),
		"og:type":  "website",
	})
	v.HTML(data.Code, "pages/error.html")
}
