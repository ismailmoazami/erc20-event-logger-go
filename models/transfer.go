package models

import (
	"gorm.io/gorm"
)

type TransferRecord struct {
	gorm.Model
	From  string `gorm:"type:text"`
	To    string `gorm:"type:text"`
	Value string `gorm:"type:text"`
}
