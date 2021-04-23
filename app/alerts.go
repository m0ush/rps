package app

// TODO: Include Secruity Attributes (e.g. Sedol, Name)
type Alert struct {
	SecID int
	Limit
}

func NewAlert(id int, lim Limit) Alert {
	return Alert{
		SecID: id,
		Limit: lim,
	}
}
