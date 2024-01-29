package cache

import "github.com/tschuyebuhl/livesession/data"

type Cache interface {
	Get(key string) (*data.User, bool)
	Put(value *data.User) error
	Delete(key string)
	Nuke(sure bool)
}
