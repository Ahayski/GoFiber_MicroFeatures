package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	// Membuat string koneksi (DSN) dari konfigurasi
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s ",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Membuka koneksi ke database menggunakan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Mengembalikan instance *gorm.DB dan error (jika ada)
	return db, err
}
