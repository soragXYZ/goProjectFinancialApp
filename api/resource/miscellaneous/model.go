package miscellaneous

type Category struct {
	code        string
	parent_code string
}

type Currency struct {
	id        string
	name      string
	symbol    string
	precision uint
}

type PaginationLinks struct {
	self string
	prev string
	next string
}
