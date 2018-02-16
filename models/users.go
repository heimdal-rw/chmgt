package models

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
