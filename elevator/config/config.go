package config

import (
	"path/filepath"
	"os"
	"github.com/spf13/viper"
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)
	dataset := filepath.Join(exPath, "..", "..", "..", "dataset")
	logs := filepath.Join(exPath, "..", "..", "..", "logs")

	viper.SetEnvPrefix("ELEVATOR")
	viper.AutomaticEnv()

	viper.SetDefault("MaxPassengers", 8)
	viper.SetDefault("ListenAddr", ":8000")

	viper.SetDefault("DatasetDir", dataset)
	viper.SetDefault("LogDir", logs)
}
