package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"sync"
)

// Env values
type Env struct {
	Server Server
	Logger Logger
	Doc    Doc
}

// Server config
type Server struct {
	AppVersion string
	Port       string
}

// Logger config
type Logger struct {
	Enabled  bool
	Level    string
	FilePath string
}

// Doc - swagger information
type Doc struct {
	Title       string
	Description string
	Version     string
	BasePath    string
}

var (
	env  *Env
	once sync.Once
)

// GetEnv returns env values
func GetEnv() *Env {

	once.Do(func() {
		viper.AutomaticEnv()
		if err := godotenv.Load("config/.env"); err != nil {
			//log.Println(err)
		}
		env = new(Env)
		env.Server.AppVersion = viper.GetString("SERVER_VERSION")
		env.Server.Port = viper.GetString("SERVER_PORT")

		env.Logger.Level = viper.GetString("LOG_LEVEL")
		env.Logger.Enabled = viper.GetBool("LOG_ENABLED")

		env.Doc.Title = viper.GetString("DOC_TITLE")
		env.Doc.Description = viper.GetString("DOC_DESCRIPTION")
		env.Doc.Version = viper.GetString("DOC_VERSION")
		env.Doc.BasePath = viper.GetString("DOC_BASE_PATH")
	})
	return env
}
