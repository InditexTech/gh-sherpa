package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

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
