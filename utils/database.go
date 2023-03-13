package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var once sync.Once

var instance *gorm.DB

func NewDatabase() *gorm.DB {
	once.Do(func() {
		instance = initializeDB()
	})
	return instance
}

func CloseDatabase() {
	database := NewDatabase()
	db, err := database.DB()
	PanicError("Error Getting Database", err)
	err = db.Close()
	PanicError("Error Closing Database", err)
}

type database struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
}

func loadVariables() (db database) {
	//err := godotenv.Load()
	//PanicError("Error loading .env file", err)
	myEnv, err := godotenv.Read()
	PanicError("Error Reading .env file", err)
	db.Host = myEnv["HOST"]
	db.Port = myEnv["DBPORT"]
	db.User = myEnv["USER"]
	db.Name = myEnv["NAME"]
	db.Password = myEnv["PASSWORD"]
	return
}

func initializeDB() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	dataBase := loadVariables()
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", dataBase.Host, dataBase.User, dataBase.Name, dataBase.Password, dataBase.Port)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURI,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	PanicError("Cannot Connect to Database", err)
	return db
}
