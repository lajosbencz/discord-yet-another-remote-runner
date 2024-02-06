package discordyetanotherremoterunner

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConfigCommand struct {
	Cmd  string
	Args []string
}

func (c *ConfigCommand) String() string {
	return strings.TrimSpace(c.Cmd + " " + strings.Join(c.Args, " "))
}

type ConfigServer struct {
	Name  string
	Start ConfigCommand
	Stop  ConfigCommand
}

type Config struct {
	Guild    string
	Commands map[string]ConfigCommand
	Servers  map[string]ConfigServer
}

func ReadConfig(filePath string) (Config, error) {
	config := Config{}

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
