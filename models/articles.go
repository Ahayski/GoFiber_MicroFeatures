// models/articles.go

package models

import (
	"gorm.io/gorm"
)

type Articles struct {
	ID          uint   `json:"id" gorm:"primaryKey:autoIncrement"`
	ImgArticle  string `json:"imgarticle" gorm:"type: varchar(255)"`
	ArticleDate string `json:"articledate" gorm:"type: Date"`
	Title       string `json:"title" gorm:"type: varchar(255)"`
	Description string `json:"description" gorm:"type: varchar(255)"`

	// Kunci asing untuk mengaitkan artikel dengan user
	UserID uint `json:"user_id" gorm:"index"`

	User Users `json:"user" gorm:"foreignKey:UserID"`
}

func MigrateArticles(db *gorm.DB) error {
	err := db.AutoMigrate(&Articles{})
	return err
}
