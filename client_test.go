package ovoclient

import(
	"testing"
)

func TestConfigurationLoad(t *testing.T) {
	client := NewClientFromConfigPath("config.json")
	if client == nil {
		t.Fail()
	}
}