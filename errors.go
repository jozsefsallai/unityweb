package unityweb

// ParseError is an error that occurs while parsing a Unity Web data file.
type ParseError int

const (
	// ErrInvalidMagicHeader is returned when the magic header is invalid. A
	// header is valid if it's a sequence of 16-bytes, consisting of the
	// null-terminated string "UnityWebData1.0"
	ErrInvalidMagicHeader ParseError = iota
)

// Error returns a string representation of the error.
func (e ParseError) Error() string {
	switch e {
	case ErrInvalidMagicHeader:
		return "invalid magic header"
	}
	return "parse error"
}
