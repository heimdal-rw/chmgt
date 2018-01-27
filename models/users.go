package models

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

var ErrNoRows = errors.New("datasource: no records returned")

type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	FirstName string        `json:"firstname"`
	LastName  string        `json:"lastname"`
	UserName  string        `json:"username"`
	Email     string        `json:"email"`
	Password  string        `json:"password"`
}

func (u *User) SetID(id string) {
	u.ID = bson.ObjectIdHex(id)
}

func (d *Datasource) GetUsers(id string) ([]User, error) {
	c := d.Session.DB(d.DatabaseName).C("Users")
	var (
		users []User
		err   error
	)
	if id != "" {
		query := c.FindId(bson.ObjectIdHex(id))
		num, err := query.Count()
		if err != nil {
			return nil, err
		}
		if num <= 0 {
			return nil, ErrNoRows
		}
		err = query.All(&users)
	} else {
		err = c.Find(nil).All(&users)
	}
	if err != nil {
		return nil, err
	}

	return users, err
}

func (d *Datasource) InsertUser(user *User) error {
	c := d.Session.DB(d.DatabaseName).C("Users")
	info, err := c.Upsert(new(User), user)
	if err != nil {
		return err
	}
	if info.UpsertedId != nil {
		user.ID = info.UpsertedId.(bson.ObjectId)
	} else {
		return errors.New("datasource: unknown error inserting user")
	}
	return nil
}

func (d *Datasource) RemoveUser(user *User) error {
	c := d.Session.DB(d.DatabaseName).C("Users")
	return c.RemoveId(user.ID)
}

func (d *Datasource) UpdateUser(user *User) error {
	c := d.Session.DB(d.DatabaseName).C("Users")
	return c.UpdateId(user.ID, user)
}
