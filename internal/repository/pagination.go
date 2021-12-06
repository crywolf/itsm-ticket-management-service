package repository

// Pagination encapsulates information about pagination data
type Pagination struct {
	Total             int
	Size              int
	Page              int
	Prev              int
	Next              int
	First             int
	Last              int
	FirstElementIndex int
	LastElementIndex  int
}

// NewPagination creates new initialized pagination object
func NewPagination(collectionLength int, page, perPage uint) *Pagination {
	total := uint(collectionLength)
	next := page + 1
	prev := page - 1

	if page == 1 && total < perPage {
		prev = 0
	}

	start := (page - 1) * perPage
	if start >= total {
		start = total
	}

	end := start + perPage
	if end > total {
		end = total
		next = 0
	}

	lastIndex := int(end) - 1
	if lastIndex < 0 {
		lastIndex = 0
	}

	last := (total / perPage) + 1
	if (total % perPage) == 0 {
		last--
	}

	if next > last {
		next = 0
	}

	if last == 0 {
		last = 1
	}

	return &Pagination{
		Total:             int(total),
		Size:              int(end - start),
		Page:              int(page),
		Prev:              int(prev),
		Next:              int(next),
		First:             1,
		Last:              int(last),
		FirstElementIndex: int(start),
		LastElementIndex:  lastIndex,
	}
}
