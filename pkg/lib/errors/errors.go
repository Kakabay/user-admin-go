package errors

// auth
const (
	InvalidRequestFormat    = "Invalid request format"
	AdminNotFound           = "Admin not found"
	InvalidCredentials      = "Invalid credentials"
	RefreshTokenNotProvided = "Refresh token not provided"
	InvalidRefreshToken     = "Invalid refresh token"
)

// user
const (
	InternalServerError = "Internal server error"
	InvalidID           = "Invalid ID"
	InvalidRequestBody  = "Invalid request body"
)

// middleware
const (
	AuthorizationTokenNotProvided = "Authorization token not provided"
	RoleNotFoundInTokenClaims     = "Role not found in token claims"
	InsufficientPermission        = "Insufficient permissions"
	TokenClaimsNotFound           = "Token claims not found"
)
