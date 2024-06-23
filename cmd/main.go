package main

import (
	"TalkBoard/pkg/handler"
	"TalkBoard/pkg/repository"
	"TalkBoard/pkg/service"
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	//storageType := flag.String("storage", "memory", "Type of storage: memory or postgres")
	//flag.Parse()
	storageType := "memory"

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	var repos *repository.Repository
	//*storageType
	switch storageType {
	case "postgres":
		fmt.Println("Starting application with PostgreSQL storage")
		repos = repository.NewRepository(db, false)
	case "memory":
		fmt.Println("Starting application with in-memory storage")
		repos = repository.NewRepository(db, true)
	default:
		fmt.Println("Unknown storage type, starting with default PostgreSQL storage")
		repos = repository.NewRepository(db, false)
	}

	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	//middlewares := handler.NewMiddleware(services)

	//http.Handle("/graphql", middlewares.AuthMiddleware(handlers.InitGraphQL()))
	http.Handle("/graphql", handlers.InitGraphQL())
	http.Handle("/playground", playground.Handler("GraphQL", "/graphql"))

	port := viper.GetString("port")
	logrus.Infof("connect to http://localhost:%s/playground for GraphQL playground", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
