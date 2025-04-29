package models

type Mission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
}
