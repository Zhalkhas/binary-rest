package app

import (
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	// Port is the port number to listen HTTP server on
	Port string `koanf:"PORT"`
	// LogLevel configures max log level
	LogLevel string `koanf:"LOG_LEVEL"`
	// InputFileName is the name of the file with indices
	InputFileName string `koanf:"INPUT_FILE_NAME"`
}

// defaultConfig is the default configuration values,
// which are used if not provided by .env file
var defaultConfig = Config{
	Port:          "8080",
	LogLevel:      "INFO",
	InputFileName: "input.txt",
}

// NewConfig creates a new Config instance from .env file,
// merging it with defaultConfig
func NewConfig() (Config, error) {
	k := koanf.New("\n")
	// try to load default config struct
	err := k.Load(structs.Provider(defaultConfig, "koanf"), nil)
	if err != nil {
		slog.Error("error while parsing default config struct", "err", err, "config", defaultConfig)
		return Config{}, err
	}
	// try to load .env file
	err = k.Load(file.Provider(".env"), dotenv.Parser())
	if err != nil {
		slog.Error("error while parsing .env file", "err", err)
		return Config{}, err
	}
	// unmarshal merged loaded configs into Config struct
	var config Config
	err = k.Unmarshal("", &config)
	if err != nil {
		slog.Error("error while unmarshalling config", "err", err)
		return Config{}, err
	}

	logLevel := slog.LevelInfo
	switch strings.ToUpper(config.LogLevel) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "ERROR":
		logLevel = slog.LevelError
	default:
		// logging before setting default logger???
		slog.Info("unknown log level, using INFO", "log_level", config.LogLevel)
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	return config, nil
}
