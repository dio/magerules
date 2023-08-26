package installable

import "errors"

var ErrNotFound = errors.New("Not found")
var ErrAlreadyInstalled = errors.New("Already installed")
var ErrInvalid = errors.New("Invalid")
