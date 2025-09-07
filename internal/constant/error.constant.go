package constant

import "errors"

// ErrNotFound digunakan saat resource tidak ditemukan
var ErrNotFound = errors.New("requested resource not found")

// ErrUnauthorized digunakan saat kredensial/token tidak valid
var ErrUnauthorized = errors.New("invalid credentials or token")

// ErrForbidden digunakan saat user tidak punya permission
var ErrForbidden = errors.New("user does not have the required permission")

// ErrConflict digunakan saat resource sudah ada
var ErrConflict = errors.New("resource already exists")

// ErrInvalidInput digunakan saat input tidak valid
var ErrInvalidInput = errors.New("invalid input provided")
