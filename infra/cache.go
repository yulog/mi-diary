package infra

import (
	"log"

	"github.com/yulog/mi-diary/domain/repository"
)

type CacheInfra struct {
	dao *DataBase
}

func (i *Infra) NewCacheInfra() repository.CacheRepositorier {
	return &CacheInfra{dao: i.dao}
}

func (i *CacheInfra) Set(host, key string, value []byte) error {
	return i.dao.Cache(host).Put([]byte(key), value)
}

func (i *CacheInfra) Get(host, key string) ([]byte, error) {
	log.Println("cache get")
	return i.dao.Cache(host).Get([]byte(key))
}

func (i *CacheInfra) Close() {
	i.dao.cache.Range(func(key, value any) bool {
		value.(*pogreb.DB).Close()
		return true
	})
}
