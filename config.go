package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Gittoken          string   `yaml:"gittoken"`
	Directory         string   `yaml:"directory"`
	Extensions        []string `yaml:"extensions"`
	ExcludeFolders    []string `yaml:"exclude"`
	TodoCommentSymbol string   `yaml:"todocommentsymbol"`
	TodoKeyword       string   `yaml:"todokeyword"`
}

func ParseConfigFile() (*Config, error) {
	var cfg Config
	data, err := os.ReadFile("klaw.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Gittoken == "" {
		return nil, fmt.Errorf("error: gittoken is missing in klaw.yml")
	}
	if cfg.Directory == "" {
		cfg.Directory = "." // defaults to cwd
	}
	if len(cfg.Extensions) == 0 {
		return nil, fmt.Errorf("error: file extensions missing in klaw.yml")
	}
	if cfg.TodoCommentSymbol == "" {
		cfg.TodoCommentSymbol = "//"
	}
	if cfg.TodoKeyword == "" {
		cfg.TodoKeyword = "todo"
	}
	return &cfg, nil
}
