package main

import (
	"log"
	"os"

	"github.com/Relax-FM/todo-app-go"
	"github.com/Relax-FM/todo-app-go/pkg/handler"
	"github.com/Relax-FM/todo-app-go/pkg/repository"
	"github.com/Relax-FM/todo-app-go/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Password: 	os.Getenv("DB_PASSWORD"),
		Host: 		viper.GetString("db.host"),
		Port: 		viper.GetString("db.port"),
		Username: 	viper.GetString("db.username"),
		DBName: 	viper.GetString("db.dbname"),
		SSLMode: 	viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err:=srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occurred while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}