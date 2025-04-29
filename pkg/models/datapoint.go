package models

import "time"

type Datapoint struct {
	ID        uint      `gorm:"primaryKey"`
	MissionId string    `json:"id"`
	Msg       string    `json:"msg"`
	Lat       float64   `json:"lat"`
	NS        string    `json:"ns"`
	Lon       float64   `json:"lon"`
	EW        string    `json:"ew"`
	Alt       float64   `json:"alt"`
	Hdg       float64   `json:"hdg"`
	Spd       float64   `json:"spd"`
	Sats      int       `json:"sats"`
	Vbat      float64   `json:"vbat"`
	Baro      float64   `json:"baro"`
	Tin       float64   `json:"tin"`
	Tout      float64   `json:"tout"`
	Arate     float64   `json:"arate"`
	Date      string    `json:"date"`
	Time      string    `json:"time"`
	Sep       string    `json:"sep"`
	DateTime  time.Time `json:"datetime"`
	Hpwr      bool      `json:"hpwr"`
}
