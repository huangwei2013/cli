package utils

import (
	"github.com/ghodss/yaml"
	projectClient "github.com/rancher/types/client/project/v3"
	"strings"
)

func PipelineConfigFromYaml(content []byte) (*projectClient.PipelineConfig, error) {

	var out projectClient.PipelineConfig
	err := yaml.Unmarshal(content, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}


func SplitOnColon(s string) []string {
	return strings.Split(s, ":")
}
