package handler

import (
	"github.com/arthur404dev/gopportunities/config"
	"gorm.io/gorm"
)

var (
	logger *config.Logger
	db     *gorm.DB
)

func InitializeHandler() {
	logger = config.GetLogger("handler")
	db = config.GetSQLite()
}

// SetDB injects the database dependency used by the handlers. It exists so the
// package can be exercised against an isolated database without booting the
// full application. The logger is initialized if it has not been already.
func SetDB(database *gorm.DB) {
	if logger == nil {
		logger = config.GetLogger("handler")
	}
	db = database
}
