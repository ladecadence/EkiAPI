package database

import (
	//"fmt"

	"errors"
	"fmt"

	"github.com/ladecadence/EkiAPI/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type SQLite struct {
	db *gorm.DB
}

func (s *SQLite) Open(fileName string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	s.db = database
	return s.db, nil
}

func (s *SQLite) Init() error {
	err := s.db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	err = s.db.AutoMigrate(&models.Mission{})
	if err != nil {
		return err
	}

	// err = s.db.AutoMigrate(&models.Datapoint{})
	// if err != nil {
	// 	return err
	// }

	// err = s.db.AutoMigrate(&models.Problem{})
	// if err != nil {
	// 	return err
	// }
	// return s.db.AutoMigrate(&models.Wall{})
	return nil
}

func (s *SQLite) UpsertUser(u models.User) error {
	result := s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&u)
	return result.Error
}

func (s *SQLite) DeleteUser(models.User) error {
	// TODO
	return nil
}

func (s *SQLite) GetUsers() ([]models.User, error) {
	var users []models.User
	result := s.db.Find(&users)
	return users, result.Error
}

func (s *SQLite) GetUser(name string) (models.User, error) {
	var user models.User
	result := s.db.Where("name=?", name).First(&user)
	return user, result.Error
}

func (s *SQLite) InsertMission(m models.Mission) error {
	result := s.db.Create(&m)
	return result.Error
}

func (s *SQLite) GetMissions() ([]models.Mission, error) {
	var missions []models.Mission
	result := s.db.Find(&missions)
	return missions, result.Error
}

func (s *SQLite) GetMission(name string) (models.Mission, error) {
	var mission models.Mission
	result := s.db.Where("name=?", name).First(&mission)
	return mission, result.Error
}

func (s *SQLite) MissionExists(missionName string) (bool, error) {
	// check if mission table exists
	missions, err := s.GetMissions()
	if err != nil {
		return false, err
	}
	found := false
	for _, mission := range missions {
		if mission.Name == missionName {
			found = true
			break
		}
	}
	return found, nil
}

func (s *SQLite) CreateMissionTable(missionName string) error {
	// create a table for the datapoint model and rename it to the mission name
	err := s.db.Migrator().CreateTable(&models.Datapoint{})
	if err != nil {
		return err
	}
	err = s.db.AutoMigrate(&models.Datapoint{})
	if err != nil {
		return err
	}
	err = s.db.Migrator().RenameTable("datapoints", missionName)
	return err
}

func (s *SQLite) GetMissionData(missionName string) ([]models.Datapoint, error) {
	// check if mission table exists
	found, err := s.MissionExists(missionName)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("No such mission")
	}
	// ok
	data := []models.Datapoint{}
	result := s.db.Raw("SELECT * FROM " + missionName + " ORDER BY id;").Scan(&data)
	return data, result.Error
}

func (s *SQLite) InsertDatapoint(data models.Datapoint) error {
	fmt.Printf("ID: %s\n", data.MissionId)
	found, err := s.MissionExists(data.MissionId)
	if err != nil {
		return err
	}
	// if table not found, create it
	if !found {
		err := s.CreateMissionTable(data.MissionId)
		if err != nil {
			return err
		}
		// and create mission
		mission := models.Mission{Name: data.MissionId}
		err = s.InsertMission(mission)
		if err != nil {
			return err
		}
	}

	// ok, insert data
	sql := "INSERT INTO " +
		data.MissionId +
		" (mission_id, msg, lat, ns, lon, ew, alt, hdg, spd, sats, vbat, baro, tin, tout, arate, date, time, sep, date_time, hpwr)" +
		" VALUES ( " +
		"'" + data.MissionId + "', " +
		"'" + data.Msg + "', " +
		fmt.Sprintf("%07.2f", data.Lat) + ", " +
		"'" + data.NS + "', " +
		fmt.Sprintf("%08.2f", data.Lon) + ", " +
		"'" + data.EW + "', " +
		fmt.Sprintf("%.1f", data.Alt) + ", " +
		fmt.Sprintf("%.1f", data.Hdg) + ", " +
		fmt.Sprintf("%.1f", data.Spd) + ", " +
		fmt.Sprintf("%d", data.Sats) + ", " +
		fmt.Sprintf("%.1f", data.Vbat) + ", " +
		fmt.Sprintf("%.1f", data.Baro) + ", " +
		fmt.Sprintf("%.1f", data.Tin) + ", " +
		fmt.Sprintf("%.1f", data.Tout) + ", " +
		fmt.Sprintf("%.1f", data.Arate) + ", " +
		"'" + data.Date + "', " +
		"'" + data.Time + "', " +
		"'" + data.Sep + "', " +
		"'" + data.DateTime.Format("2006-01-02 15:04:05.00") + "', " +
		func() string {
			if data.Hpwr {
				return "1"
			} else {
				return "0"
			}
		}() + ");"
	result := s.db.Exec(sql)
	return result.Error
}
