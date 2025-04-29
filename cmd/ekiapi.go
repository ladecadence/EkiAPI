package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"

	// "os"
	"net/http"

	"github.com/ladecadence/EkiAPI/pkg/color"
	"github.com/ladecadence/EkiAPI/pkg/config"
	"github.com/ladecadence/EkiAPI/pkg/database"
	"github.com/ladecadence/EkiAPI/pkg/models"
	"github.com/ladecadence/EkiAPI/pkg/routes"
)

func initData(db database.Database) {
	user := models.User{
		Name:     "testuser",
		Password: fmt.Sprintf("%x", sha256.Sum256([]byte("testpassword"))),
		Email:    "test@email.com",
		Role:     models.UserRoleUser,
	}
	err := db.UpsertUser(user)
	if err != nil {
		panic("Error Inserting initial data: " + err.Error())
	}
	user = models.User{
		Name:     "testadmin",
		Password: fmt.Sprintf("%x", sha256.Sum256([]byte("testAdminpassword"))),
		Email:    "admin@email.com",
		Role:     models.UserRoleAdmin,
	}
	err = db.UpsertUser(user)
	if err != nil {
		panic("Error Inserting initial data: " + err.Error())
	}
}

func main() {
	// flags
	testDataFlag := flag.Bool("testdata", false, "load test data into database")
	configFileFlag := flag.String("conf", "config.toml", "config file path")
	flag.Parse()

	// read configuration
	config := config.Config{ConfFile: *configFileFlag}
	config.GetConfig()

	// open and init database
	database := &database.SQLite{}
	_, err := database.Open(config.Database)
	if err != nil {
		panic("Error opening DB: " + err.Error())
	}
	err = database.Init()
	if err != nil {
		panic("Error initializing DB: " + err.Error())
	}

	// load testa data in database if -testdata arg
	if *testDataFlag {
		initData(database)
	}

	// create server mux and routes
	mux := http.NewServeMux()
	routes.RegisterRoutes(database, config, mux)

	const logo = `
	
	███████╗██╗  ██╗██╗ █████╗ ██████╗ ██╗
	██╔════╝██║ ██╔╝██║██╔══██╗██╔══██╗██║
	█████╗  █████╔╝ ██║███████║██████╔╝██║
	██╔══╝  ██╔═██╗ ██║██╔══██║██╔═══╝ ██║
	███████╗██║  ██╗██║██║  ██║██║     ██║
	╚══════╝╚═╝  ╚═╝╚═╝╚═╝  ╚═╝╚═╝     ╚═╝						  
																									 
	`
	fmt.Fprintf(os.Stderr, "%s", color.Purple+logo+color.Reset)
	log.Printf("📡 "+color.Green+"EkiAPI version "+color.Purple+"%s"+color.Green+" listening on port "+color.Yellow+"%d"+color.Reset, config.Version, config.Port)

	// launch server
	addr := fmt.Sprintf("%s:%d", config.Addr, config.Port)
	log.Fatal(http.ListenAndServe(addr, mux))

}
