package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(path string) (*gorm.DB, error) {
	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := conn.AutoMigrate(&ModelSpec{}, &RoutingLog{}, &AgentKey{}); err != nil {
		return nil, err
	}

	var count int64
	if err := conn.Model(&ModelSpec{}).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		if err := conn.Create(DefaultModels()).Error; err != nil {
			return nil, err
		}
	}

	return conn, nil
}
