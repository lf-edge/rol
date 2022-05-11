package domain

//DbConfig structure describing the database configuration
type DbConfig struct {
	DbName     string `yaml:"dbName"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Protocol   string `yaml:"protocol"`
	Hostname   string `yaml:"hostname"`
	Port       string `yaml:"port"`
	Parameters string `yaml:"parameters"`
}

//AppConfig application config structure
type AppConfig struct {
	HTTPServer struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"httpServer"`
	Database struct {
		Entity DbConfig `yaml:"entity"`
		Log    DbConfig `yaml:"log"`
	} `yaml:"database"`
	Logger struct {
		Level          string `yaml:"level"`
		LogsToDatabase bool   `yaml:"logsToDatabase"`
	} `yaml:"logger"`
}
