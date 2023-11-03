package main

import (
	"context"
	"golangToDo/todo-app"
	"golangToDo/todo-app/pkg/handlers"
	"golangToDo/todo-app/pkg/repository"
	"golangToDo/todo-app/pkg/service"
	"os/signal"
	"syscall"

	// "log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка инициализации конфига: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка загрузки env переменной: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{Host: viper.GetString("db.host"), Port: viper.GetString("db.port"), Username: viper.GetString("db.username"), Password: os.Getenv("DB_PASSWORD"), DBName: viper.GetString("db.dbname"), SSLMode: viper.GetString("db.sslmode")})
	if err != nil {
		logrus.Fatalf("Ошибка инициализации БД: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewServices(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)

	go func() {

		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Ошибка: %s", err.Error())
		}
	}()

	logrus.Print("Started app")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Shutdown app")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
