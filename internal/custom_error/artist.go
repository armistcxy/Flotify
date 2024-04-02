package custom_error

type NonExistArtistError struct{}

func (e NonExistArtistError) Error() string {
	return "non exist artist record in database"
}
