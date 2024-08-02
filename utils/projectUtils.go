package utils

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type ProjectConfig struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Options     struct {
		Refresh string `yaml:"refresh"`
	} `yaml:"options"`
	Runtime  string `yaml:"runtime"`
	Template struct {
		Config struct {
			AWSRegion struct {
				Default     string `yaml:"default"`
			} `yaml:"aws:region"`
			PulumiTags struct {
				AWSRegionDeployed string `yaml:"awsRegionDeployed"`
				ProjectName       string `yaml:"projectName"`
			} `yaml:"pulumi:tags"`
		} `yaml:"config"`
		Description string `yaml:"description"`
		DisplayName string `yaml:"displayName"`
		Metadata struct {
			Cloud struct {
				Provider string `yaml:"cloud:provider"`
			} `yaml:"metadata:cloud"`
		} `yaml:"metadata"`
	} `yaml:"template"`
}

func BuildProjectConfig(projectName, awsRegion, cloudProvider string) (ProjectConfig, error) {
	return ProjectConfig{
		Description: "Cloud project created by Maos Corp.",
		Name:        projectName,
		Runtime:     "go",
		Options: struct {
			Refresh string `yaml:"refresh"`
		}{
			Refresh: "always",
		},
		Template: struct {
			Config struct {
				AWSRegion struct {
					Default     string `yaml:"default"`
				} `yaml:"aws:region"`
				PulumiTags struct {
					AWSRegionDeployed string `yaml:"awsRegionDeployed"`
					ProjectName       string `yaml:"projectName"`
				} `yaml:"pulumi:tags"`
			} `yaml:"config"`
			Description string `yaml:"description"`
			DisplayName string `yaml:"displayName"`
			Metadata    struct {
				Cloud struct {
					Provider string `yaml:"cloud:provider"`
				} `yaml:"metadata:cloud"`
			} `yaml:"metadata"`
		}{
			Config: struct {
				AWSRegion struct {
					Default     string `yaml:"default"`
				} `yaml:"aws:region"`
				PulumiTags struct {
					AWSRegionDeployed string `yaml:"awsRegionDeployed"`
					ProjectName       string `yaml:"projectName"`
				} `yaml:"pulumi:tags"`
			}{
				AWSRegion: struct {
					Default     string `yaml:"default"`
				}{
					Default:     awsRegion,
				},
				PulumiTags: struct {
					AWSRegionDeployed string `yaml:"awsRegionDeployed"`
					ProjectName       string `yaml:"projectName"`
				}{
					AWSRegionDeployed: awsRegion,
					ProjectName:       projectName,
				},
			},
			Description: "A brief description of the Environment name",
			DisplayName: "Environment Name, Prod,Staging",
			Metadata: struct {
				Cloud struct {
					Provider string `yaml:"cloud:provider"`
				} `yaml:"metadata:cloud"`
			}{
				Cloud: struct {
					Provider string `yaml:"cloud:provider"`
				}{
					Provider: cloudProvider,
				},
			},
		},
	}, nil
}


// GetRootDir returns the root directory of the project
func GetRootDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, ".."), nil // Adjusting the path to the project directory
}

func ConvertProjectConfigToAutoConfigMap(projectConfig ProjectConfig) auto.ConfigMap {
	return map[string]auto.ConfigValue{
		"project:description": {Value: "Cloud project created by Maos Corp."},
		"project:name":        {Value: "eks-test-2452"},
		"project:runtime":     {Value: "go"},
		"options:refresh":     {Value: "always"},
		"aws:region":          {Value: "us-east-2"},
		"pulumiTags:awsRegionDeployed": {Value: "us-east-2"},
		"pulumiTags:projectName":       {Value: "eks-test-2452"},
	}

}

func GeneratePulumiYAML(config ProjectConfig, filepath string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("could not marshal YAML: %v", err)
	}
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}
	return nil
}
