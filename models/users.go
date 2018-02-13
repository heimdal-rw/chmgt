package models

import (
	"encoding/json"
	"errors"

	"gopkg.in/mgo.v2/bson"
)

var CollectionUsers = "Users"

type User struct {
	ID        bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	FirstName string                 `json:"firstname"`
	LastName  string                 `json:"lastname"`
	UserName  string                 `json:"username"`
	Email     string                 `json:"email"`
	Password  string                 `json:"password"`
	Extras    map[string]interface{} `bson:",inline" json:"-"`
}

func (u *User) SetID(id string) {
	u.ID = bson.ObjectIdHex(id)
}

func (u *User) MarshalJSON() ([]byte, error) {
	var j interface{}
	data, err := bson.Marshal(u)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(data, &j)
	return json.Marshal(&j)
}

func (u *User) UnmarshalJSON(p []byte) error {
	var j map[string]interface{}
	json.Unmarshal(p, &j)
	data, err := bson.Marshal(&j)
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, u)
}

func (d *Datasource) GetUsers(id string) ([]User, error) {
	c := d.Session.DB(d.DatabaseName).C(CollectionUsers)
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
	c := d.Session.DB(d.DatabaseName).C(CollectionUsers)
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
	c := d.Session.DB(d.DatabaseName).C(CollectionUsers)
	return c.RemoveId(user.ID)
}

func (d *Datasource) UpdateUser(user *User) error {
	c := d.Session.DB(d.DatabaseName).C(CollectionUsers)
	return c.UpdateId(user.ID, user)
}
