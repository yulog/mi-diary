package repository

type CacheRepositorier interface {
	Set(host, key string, value []byte) error
	Get(host, key string) ([]byte, error)

	Close()
}
