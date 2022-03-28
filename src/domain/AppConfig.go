package domain

//AppConfig application config structure
type AppConfig struct {
	HTTPServer struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"httpServer"`
	Database struct {
		EntityConnectionString string `yaml:"entityConnectionString"`
		EntityDbName           string `yaml:"entityDbName"`
		EntityDbParams         string `yaml:"entityDbParams"`
		LogConnectionString    string `yaml:"logConnectionString"`
		LogDbName              string `yaml:"logDbName"`
		LogDbParams            string `yaml:"logDbParams"`
	} `yaml:"database"`
	Logger struct {
		Level          string `yaml:"level"`
		LogsToDatabase bool   `yaml:"logsToDatabase"`
	} `yaml:"logger"`
}
