package db

import "gorm.io/gorm"

func InitWorkSpaceDb(p string) error {
	_db, err := GetWorkSpaceGormDb(p)
	if err != nil {
		return err
	}
	return Migrate(_db)
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&DbApiConfigCollection{},
		&DbApiConfig{},
		&DbData{},
	)
}
