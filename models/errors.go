package models

import "errors"

// ErrNoRows is the error to return when no records were found
var ErrNoRows = errors.New("datasource: no records returned")

// ErrObjID is the error to return when an Object ID is invalid
var ErrObjID = errors.New("invalid object id")

// ErrInvalidCollection is returned when the collection specified is invalid
var ErrInvalidCollection = errors.New("invalid collection")
