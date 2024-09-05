package mirror

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestConvertSVCPorts(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    []v1.ServicePort
		expected []v1.ServicePort
	}{
		{
			name:     "No ports to convert",
			input:    []v1.ServicePort{},
			expected: []v1.ServicePort{},
		},
		{
			name: "Convert a single port",
			input: []v1.ServicePort{
				{
					Name:       "https",
					Port:       443,
					TargetPort: intstr.FromInt(443),
					Protocol:   v1.ProtocolTCP,
				},
			},
			expected: []v1.ServicePort{
				{
					Name:       "https",
					Port:       6443,
					TargetPort: intstr.FromInt(443),
					Protocol:   v1.ProtocolTCP,
				},
			},
		},
		{
			name: "Convert multiple ports",
			input: []v1.ServicePort{
				{
					Name:       "https",
					Port:       443,
					TargetPort: intstr.FromInt(443),
					Protocol:   v1.ProtocolTCP,
				},
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   v1.ProtocolTCP,
				},
			},
			expected: []v1.ServicePort{
				{
					Name:       "https",
					Port:       6443,
					TargetPort: intstr.FromInt(443),
					Protocol:   v1.ProtocolTCP,
				},
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   v1.ProtocolTCP,
				},
			},
		},
		{
			name: "No port named 'https'",
			input: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   v1.ProtocolTCP,
				},
			},
			expected: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   v1.ProtocolTCP,
				},
			},
		},
		{
			name: "Port named 'https' but not on port 443",
			input: []v1.ServicePort{
				{
					Name:       "https",
					Port:       8443,
					TargetPort: intstr.FromInt(8443),
					Protocol:   v1.ProtocolTCP,
				},
			},
			expected: []v1.ServicePort{
				{
					Name:       "https",
					Port:       8443,
					TargetPort: intstr.FromInt(8443),
					Protocol:   v1.ProtocolTCP,
				},
			},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := convertSVCPorts(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}
