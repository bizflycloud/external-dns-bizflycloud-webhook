package dnsprovider

import (
	"fmt"
	"testing"

	"github.com/bizflycloud/external-dns-bizflycloud-webhook/cmd/webhook/init/configuration"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	cases := []struct {
		name          string
		config        configuration.Config
		env           map[string]string
		providerType  string
		expectedError string
	}{
		{
			name:   "minimal config",
			config: configuration.Config{},
			env: map[string]string{
				"BFC_APP_CREDENTIAL_ID":     "e5d084f79fd5407da705f8df97332090",
				"BFC_APP_CREDENTIAL_SECRET": "Fhkp66ClGHPFTnjHQId1RWrUoG14qMIK8IWT4GxURTNLfY2cAkVmpTwyJ-xYsDlEqrS0lqBSGG8DEGZma6OQgQ"},
			expectedError: "",
		},
		{
			name:   "config with domain filter",
			config: configuration.Config{},
			env: map[string]string{
				"DOMAIN_FILTER":             "vietquocxa.online",
				"BFC_APP_CREDENTIAL_ID":     "e5d084f79fd5407da705f8df97332090",
				"BFC_APP_CREDENTIAL_SECRET": "Fhkp66ClGHPFTnjHQId1RWrUoG14qMIK8IWT4GxURTNLfY2cAkVmpTwyJ-xYsDlEqrS0lqBSGG8DEGZma6OQgQ"},
			expectedError: "",
		},
		{
			name:          "without credential or secret not able to create provider",
			config:        configuration.Config{},
			expectedError: "reading bizflycloudConfig failed: env: environment variable \"BFC_APP_CREDENTIAL_ID\" should not be empty; environment variable \"BFC_APP_CREDENTIAL_SECRET\" should not be empty",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			dnsProvider, err := Init(tc.config)
			fmt.Printf("%+v\n", err)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError, "expecting error")
				return
			}
			assert.NoErrorf(t, err, "error creating provider")
			assert.NotNil(t, dnsProvider)
		})
	}
}
