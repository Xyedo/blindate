package pagination


type Pagination struct {
	Page int
	Limit int
}

func (p Pagination) Offset() int {
	return p.Page * p.Limit - p.Limit
}