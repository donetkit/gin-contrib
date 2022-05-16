package cache

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"github.com/donetkit/gin-contrib/utils/cache"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	PageCachePrefix = "gincontrib.page.cache"
	ErrCacheMiss    = errors.New("cache: key not found.")
	ErrNotStored    = errors.New("cache: not stored.")
	ErrNotSupport   = errors.New("cache: not support.")
)

type responseCache struct {
	status int
	header http.Header
	data   []byte
}

type cachedWriter struct {
	c *gin.Context
	gin.ResponseWriter
	status  int
	written bool
	store   cache.IShortCache
	expire  time.Duration
	key     string
}

func SiteCache(store cache.IShortCache, expire time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL
		key := urlEscape(PageCachePrefix, url.RequestURI())
		var cacheData = store.WithDB(0).WithContext(c.Request.Context()).Get(key)
		if cacheData == nil {
			c.Next()
			return
		}
		cache, ok := cacheData.(responseCache)
		if !ok {
			c.Next()
		}
		c.Writer.WriteHeader(cache.status)
		for k, val := range cache.header {
			for _, v := range val {
				c.Writer.Header().Add(k, v)
			}
		}
		c.Writer.Write(cache.data)
	}
}

func CachePage(store cache.IShortCache, expire time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL
		key := urlEscape(PageCachePrefix, url.RequestURI())
		cacheData := store.WithDB(0).WithContext(c.Request.Context()).Get(key)
		if cacheData == nil { // replace writer
			writer := newCachedWriter(c, store, expire, c.Writer, key)
			c.Writer = writer
			handle(c)
		}
		cache, ok := cacheData.(responseCache) // replace writer
		if ok {
			c.Writer.WriteHeader(cache.status)
			for k, vals := range cache.header {
				for _, v := range vals {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Write(cache.data)
		}
	}
}

func NewPageCache(store cache.IShortCache, expire time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL
		key := urlEscape(PageCachePrefix, url.RequestURI())
		cacheData := store.WithDB(0).WithContext(c.Request.Context()).Get(key)
		if cacheData == nil {
			writer := newCachedWriter(c, store, expire, c.Writer, key) // replace writer
			c.Writer = writer
			handle(c)
			return
		}
		cache, ok := cacheData.(responseCache)
		if ok {
			c.Writer.WriteHeader(cache.status)
			for k, vals := range cache.header {
				for _, v := range vals {
					c.Writer.Header().Add(k, v)
				}
			}
			c.Writer.Write(cache.data)
		}

	}
}

func newCachedWriter(c *gin.Context, store cache.IShortCache, expire time.Duration, writer gin.ResponseWriter, key string) *cachedWriter {
	return &cachedWriter{c, writer, 0, false, store, expire, key}
}

func (w *cachedWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachedWriter) Status() int {
	return w.status
}

func (w *cachedWriter) Written() bool {
	return w.written
}

func (w *cachedWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err == nil {
		//cache response
		store := w.store
		val := responseCache{
			w.status,
			w.Header(),
			data,
		}
		err = store.WithDB(0).WithContext(w.c.Request.Context()).Set(w.key, val, w.expire)
		if err != nil {
			// need logger
		}
	}
	return ret, err
}

func urlEscape(prefix string, u string) string {
	key := url.QueryEscape(u)
	if len(key) > 200 {
		h := sha1.New()
		io.WriteString(h, u)
		key = string(h.Sum(nil))
	}
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}
