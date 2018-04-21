package models

import (
	"errors"

	"gopkg.in/mgo.v2"
)

// ErrNotFound is the error to return when no records were found
var ErrNotFound = mgo.ErrNotFound

// ErrObjID is the error to return when an Object ID is invalid
var ErrObjID = errors.New("invalid object id")

// ErrInvalidCollection is returned when the collection specified is invalid
var ErrInvalidCollection = errors.New("invalid collection")
