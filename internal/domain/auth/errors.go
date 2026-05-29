package auth

import "errors"

// EmailAlreadyTaken indicates that an account with the requested email already exists.
var EmailAlreadyTaken = errors.New("email already taken")
