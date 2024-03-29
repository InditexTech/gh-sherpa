package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/InditexTech/gh-sherpa/internal/interactive"
	"github.com/InditexTech/gh-sherpa/internal/logging"
	"github.com/InditexTech/gh-sherpa/pkg/metadata"
	"github.com/InditexTech/gh-sherpa/pkg/validator"
	"github.com/spf13/viper"
)

const (
	configPath = ".config/sherpa"
	configName = "config"
	configType = "yml"
)

//go:embed default-config.yml
var defaultConfigBuff []byte

var (
	cfg *Configuration
	vip *viper.Viper
)

type Configuration struct {
	Jira     Jira
	Github   Github `validate:"required"`
	Branches Branches
}

// Validates the configuration
func (c Configuration) Validate() error {
	if err := validator.Struct(c); err != nil {
		return fmt.Errorf("configuration is invalid:\n%w", err)
	}

	return nil
}

// GetConfig returns the configuration
func GetConfig() Configuration {
	if cfg == nil {
		panic("Configuration not initialized")
	}

	return *cfg
}

type ConfigFile struct {
	Path string
	Name string
	Type string
}

func (cf ConfigFile) getFilePath() string {
	return filepath.Join(cf.Path, fmt.Sprintf("%s.%s", cf.Name, cf.Type))
}

var GetConfigFile = func() (cfgFile ConfigFile, err error) {
	var homeDir string
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return
	}

	cfgFile = ConfigFile{
		Path: filepath.Join(homeDir, configPath),
		Name: configName,
		Type: configType,
	}

	return
}

// Initialize initializes the configuration
func Initialize(isInteractive bool) error {
	vip = viper.New()

	// Load default configuration into viper
	vip.SetConfigType(configType)
	if err := vip.MergeConfig(bytes.NewBuffer(defaultConfigBuff)); err != nil {
		return err
	}

	cfgFile, err := GetConfigFile()
	if err != nil {
		return err
	}

	logging.Debugf("Reading config file from %s", cfgFile.getFilePath())
	vip.AddConfigPath(cfgFile.Path)
	vip.SetConfigName(cfgFile.Name)
	vip.SetConfigType(cfgFile.Type)

	if err := vip.MergeInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			if err := generateConfigurationFile(cfgFile, isInteractive); err != nil {
				return err
			}
		default:
			return err
		}
	}

	// Unmarshal configuration into target struct
	if err := vip.Unmarshal(&cfg); err != nil {
		return err
	}

	return cfg.Validate()
}

func generateConfigurationFile(cfgFile ConfigFile, isInteractive bool) error {
	// If the file doesn't exist, create it
	logging.PrintWarn(fmt.Sprintf("Config file not found, generating a new configuration in %s", cfgFile.getFilePath()))
	if isInteractive {
		if err := askJiraConfiguration(); err != nil {
			return err
		}
	} else {
		logging.PrintInfo("Skipping interactive configuration")
	}

	// Write configuration into file
	if err := writeConfigurationFile(cfgFile); err != nil {
		return err
	}

	return nil
}

func askJiraConfiguration() error {
	shouldConfigureJira, err := interactive.AskUserForConfirmation("Will you use Jira issues?", false)
	if err != nil {
		return err
	}
	if shouldConfigureJira {
		if err := configureJira(); err != nil {
			return err
		}
	} else {
		logging.PrintInfo("Skipping Jira configuration. You can configure it later in the configuration file")
	}

	return nil
}

func writeConfigurationFile(cfgFile ConfigFile) error {
	if err := os.MkdirAll(cfgFile.Path, os.ModePerm); err != nil {
		return err
	}
	filePath := cfgFile.getFilePath()
	vip.SetConfigFile(filePath)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := vip.Unmarshal(&cfg); err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	configFileTemplateData := configFileTemplateData{
		Metadata: MetadataConfiguration{
			Version:     metadata.Version,
			GeneratedAt: time.Now(),
		},
		JiraData: JiraTemplateConfiguration{
			Jira: cfg.Jira,
		},
	}
	if err := writeTemplatedConfigFile(f, configFileTemplateData); err != nil {
		return err
	}

	return nil
}
