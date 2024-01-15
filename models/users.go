// models/user.go

package models

import (
	"gorm.io/gorm"
)

type Users struct {
	ID       uint   `json:"id" gorm:"primaryKey:autoIncrement"`
	FullName string `json:"fullname" gorm:"type: varchar(255)"`
	Address  string `json:"address" gorm:"type: varchar(255)"`
	Gender   string `json:"Gender" gorm:"type: varchar(255)"`
	UserName string `json:"username" gorm:"type: varchar(255)"`
	Password string `json:"password" gorm:"type: varchar(255)"`

	// Relasi dengan Articles, satu user bisa memiliki banyak articles
	Articles []Articles `json:"articles" gorm:"foreignKey:UserID"`
	Votes    []Votes    `json:"votes" gorm:"foreignKey:UserID"` // Relasi ke model Vote
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}
