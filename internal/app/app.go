package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	handlers "people/internal/handler/router"
	"people/internal/repository/enrichment"
	"people/internal/repository/storage"
	"people/internal/types"
	"people/internal/usecase"
)

func Init() {
	logger := logrus.New()

	logger.Info("Starting people")

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	})
	logger.SetLevel(logrus.DebugLevel)

	cfg, err := LoadConfig(".", logger)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	serverHost := cfg.Server.Host + ":" + cfg.Server.Port

	db := cfg.Database.DB
	user := cfg.Database.User
	password := cfg.Database.Password
	dbHost := cfg.Database.Host + ":" + cfg.Database.Port
	dbName := cfg.Database.DBName

	dbUrl := fmt.Sprintf("%s://%s:%s@%s/%s", db, user, password, dbHost, dbName)
	log.Printf("Connecting to %s", dbUrl)

	ageUrl := cfg.Enrichment.AgeUrl
	genderUrl := cfg.Enrichment.GenderUrl
	nationalityUrl := cfg.Enrichment.NationalityUrl

	ctx := context.Background()

	store, err := storage.New(ctx, dbUrl, logger)
	if err != nil {
		logger.Fatalf("Failed start storage. Error: %v", err)
	}

	enrichments, err := enrichment.New(ageUrl, genderUrl, nationalityUrl, logger)
	if err != nil {
		logger.Fatalf("Failed start enrichment. Error: %v", err)
	}

	useCase := usecase.New(store, enrichments, logger)

	server := handlers.New(useCase, logger)

	router := handlers.Router(server)

	if err = router.Run(serverHost); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

func LoadConfig(path string, log *logrus.Logger) (types.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Error("Failed loading config")
		return types.Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg types.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.WithError(err).Error("Failed unmarshalling config")
		return types.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
