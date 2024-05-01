package util

// VerifyServiceType verifies the service type
func VerifyServiceType(serviceType string) bool {
	switch serviceType {
	case "ClusterIP", "NodePort", "LoadBalancer", "ExternalName":
		return true
	default:
		return false
	}
}
