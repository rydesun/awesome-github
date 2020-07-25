package config

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	config := Config{
		AccessToken: "123456",
		StartPoint: StartPoint{
			Path: "tester/awesome-test",
		},
		Github: Github{
			HTMLHost: "https://github.com",
			ApiHost:  "https://api.github.com",
		},
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("expect no error, but got a error: %v", err)
	}

	config.AccessToken = ""
	err = config.Validate()
	if err == nil {
		t.Errorf("expect a error")
	}
	config.Github.ApiHost = "invalid"
	err = config.Validate()
	if err == nil {
		t.Errorf("expect a error")
	}
}
