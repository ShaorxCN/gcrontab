package cache

var (
	saltCache Cache
)

func SaltCacheInit(max int) (err error) {
	saltCache, err = NewLruCache(max)
	return err
}

func RemoveSalt(key string) {
	saltCache.Remove(key)
}

// SetSalt 存放uid 以及salt
func SetSalt(key string, value interface{}) {
	saltCache.Set(key, value)
}

func GetSaltByUID(id string) (interface{}, bool) {
	return saltCache.Get(id)
}
