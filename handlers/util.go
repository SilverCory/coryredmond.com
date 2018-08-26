package handlers

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dineshappavoo/basex"
)

func GetPostURL(id, title string) string {
	return url.PathEscape(strings.Replace(strings.TrimSpace(title), " ", "-", -1) + "-" + id)
}

func GetPostIDFromURL(urlPath string) string {
	urlPath, err := url.PathUnescape(urlPath)
	if err != nil {
		return ""
	}

	parts := strings.Split(urlPath, "-")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}

func GeneratePostID() (uint64, string, error) {
	rand.Seed(time.Now().Unix())
	id := rand.Uint64()
	idString, err := basex.Encode(strconv.FormatUint(id, 10))
	return id, idString, err
}

func DecodeID(id string) (uint64, error) {
	str, err := basex.Decode(id)
	if err != nil {
		return uint64(0), err
	}

	return strconv.ParseUint(str, 10, 64)
}
