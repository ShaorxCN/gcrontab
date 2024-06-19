package cache

var (
	tokenCache Cache
)

func TokenCacheInit(max int) (err error) {
	tokenCache, err = NewLruCache(max)
	return err
}

func RemoveToken(key string) {
	tokenCache.Remove(key)
}

// SetToken 存放tokenStr 以及salt
func SetToken(key string, value interface{}) {
	tokenCache.Set(key, value)
}

func GetSaltByToken(token string) (interface{}, bool) {
	return tokenCache.Get(token)
}
