package models

import (
	"gorm.io/gorm"
)

type Votes struct {
	ID            uint `json:"id" gorm:"primaryKey:autoIncrement"`
	UserID        uint `json:"user_id"`
	VotedPaslonID uint `json:"voted_paslon_id"`

	User   Users   `json:"user" gorm:"foreignKey:UserID"`          // Relasi ke model Users
	Paslon Paslons `json:"paslon" gorm:"foreignKey:VotedPaslonID"` // Relasi ke model Paslons
}

func MigrateVotes(db *gorm.DB) error {
	err := db.AutoMigrate(&Votes{})
	return err
}
