package domain

//MySQL structure with mysql db connection parameters
type MySQL struct {
	DbName     string `yaml:"dbName"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Protocol   string `yaml:"protocol"`
	Hostname   string `yaml:"hostname"`
	Port       string `yaml:"port"`
	Parameters string `yaml:"parameters"`
}

//SQLite structure for sqlite db connection
type SQLite struct {
	Filename string `yaml:"filename"`
}

//DbConfig structure describing the database configuration
type DbConfig struct {
	Driver string `yaml:"driver"`
	MySQL  `yaml:"mysql"`
	SQLite `yaml:"sqlite"`
}

//AppConfig application config structure
type AppConfig struct {
	HTTPServer struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"httpServer"`
	Network struct {
		Interface string `yaml:"interface"`
	} `yaml:"network"`
	Database struct {
		Entity DbConfig `yaml:"entity"`
		Log    DbConfig `yaml:"log"`
	} `yaml:"database"`
	Logger struct {
		Level          string `yaml:"level"`
		LogsToDatabase bool   `yaml:"logsToDatabase"`
	} `yaml:"logger"`
}
