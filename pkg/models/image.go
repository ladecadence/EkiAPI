package models

import "time"

type Image struct {
	ID       uint   `gorm:"primaryKey"`
	FileName string `json:"filename" gorm:"unique"`
	DateTime time.Time
	Mission  string
}
