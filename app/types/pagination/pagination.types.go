package pagination

type Meta struct {
	Limit       int  `json:"limit"`
	Offset      int  `json:"offset"`
	Skip        int  `json:"skip"`
	Page        int  `json:"page"`
	NextPage    int  `json:"nextPage"`
	PrevPage    int  `json:"prevPage"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"totalPages"`
	HasNextPage bool `json:"hasNextPage"`
	HasPrevPage bool `json:"hasPrevPage"`

	Filters []FilterEntry `json:"-"`
	Search  string        `json:"-"`
	Order   string        `json:"-"`
}

type FilterEntry struct {
	Key      string
	Operator string
	Value    string
}
