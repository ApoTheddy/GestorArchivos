package models

import "gorm.io/gorm"

type FileModel struct {
	gorm.Model

	Filename   string `gorm:"type:varchar(255);not null" json:"filename"`
	FilePath   string `gorm:"type:text;not null" json:"filepath"`
	UniqueName string `gorm:"type:text;not null" json:"uniquename"`
	Dir        string `gor:"type:varchar(150);not null" json:"dirname"`
}

func (FileModel) TableName() string {
	return "File"
}
