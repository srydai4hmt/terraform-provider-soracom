package soracom

import "github.com/soracom/soracom-sdk-go"

// Config type of SORACOM Config
type Config struct {
	AuthKeyID     string
	AuthKeySecret string
}

// NewClient returns new APIClient for SORACOM
func (c *Config) NewClient() (interface{}, error) {
	client := soracom.NewAPIClient(nil)
	err := client.AuthWithAuthKey(c.AuthKeyID, c.AuthKeySecret)
	if err != nil {
		return nil, err
	}

	return client, nil
}
