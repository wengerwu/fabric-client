package parse

var DB DBConfig

type DBConfig struct {
	MasterDB DBYamlConfig `yaml:"master"`
	Slave	DBYamlConfig	`yaml:"slave"`
}

type DBYamlConfig struct {
	Dialect	string `yaml:"dialect"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Database string `yaml:"database"`
	Charset string `yaml:"charset"`
	ShowSql bool `yaml:"showSql"`
	LogLevel string `yaml:"logLevel"`
	MaxIdleConns int `yaml:"maxIdleConns"`
	MaxOpenConns int `yaml:"maxOpenConns"`
}

