package mailer

// Custom error type for testing
type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}

// Helper function to check if string contains substring
func contains(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
