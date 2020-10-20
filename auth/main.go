package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/riandigitalent/microservice1/auth/config"
	"github.com/riandigitalent/microservice1/auth/database"
	"github.com/riandigitalent/microservice1/auth/handler"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(cfg)
	}
	router := mux.NewRouter()

	db, err := initDB(cfg.Database)
	authHandler := handler.Auth{Db: db}
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("sukses konek DB")
	}
	router.Handle("/auth/validate", http.HandlerFunc(authHandler.ValidateAuth))
	router.Handle("/auth/signup", http.HandlerFunc(authHandler.SignUp))
	router.Handle("/auth/login", http.HandlerFunc(authHandler.Login))

	fmt.Printf("Auth service listen on :8001")
	log.Panic(http.ListenAndServe(":8001", router))
}

func getConfig() (config.Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.SetConfigName("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return config.Config{}, err
	}

	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return config.Config{}, err
	}
	return cfg, nil
}

func initDB(cfg config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.Config)
	log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&database.Auth{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
