package exporters

// getEnvironmentName returns a formatted environment name
func getEnvironmentName(env string) string {
	switch env {
	case "dev":
		return "Development"
	case "staging":
		return "Staging"
	case "prod":
		return "Production"
	default:
		return "Development"
	}
}

// getBaseURL returns the base URL for the environment
func getBaseURL(env string) string {
	switch env {
	case "dev":
		return "http://localhost:3000"
	case "staging":
		return "https://staging-api.example.com"
	case "prod":
		return "https://api.example.com"
	default:
		return "{{base_url}}"
	}
}

