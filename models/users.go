package models

import (
	"github.com/divideandconquer/go-merge/merge"
	"gopkg.in/mgo.v2/bson"
)

// TBLUSERS is the name of the collection for users
var TBLUSERS = "Users"

// Users is a slice of User objects
type Users []User

// User is an object containing user information
type User struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserName   string        `bson:"username" json:"username"`
	FirstName  string        `bson:"firstname" json:"firstname"`
	MiddleName string        `bson:"middlename" json:"middlename"`
	LastName   string        `bson:"lastname" json:"lastname"`
	Email      string        `bson:"email" json:"email"`
	Password   string        `bson:"password" json:"password"`
}

// GetUsers returns the specified user objects
func (d *Datasource) GetUsers(id string) (Users, error) {
	var users Users
	c := d.Database.C(TBLUSERS)

	if id == "" {
		if err := c.Find(nil).All(&users); err != nil {
			return nil, err
		}
	} else {
		if !bson.IsObjectIdHex(id) {
			return nil, ErrObjID
		}
		query := c.FindId(bson.ObjectIdHex(id))
		num, err := query.Count()
		if err != nil {
			return nil, err
		}
		if num <= 0 {
			return nil, ErrNotFound
		}
		if err := query.All(&users); err != nil {
			return nil, err
		}
	}
	return users, nil
}

// InsertUser inserts a new user object
func (d *Datasource) InsertUser(user User) (bson.ObjectId, error) {
	user.ID = bson.NewObjectId()

	if err := d.Database.C(TBLUSERS).Insert(&user); err != nil {
		return "", err
	}
	return user.ID, nil
}

// UpdateUser updates a specified user object
func (d *Datasource) UpdateUser(user User) error {
	var currUser User
	c := d.Database.C(TBLUSERS)

	// get the current user data
	if err := c.FindId(user.ID).One(&currUser); err != nil {
		return err
	}

	// merge the new data with the current
	result := merge.Merge(currUser, user)

	return c.UpdateId(user.ID, result.(*User))
}

// RemoveUser removes a specified user object
func (d *Datasource) RemoveUser(id string) error {
	if !bson.IsObjectIdHex(id) {
		return ErrObjID
	}
	return d.Database.C(TBLUSERS).RemoveId(bson.ObjectIdHex(id))
}
