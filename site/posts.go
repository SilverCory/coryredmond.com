package site

import (
	"fmt"

	"github.com/SilverCory/coryredmond.com/data"
)

func (b *Blog) LoadPostsViaDate(page int) []data.Post {
	var postsFound []data.Post
	cacheKey := fmt.Sprint("posts.pagination.cache:LoadPostsViaDate_", page)
	if err := b.StaticStore.Get(cacheKey, &postsFound); err == nil || len(postsFound) > 0 {
		return postsFound
	}
	return postsFound
}

// TODO load in to redis
// TODO load in to redis if len(postsFound) == 0
// ヽ( •_)ᕗ
