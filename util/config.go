package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

const nameRegexS = `^[a-z0-9][a-z0-9-]*[a-z0-9]$`

var nameRegex *regexp.Regexp

func init() {
	nameRegex = regexp.MustCompile(nameRegexS)
}

type JobConfig struct {
	GQLAdminURL string
	Schedule    string
}

func (m *JobConfig) String() string {
	return fmt.Sprintf(
		`
			JobConfig{
				GQLAdminUrl: %v
				Schedule: %v
			} 	
		`,
		m.GQLAdminURL,
		m.Schedule,
	)
}

type Config struct {
	ExportDestination string                   `mapstructure:"export-destination"`
	ExportUseSSL      bool                     `mapstructure:"export-use-ssl"`
	ExportAccessKey   string                   `mapstructure:"export-access-key"`
	ExportSecretKey   string                   `mapstructure:"export-secret-key"`
	PrometheusPort    uint                     `mapstructure:"prometheus-port"`
	ExportJobsRaw     []map[string]interface{} `mapstructure:"export-jobs"`
	ExportJobs        map[string]*JobConfig
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(filePath string) (*Config, error) {
	viper.SetConfigFile(filePath)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	config.ExportJobs, err = processAdminUrls(config.ExportJobsRaw)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (m *Config) GetExportURL(path string) string {
	url := m.ExportDestination
	if path != "" {
		url = fmt.Sprintf("%v/%v", strings.TrimRight(url, "/"), path)
	}
	if !m.ExportUseSSL {
		url = fmt.Sprintf("%v?secure=false", url)
	}
	return url
}

func processAdminUrls(raw []map[string]interface{}) (map[string]*JobConfig, error) {
	jobs := make(map[string]*JobConfig)
	for _, mapping := range raw {
		name := mapping["name"].(string)
		if !nameRegex.MatchString(name) {
			return nil, fmt.Errorf("invalid job name: %v, valid names are: %v", name, nameRegexS)
		}
		jobs[name] = &JobConfig{
			GQLAdminURL: mapping["gql-admin-url"].(string),
			Schedule:    mapping["schedule"].(string),
		}
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("at least one GQL Admin URL must be specified ")
	}
	return jobs, nil
}

func (m *Config) String() string {
	return fmt.Sprintf(
		`
			Config {
				ExportJobsRaw: %v
				ExportDestination: %v
				ExportAccessKey: %v
				ExportJobs: %v
			}
		`,
		m.ExportJobsRaw,
		m.ExportDestination,
		m.ExportAccessKey,
		m.ExportJobs,
	)
}
