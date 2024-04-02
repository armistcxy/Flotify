package custom_error

type DuplicateUsernameError struct{}

func (e DuplicateUsernameError) Error() string {
	return "this username has been used"
}
