package models

import (
	"gorm.io/gorm"
)

type Partais struct {
	ID uint `json:"id" gorm:"primaryKey:autoIncrement"`

	PaslonID uint `json:"paslon_id" gorm:"index"` // Kunci asing ke Paslons

	PartaiImg string `json:"partai_image" gorm:"type: varchar(255)"`
	Nama      string `json:"nama" gorm:"type: varchar(255)"`
	KetuaUmum string `json:"ketua_umum" gorm:"type: varchar(255)"`
	VisiMisi  string `json:"visi_misi" gorm:"type: varchar(255)"`
	Alamat    string `json:"alamat" gorm:"type: varchar(255)"`
}

func MigratePartais(db *gorm.DB) error {
	err := db.AutoMigrate(&Partais{})
	return err
}
