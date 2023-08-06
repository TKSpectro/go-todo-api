package pagination

type QueryParams struct {
	Page    int    `json:"page" minimum:"1" default:"1" example:"1"`
	Limit   int    `json:"limit" minimum:"1" default:"10" example:"10"`
	Order   string `json:"order" default:"id asc" example:"id asc"`
	Search  string `json:"search" example:"test@test.com" description:"Searches for the given string in all searchable fields (marked with x-search)"`
	Filters string `json:"filters" example:"[amount][gte]=5 or [fk_id]=5. This can be given multiple times" description:"Filters the result by the given filters. The format is [key][operator]=[value]"`
}

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
