package custom_error

type PasswordLengthError struct{}

func (e PasswordLengthError) Error() string {
	return "password length must be between 8 and 64 character"
}

type MismatchError struct{}

func (e MismatchError) Error() string {
	return "email or password does not match"
}

type OldPasswordMismatchError struct{}

func (e OldPasswordMismatchError) Error() string {
	return "old password does not match"
}

type InvalidTokenError struct{}

func (e InvalidTokenError) Error() string {
	return "invalid token"
}

type AccessTokenExpiredError struct{}

func (e AccessTokenExpiredError) Error() string {
	return "access token is expired"
}

type RefreshTokenExpired struct{}

func (e RefreshTokenExpired) Error() string {
	return "refresh token is expired"
}
