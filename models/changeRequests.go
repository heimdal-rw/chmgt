package models

var CollectionChangeRequests = "ChangeRequests"

func (d *Datasource) GetChangeRequests(id string) ([]Item, error) {
	return d.GetItems(id, CollectionChangeRequests)
}

func (d *Datasource) InsertChangeRequest(cr Item) error {
	return d.InsertItem(cr, CollectionChangeRequests)
}

func (d *Datasource) RemoveChangeRequest(cr Item) error {
	return d.RemoveItem(cr, CollectionChangeRequests)
}

func (d *Datasource) UpdateChangeRequest(cr Item) error {
	return d.UpdateItem(cr, CollectionChangeRequests)
}
