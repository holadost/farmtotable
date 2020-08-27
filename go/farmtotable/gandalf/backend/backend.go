package backend

import "github.com/jinzhu/gorm"

type SqlBackend interface {
	GetDB() *gorm.DB
}
