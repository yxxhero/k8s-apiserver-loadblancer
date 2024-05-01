package util

import (
	"testing"
)

func TestVerifyServiceType(t *testing.T) {
	tests := []struct {
		serviceType string
		expected    bool
	}{
		{"ClusterIP", true},
		{"NodePort", true},
		{"LoadBalancer", true},
		{"ExternalName", true},
		{"InvalidType", false},
	}

	for _, test := range tests {
		result := VerifyServiceType(test.serviceType)
		if result != test.expected {
			t.Errorf("VerifyServiceType(%s) = %t, expected %t", test.serviceType, result, test.expected)
		}
	}
}
