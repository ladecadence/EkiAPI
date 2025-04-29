package database

import (
	"github.com/ladecadence/EkiAPI/pkg/models"
	"gorm.io/gorm"
)

type Database interface {
	Open(string) (*gorm.DB, error)
	Init() error
	UpsertUser(models.User) error
	DeleteUser(models.User) error
	GetUsers() ([]models.User, error)
	GetUser(string) (models.User, error)
	GetMissions() ([]models.Mission, error)
	GetMission(string) (models.Mission, error)
	MissionExists(string) (bool, error)
	InsertMission(models.Mission) error
	CreateMissionTable(string) error
	GetMissionData(string) ([]models.Datapoint, error)
	InsertDatapoint(models.Datapoint) error
}
