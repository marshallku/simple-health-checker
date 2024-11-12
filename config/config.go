package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WebhookURL    string `yaml:"webhook_url"`
	Timeout       int    `yaml:"timeout"`
	Pages         []Page `yaml:"pages"`
	CheckInterval int    `yaml:"check_interval"`
}

type Page struct {
	URL           string   `yaml:"url"`
	Status        int      `yaml:"status"`
	TextToInclude string   `yaml:"text_to_include"`
	Speed         int      `yaml:"speed"`
	Request       *Request `yaml:"request,omitempty"`
}

type Request struct {
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	data := buf.Bytes()

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
