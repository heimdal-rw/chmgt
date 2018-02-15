package models

var CollectionUsers = "Users"

func (d *Datasource) GetUsers(id string) ([]Item, error) {
	return d.GetItems(id, CollectionUsers)
}

func (d *Datasource) InsertUser(user Item) error {
	return d.InsertItem(user, CollectionUsers)
}

func (d *Datasource) RemoveUser(user Item) error {
	return d.RemoveItem(user, CollectionUsers)
}

func (d *Datasource) UpdateUser(user Item) error {
	return d.UpdateItem(user, CollectionUsers)
}
