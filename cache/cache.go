package cache

import "github.com/tschuyebuhl/livesession/data"

type Cache interface {
	Get(key data.ID) (*data.User, bool, bool)
	Put(value *data.User)
	Delete(key data.ID)
	Nuke(sure bool)
}
