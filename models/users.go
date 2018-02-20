package models

import "gopkg.in/mgo.v2/bson"

// CollectionUsers is the name of the collection for users
var CollectionUsers = "Users"

// GetUsers is a wrapper to get user items
func (d *Datasource) GetUsers(id string) ([]Item, error) {
	return d.GetItems(id, CollectionUsers)
}

// InsertUser is a wrapper to insert a user
func (d *Datasource) InsertUser(user Item) error {
	return d.InsertItem(user, CollectionUsers)
}

// RemoveUser is a wrapper to remove a user
func (d *Datasource) RemoveUser(user Item) error {
	return d.RemoveItem(user, CollectionUsers)
}

// UpdateUser is a wrapper to update a user
func (d *Datasource) UpdateUser(user Item) error {
	return d.UpdateItem(user, CollectionUsers)
}

// ValidateUser checks the user credentials in the database
func (d *Datasource) ValidateUser(username, password string) (bool, error) {
	c := d.Session.DB(d.DatabaseName).C(CollectionUsers)
	query := bson.M{
		"username": username,
		"password": password,
	}
	num, err := c.Find(query).Count()
	if err != nil {
		return false, err
	}
	// Since usernames are unique, there should be only one record
	if num != 1 {
		return false, nil
	}
	return true, nil
}
