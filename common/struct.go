package common

type ListParam struct {
	Pa         Pagiation        `json:"pa"`
	Filter     string           `json:"filter"`
	Order      string           `json:"order"`
	Status      string           `json:"status"`
}
type Pagiation struct {
	Pn int `json:"pn"`
	Ps int `json:"ps"`
	Total int `json:"total"`
}
