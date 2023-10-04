package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ConfigurationTestSuite struct {
	suite.Suite
	getConfigFile func() (ConfigFile, error)
	tempDirName   string
}

func TestGetConfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigurationTestSuite))
}

func (s *ConfigurationTestSuite) SetupSuite() {
	s.getConfigFile = GetConfigFile

	tempDirName, err := os.MkdirTemp("", "sherpa-config-test")
	s.tempDirName = tempDirName
	s.Require().NoError(err)

	GetConfigFile = func() (ConfigFile, error) {
		return ConfigFile{
			Path: filepath.Join(tempDirName, configPath),
			Name: configName,
			Type: configType,
		}, nil
	}
}

func (s *ConfigurationTestSuite) TearDownSuite() {
	GetConfigFile = s.getConfigFile

	os.RemoveAll(s.tempDirName)
}

func (s *ConfigurationTestSuite) TestGetConfiguration() {

	resetConfigInitialization := func() {
		cfg = nil
		vip = nil
	}

	parseConfiguration := func(configByte []byte) (loadedConfig Configuration, err error) {
		loadedConfigMap := make(map[string]any)
		if err = yaml.Unmarshal(configByte, &loadedConfigMap); err != nil {
			return
		}

		if err = mapstructure.Decode(loadedConfigMap, &loadedConfig); err != nil {
			return
		}

		return
	}

	s.Run("Should panic if configuration is not initialized", func() {
		resetConfigInitialization()

		s.Panics(func() {
			GetConfig()
		})
	})

	s.Run("Should return configuration without panic", func() {
		resetConfigInitialization()

		err := Initialize(false)
		s.Require().NoError(err)

		s.NotPanics(func() {
			GetConfig()
		})
	})

	s.Run("Loads configuration from file if file exists", func() {
		resetConfigInitialization()

		oldGetConfigFile := GetConfigFile
		defer func() {
			GetConfigFile = oldGetConfigFile
		}()

		GetConfigFile = func() (ConfigFile, error) {
			return ConfigFile{
				Path: "./testdata",
				Name: "test-configuration",
				Type: "yml",
			}, nil
		}

		fileContent, err := os.ReadFile("./testdata/test-configuration.yml")
		s.Require().NoError(err)
		loadedConfig, err := parseConfiguration(fileContent)
		s.Require().NoError(err)

		err = Initialize(false)
		s.Require().NoError(err)

		configuration := GetConfig()

		s.Truef(reflect.DeepEqual(loadedConfig, configuration), "Expected: %#v\nActual: %#v", loadedConfig, configuration)
	})

}

func TestGetConfigFile(t *testing.T) {

	t.Run("Returns the configuration file", func(t *testing.T) {
		expectedConfigFile := ConfigFile{
			Path: filepath.Join(os.Getenv("HOME"), configPath),
			Name: configName,
			Type: configType,
		}

		configFile, err := GetConfigFile()
		require.NoError(t, err)

		assert.Equal(t, expectedConfigFile, configFile)
	})
}

type ValidateConfigTestSuite struct {
	suite.Suite
}

func TestValidateConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateConfigTestSuite))
}

func (s *ValidateConfigTestSuite) getNonExistentConfigFile() (ConfigFile, error) {
	tmpDir, err := os.MkdirTemp("", "sherpa-config-test")
	s.Require().NoError(err)
	return ConfigFile{
		Path: tmpDir,
		Name: "non-existent-config",
		Type: "yml",
	}, nil
}

func (s *ValidateConfigTestSuite) getValidConfig() Configuration {
	return Configuration{
		Jira: Jira{
			Auth: JiraAuth{
				Host:        "https://jira.example.com",
				Token:       "token",
				InsecureTLS: false,
			},
			IssueTypes: JiraIssueTypes{
				"bugfix":  {"2"},
				"feature": {"1"},
			},
		},
		Github: Github{
			IssueLabels: GithubIssueLabels{
				"bugfix":  {"bugfix"},
				"feature": {"feature"},
			},
		},
		Branches: Branches{
			Prefixes: BranchesPrefixes{
				"bugfix":  "bugfix",
				"feature": "feat",
			},
		},
	}
}

func (s *ValidateConfigTestSuite) TestConfigurationValidations() {
	s.Run("Default configuration should be valid", func() {
		oldGet := GetConfigFile
		defer func() {
			GetConfigFile = oldGet
		}()
		GetConfigFile = s.getNonExistentConfigFile

		err := Initialize(false)
		s.Require().NoError(err)

		conf := GetConfig()

		err = conf.Validate()

		s.NoError(err)
	})

	s.Run("Should return error if configuration is empty struct", func() {
		config := Configuration{}

		err := config.Validate()

		s.Error(err)
	})

	s.Run("Should return error if jira host is not an url", func() {
		tCfg := s.getValidConfig()
		tCfg.Jira.Auth.Host = "not an url"

		err := tCfg.Validate()

		s.Error(err)
	})

	s.Run("Should not return error if jira host is an url", func() {
		tCfg := s.getValidConfig()
		tCfg.Jira.Auth.Host = "https://jira.example.com"

		err := tCfg.Validate()

		s.NoError(err)
	})

	s.Run("Should not return error if jira host is an http IP", func() {
		tCfg := s.getValidConfig()
		tCfg.Jira.Auth.Host = "http://127.0.0.1/jira"

		err := tCfg.Validate()

		s.NoError(err)
	})

	s.Run("Should return error if jira issue types keys are not valid", func() {
		tCfg := s.getValidConfig()
		tCfg.Jira.IssueTypes = JiraIssueTypes{
			"not a valid type":  {"2"},
			issue_types.Feature: {"1"},
		}

		err := tCfg.Validate()

		s.Error(err)
	})

	s.Run("Should return error if jira issue types IDs are not unique", func() {
		tCfg := s.getValidConfig()
		tCfg.Jira.IssueTypes = JiraIssueTypes{
			issue_types.Bug:     {"2"},
			issue_types.Feature: {"2"},
		}

		err := tCfg.Validate()

		s.Error(err)
	})

	s.Run("Should return error if github issue labels keys are not valid", func() {
		tCfg := s.getValidConfig()
		tCfg.Github.IssueLabels = GithubIssueLabels{
			"not a valid type":  {"bug"},
			issue_types.Feature: {"feature"},
		}

		err := tCfg.Validate()

		s.Error(err)
	})

	s.Run("Should return error if github issue labels are not unique", func() {
		tCfg := s.getValidConfig()
		tCfg.Github.IssueLabels = GithubIssueLabels{
			issue_types.Bug:     {"bug"},
			issue_types.Feature: {"bug"},
		}

		err := tCfg.Validate()

		s.Error(err)
	})

	s.Run("Should return error if branches prefixes keys are not valid", func() {
		tCfg := s.getValidConfig()
		tCfg.Branches.Prefixes = BranchesPrefixes{
			"not a valid type":  "bugfix",
			issue_types.Feature: "feat",
		}

		err := tCfg.Validate()

		s.Error(err)
	})

}
