package authconf

import "os"

// Defines whether the Middleware should authenticate user
func JwtVerificationDisabled() bool {
	return os.Getenv("PG_JWT_VERIFICATION_DISABLE") != ""
}
