package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func Setfromenvfile(envfile string) []string {
	_, err := os.Stat(envfile)
	if err == nil || !os.IsNotExist(err) {
		// if stat returns an error, but it does NOT say the file does not exists,
		// invoke godotenv.Load and let it return the error
		environ := os.Environ()
		err := godotenv.Load(envfile)
		if err != nil {
			log.Fatal().Err(err).Str("filename", envfile).Msg("Cannot load environment file")
		}
		return environ
	} else {
		return nil
	}
}

func Setenviron(envs []string) {
	if envs != nil {
		os.Clearenv()
		for _, env := range envs {
			envNameAndValue := strings.SplitN(env, "=", 2)
			os.Setenv(envNameAndValue[0], envNameAndValue[1])
		}
	}
}
