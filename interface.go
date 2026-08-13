package lrucache

type Cache interface {
	Set(key string, value int)
	Get(key string) (int, bool)
}
