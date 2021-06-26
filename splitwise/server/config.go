package server

import (
	"fmt"
	"os"

	"github.com/olebedev/config"
)

// Opts options for creating application
type Opts struct {
	Verbose bool
}

// GetOptsFromConfig gets opts from config
func GetOptsFromConfig(filename string) (Opts, error) {
	opts := Opts{}

	cfg := &config.Config{
		Root: make(map[string]interface{}),
	}
	for key, value := range defaults {
		cfg.Set(key, value)
	}

	if fileExists(filename) {
		cfgf, err := config.ParseYamlFile(filename)
		if err != nil {
			return opts, err
		}

		cfg, err = cfg.Extend(cfgf)
		if err != nil {
			return opts, fmt.Errorf("unable to extend cfg : %v", err)
		}
	}

	cfg.Env()
	fmt.Println(config.RenderYaml(cfg.Root))

	opts.Verbose = cfg.UBool("VERBOSE")
	return opts, nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
