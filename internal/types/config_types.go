package types

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Enrichment EnrichmentUrlsConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	DB       string
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

type EnrichmentUrlsConfig struct {
	AgeUrl         string
	GenderUrl      string
	NationalityUrl string
}
