package cache

type Cache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
	Len() int
	Remove(string)
	Cap() int
}
