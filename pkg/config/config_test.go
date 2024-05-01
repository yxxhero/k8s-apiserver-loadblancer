package config

import (
	"fmt"
	"testing"
)

func TestConfig_Verify(t *testing.T) {
	tests := []struct {
		name        string
		serviceType string
		expectedErr error
	}{
		{
			name:        "ValidServiceType",
			serviceType: "ClusterIP",
			expectedErr: nil,
		},
		{
			name:        "InvalidServiceType",
			serviceType: "InvalidType",
			expectedErr: fmt.Errorf("invalid service type: InvalidType"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := &Config{
				ServiceType: test.serviceType,
			}

			err := config.Verify()
			if err != test.expectedErr {
				t.Errorf("Verify() error = %v, expected %v", err, test.expectedErr)
			}
		})
	}
}
