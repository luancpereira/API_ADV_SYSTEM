package errors

import (
	"encoding/json"
	"os"

	"github.com/jellydator/ttlcache/v3"
	"github.com/luancpereira/APICheckout/core/config"
	log "github.com/sirupsen/logrus"
)

var C *ttlcache.Cache[string, string]

type Factory struct{}

func (f Factory) Start() {
	C = f.ttlClient()
}

func (c Factory) loadttlCache(cache *ttlcache.Cache[string, string]) {
	file, err := os.Open(config.ERROR_FILE)
	if err != nil {
		log.Error("error.loading.file.error", "could not find errors.txt file")
	}
	decoder := json.NewDecoder(file)
	var jsonData map[string]string
	decoder.Decode(&jsonData)

	for k, v := range jsonData {
		cache.Set(k, v, ttlcache.NoTTL)
	}

}

func (f *Factory) ttlClient() (ttl *ttlcache.Cache[string, string]) {
	ttl = ttlcache.New[string, string]()
	f.loadttlCache(ttl)
	return ttl
}
