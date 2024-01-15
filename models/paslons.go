package models

import (
	"gorm.io/gorm"
)

type Paslons struct {
	ID        uint   `json:"id" gorm:"primaryKey:autoIncrement"`
	Nama      string `json:"nama" gorm:"type: varchar(255)"`
	PaslonImg string `json:"paslon_image" gorm:"type: varchar(255)"`
	NoUrut    string `json:"no_urut" gorm:"type: varchar(255)"`
	VisiMisi  string `json:"visi_misi" gorm:"type: varchar(255)"`

	Partais []Partais `json:"partais" gorm:"foreignKey:PaslonID"`    // Relasi ke Partais
	Votes   []Votes   `json:"votes" gorm:"foreignKey:VotedPaslonID"` // Relasi ke model Vote
}

func MigratePaslons(db *gorm.DB) error {
	err := db.AutoMigrate(&Paslons{})
	return err
}
