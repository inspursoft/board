package commons

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type cacheStore struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type CacheStoreController struct {
	BaseController
}

func (c *CacheStoreController) Prepare() {
	c.EnableXSRF = false
}

func (c *CacheStoreController) Post() {
	var cs cacheStore
	err := c.ResolveBody(&cs)
	if err != nil {
		c.InternalError(err)
	}
	MemoryCache.Put(cs.Key, cs.Value, DefaultCacheDuration)
	logs.Debug("Successfully stored key: %s with value: %+v", cs.Key, cs.Value)
}

func (c *CacheStoreController) Get() {
	key := c.GetString("key")
	if MemoryCache.IsExist(key) {
		val := MemoryCache.Get(key)
		logs.Debug("Found store: %+v by key: %s", val, key)
		c.Data["json"] = cacheStore{Key: key, Value: val}
		c.ServeJSON()
	}
	c.CustomAbortAudit(http.StatusNotFound, fmt.Sprintf("No value found from cache store with key: %s", key))
}
