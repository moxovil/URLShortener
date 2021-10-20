package main

import (
	urlshortener "URLShortener"
	"URLShortener/pkg/handler"
	"URLShortener/pkg/repository"
	"URLShortener/pkg/repository/inmemory"
	"URLShortener/pkg/repository/postgres"
	"URLShortener/pkg/service"
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := initConfig(); err != nil {
		errorLog.Fatalf("ошибка инициализации configs: %s", err.Error())
	}

	var db repository.UrlList
	mode, exists := os.LookupEnv("MEM_MODE")
	if exists && mode == "POSTGRES" {
		var err error
		db, err = postgres.NewPostgresDB(postgres.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		})
		if err != nil {
			errorLog.Fatalf("ошибка инициализации БД: %s\n", err.Error())
		}
	} else {
		db = inmemory.NewInMemoryDB()
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, viper.GetInt("uniquestr.len"), []rune(viper.GetString("uniquestr.chars")))
	handlers := handler.NewHandler(services, errorLog, infoLog)

	srv := new(urlshortener.Server)
	infoLog.Printf("Запуск сервера на :%s\n", viper.GetString("port"))
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.Routes()); err != nil {
			errorLog.Fatal(err)
		}
	}()

	//Graceful shutdown
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc

	if err := srv.Shutdown(context.Background()); err != nil {
		errorLog.Fatalf("error shutdown server: %s\n", err.Error())
	}

	switch db.(type) {
	case *postgres.PostgresDB:
		postgresDB, ok := db.(*postgres.PostgresDB)
		if !ok {
			errorLog.Fatalln("не удалось преобразовать к типу *repository.PostgresDB")
		}
		if err := postgresDB.Close(); err != nil {
			errorLog.Fatalf("ошибка закрытия соединения БД: %s\n", err.Error())
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
