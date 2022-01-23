package paginator

const maxPageSize = 100

//Paginator is the pagination interface
type Paginator interface {
	Page() uint64
	Size() uint64
}

type paginator struct {
	page uint64
	size uint64
}

func (p paginator) Page() uint64 {
	return p.page
}

func (p paginator) Size() uint64 {
	return p.size
}

func New(page, size uint64) Paginator {
	p := paginator{
		page: page,
		size: size,
	}

	if page <= 0 {
		p.page = 1
	}

	if size <= 0 {
		p.size = 1
	}

	if size > maxPageSize {
		p.size = maxPageSize
	}

	return p
}
