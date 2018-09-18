package util

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dineshappavoo/basex"
)

// Constant definitions
const MaxUInt = ^uint(0)
const MinUInt = 0

const MaxInt = int(^uint(0) >> 1)
const MinInt = -MaxInt - 1

func GetPostURL(id, title string) string {
	title += " "
	newTitle := ""
	for k, v := range title {
		if k >= 42 {
			if newTitle == "" {
				newTitle = title[:42]
			}
			break
		} else {
			if unicode.IsSpace(v) {
				newTitle = title[:k]
			}
		}
	}

	return url.PathEscape(strings.Replace(strings.TrimSpace(newTitle), " ", "-", -1) + "-" + id)
}

func GetPostIDFromURL(urlPath string) string {
	urlPath, err := url.PathUnescape(urlPath)
	if err != nil {
		return ""
	}

	parts := strings.Split(urlPath, "-")
	if len(parts) == 1 {
		return parts[0]
	} else if len(parts) < 1 {
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
