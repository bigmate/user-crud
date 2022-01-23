package user

import (
	"user-crud/pkg/filter"
)

//FilterFactory is a factory to create filter.Filter interfaces
type FilterFactory interface {
	ByFirstName(firstName string) filter.Filter
	ByLastName(lastName string) filter.Filter
	ByCountry(country string) filter.Filter
	ByEmail(email string) filter.Filter
}

type filterFactory struct{}

func (f filterFactory) ByFirstName(firstName string) filter.Filter {
	if firstName == "" {
		return noop{}
	}

	return &byField{
		name:  "first_name",
		value: firstName,
		cmp:   filter.EQ,
	}
}

func (f filterFactory) ByLastName(lastName string) filter.Filter {
	if lastName == "" {
		return noop{}
	}

	return &byField{
		name:  "last_name",
		value: lastName,
		cmp:   filter.EQ,
	}
}

func (f filterFactory) ByCountry(country string) filter.Filter {
	if country == "" {
		return noop{}
	}

	return &byField{
		name:  "country",
		value: country,
		cmp:   filter.EQ,
	}
}

func (f filterFactory) ByEmail(email string) filter.Filter {
	if email == "" {
		return noop{}
	}

	return &byField{
		name:  "email",
		value: email,
		cmp:   filter.EQ,
	}
}

type byField struct {
	name  string
	value interface{}
	cmp   filter.Comparator
}

func (b *byField) Evaluate() (field string, value interface{}, cmp filter.Comparator) {
	return b.name, b.value, b.cmp
}

type noop struct{}

func (n noop) Evaluate() (field string, value interface{}, cmp filter.Comparator) {
	return "", nil, filter.NOOP
}
