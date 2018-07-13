package app

import (
	"github.com/jinzhu/gorm"
	"github.com/qeelyn/go-common/logger"
)

var (
	IsDebug bool
	Db      *gorm.DB
	Logger  *logger.Logger
)
