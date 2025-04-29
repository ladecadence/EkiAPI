package models

type Image struct {
	ID       uint   `gorm:"primaryKey"`
	FileName string `json:"filename" gorm:"unique"`
	Mission  string
}
