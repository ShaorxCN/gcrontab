package cache

var (
	saltCache Cache
)

func SaltCacheInit(max int) (err error) {
	saltCache, err = NewLruCache(max)
	return err
}

func RemoveUIDSalt(id string) {
	saltCache.Remove(id)
}

func RemoveByUIDAndToken(uid, token string) {
	inner, ok := saltCache.Get(uid)
	if !ok {
		return
	}
	inner.(*LruCache).Remove(token)
}

// SetSalt 存放uid token  以及salt
func SetSalt(uid, token, salt string) error {
	inner, ok := saltCache.Get(uid)
	if !ok {
		inner, err := NewLruCache(saltCache.Cap())

		if err != nil {
			return err
		}

		saltCache.Set(uid, inner)
		inner.Set(token, salt)
		return nil
	}
	inner.(*LruCache).Set(token, salt)
	return nil
}

func GetSaltByUIDAndToken(uid, token string) (string, bool) {
	var ok bool
	var inner interface{}
	var ret interface{}
	inner, ok = saltCache.Get(uid)

	if !ok {
		return "", false
	}

	ret, ok = inner.(*LruCache).Get(token)
	if !ok {
		return "", false
	}

	return ret.(string), true

}
