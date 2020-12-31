package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	maxPageSize     = 100
	defaultPageSize = 10
)

func NewDB(url string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&InchAction{})
	return db, nil
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > maxPageSize:
			pageSize = maxPageSize
		case pageSize <= 0:
			pageSize = defaultPageSize
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func GetInchAction(db *gorm.DB, page, pageSize int) ([]*InchAction, error) {
	var result []*InchAction
	if err := db.Scopes(Paginate(page, pageSize)).Find(&result).Error; err == nil {
		return result, nil
	} else {
		return nil, err
	}
}

func InsertSwapResult(db *gorm.DB, record *InchAction) error {
	return db.Create(record).Error
}
