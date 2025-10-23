package main

import (
	"fmt"
	"os"
	"strings"
	"log/slog"

	flag "github.com/spf13/pflag"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	
)

// initLogger initializes logger instance.
func initLogger(ko *koanf.Koanf) *slog.Logger {
	// Default log level
	level := slog.LevelInfo

	// Check app.log value in config
	switch ko.String("app.log") {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	// Choose output format (text, can switch to JSON if you like)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // like EnableCaller
	})

	logger := slog.New(handler)

	return logger
}

// initConfig loads config to `ko` object.
func initConfig() (*koanf.Koanf, error) {
	var (
		ko = koanf.New(".")
		f  = flag.NewFlagSet("front", flag.ContinueOnError)
	)

	// Configure Flags.
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	// Register `--config` flag.
	cfgPath := f.String("config", "config.sample.toml", "Path to a config file to load.")

	// Parse and Load Flags.
	err := f.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	err = ko.Load(file.Provider(*cfgPath), toml.Parser())
	if err != nil {
		return nil, err
	}
	err = ko.Load(env.Provider("KVSTORE_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "KVSTORE_")), "__", ".", -1)
	}), nil)
	if err != nil {
		return nil, err
	}
	return ko, nil
}
