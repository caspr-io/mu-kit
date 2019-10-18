package test

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	gotestenv "gotest.tools/v3/env"
)

func PatchEnvFromFile(envfile string) func() {
	_, err := os.Stat(envfile)
	if err == nil || !os.IsNotExist(err) {
		// if stat returns an error, but it does NOT say the file does not exists,
		// invoke godotenv.Load and let it return the error
		patchEnv, err := godotenv.Read(envfile)
		if err != nil {
			log.Fatal().Err(err).Str("filename", envfile).Msg("Cannot load environment file")
		}
		return PatchEnv(patchEnv)
	} else {
		return func() {}
	}
}

func PatchEnv(envMap map[string]string) func() {
	oldEnv := os.Environ()

	setEnv(envMap)

	return func() {
		os.Clearenv()
		setEnv(gotestenv.ToMap(oldEnv))
	}
}

func setEnv(env map[string]string) {
	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			log.Fatal().Err(err).Str("env", fmt.Sprintf("%s=%s", key, value)).Send()
		}
	}

}
