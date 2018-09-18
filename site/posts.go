package site

import (
	"fmt"
	"time"

	"github.com/SilverCory/coryredmond.com/data"
)

var PostPaginationKey = "posts.pagination.cache:LoadPostsViaDate_"

func (b *Blog) LoadPostsViaDate(page int) []data.Post {
	var postsFound []data.Post
	cacheKey := fmt.Sprint(PostPaginationKey, page)
	if err := b.CacheStore.Get(cacheKey, &postsFound); err == nil || len(postsFound) > 0 {
		return postsFound
	} else {
		panic(err)
	}
	return postsFound
}

func (b *Blog) ReloadAllPosts() {
	start := time.Now()
	i := 0
	var postsFound []data.Post
	for len(postsFound) == 10 || i == 0 {
		i++
		lastTime := time.Now()
		if len(postsFound) == 10 {
			lastTime = postsFound[len(postsFound)-1].CreatedAt
		}
		postsFound = b.loadPostsViaDateForce(lastTime, 10, i, true)
	}
	b.TotalPages = i
	fmt.Println("Loaded all posts in: ", time.Now().Sub(start))
}

func (b *Blog) loadPostsViaDateForce(previousTime time.Time, total, page int, override bool) []data.Post {
	var botsFound []data.Post
	cacheKey := fmt.Sprint(PostPaginationKey, page)
	if err := b.CacheStore.Get(cacheKey, &botsFound); err != nil || len(botsFound) < 1 || override {
		b.Data.Engine = b.Data.Engine.Debug()
		b.Data.Engine.Model(&data.Post{}).Order("created_at ASC").Where("(created_at < ?)", previousTime).Limit(total).Find(&botsFound)
		for k, v := range botsFound {
			if _, err := v.GetAuthor(b.Data); err != nil {
				panic(err)
			}
			botsFound[k] = v
		}
		if err := b.CacheStore.Set(cacheKey, botsFound, 3*time.Hour); err != nil {
			// TODO error handling?
			panic(err)
		}
	}

	return botsFound
}

// ヽ( •_)ᕗ
