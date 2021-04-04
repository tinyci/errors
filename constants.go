package errors

/*
 * These constants are mostly related to tinyCI and things that consume it.
 */

var (
	// ErrNotFound is a generic error for things that don't exist.
	ErrNotFound = New("not found")

	// ErrInvalidAuth is for authentication events that do not succeed.
	ErrInvalidAuth = New("invalid authentication")

	// ErrRunCanceled is thrown to github when the run is canceled.
	ErrRunCanceled = New("run canceled by user intervention")
)
