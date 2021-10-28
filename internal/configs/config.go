package configs

type Configs struct {
	Port  string `json:"port"`
	DB *DB `json:"db"`
}

type DB struct {
	Host string `json:"host"`
	Password string `json:"password"`
	User string `json:"user"`
	Port string `json:"port"`
	Name string `json:"name"`
}

func NewConfig() *Configs {
	return &Configs{
		DB: &DB{},
	}
}
