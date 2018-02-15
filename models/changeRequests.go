package models

var CollectionChangeRequests = "ChangeRequests"

func (d *Datasource) GetChangeRequests(id string) ([]Item, error) {
	return d.GetItems(id, CollectionChangeRequests)
}

func (d *Datasource) InsertChangeRequest(user Item) error {
	return d.InsertItem(user, CollectionChangeRequests)
}

func (d *Datasource) RemoveChangeRequest(user Item) error {
	return d.RemoveItem(user, CollectionChangeRequests)
}

func (d *Datasource) UpdateChangeRequest(user Item) error {
	return d.UpdateItem(user, CollectionChangeRequests)
}
