package config

type Config struct {
	Env             string
	BloomFilterPath string
	Server          Server
	Database        Database
	Paseto          Paseto
	Seeder          Seeder
}

type Server struct {
	Host string
	Port string
}

type Database struct {
	Host string
	Port string
	Name string
	User string
	Pass string
	Tz   string
}

type Paseto struct {
	SecretKey   string
	TokenTTLMin int
}

type Seeder struct {
	Admin Admin
}

type Admin struct {
	Name     string
	Email    string
	Password string
}
