package hgtealib

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Filter struct {
	stockedOnly bool
	samplesOnly bool
	types       map[string]struct{}
}

func (f *Filter) StockedOnly() *Filter {
	f.stockedOnly = true
	return f
}

func (f *Filter) SamplesOnly() *Filter {
	f.samplesOnly = true
	return f
}

func (f *Filter) Types(v []string) *Filter {
	if len(v) > 0 {
		for _, t := range v {
			if t != "" {
				f.Type(t)
			}
		}
	}
	return f
}

func (f *Filter) Type(v string) *Filter {
	f.types[strings.ToLower(v)] = struct{}{}
	return f
}

func NewFilter() *Filter {
	f := new(Filter)

	f.stockedOnly = false
	f.samplesOnly = false
	f.types = make(map[string]struct{})

	return f
}

type TeaDb struct {
	teas          map[int]Tea
	log           map[time.Time]Entry
	logSortedKeys TimeSlice
}

func (d *TeaDb) Teas(filter *Filter) (map[int]Tea, error) {
	teas := make(map[int]Tea)
	for k, v := range d.teas {
		// Now apply the filters
		if filter.stockedOnly && !v.Storage.Stocked {
			continue
		}

		// if filter.samplesOnly && !strings.Contains(strings.ToLower(v.Size), "sample") {
		// continue
		// }

		if len(filter.types) > 0 {
			if _, ok := filter.types[strings.ToLower(v.Type)]; !ok {
				continue
			}
		}

		teas[k] = v
	}
	return teas, nil
}

func (d *TeaDb) Tea(id int) (Tea, error) {
	if val, ok := d.teas[id]; ok {
		return val, nil
	}
	return *new(Tea), errors.New(fmt.Sprintf("Could not retrieve Tea by id: %d", id))
}

func (d *TeaDb) Log(filter *Filter) ([]Entry, error) {
	log := make([]Entry, 0)
	for _, k := range d.logSortedKeys {
		log = append(log, d.log[k])
	}
	return log, nil
}

func newTeaDb(teas []*Tea, entries []*Entry) (*TeaDb, error) {
	db := new(TeaDb)
	db.teas = make(map[int]Tea)
	db.log = make(map[time.Time]Entry)

	for _, tea := range teas {
		if tea != nil {
			db.teas[tea.Id] = *tea
		}

	}

	for _, entry := range entries {
		if entry != nil {
			db.log[entry.DateTime] = *entry
			db.logSortedKeys = append(db.logSortedKeys, entry.DateTime)
			sort.Sort(db.logSortedKeys)

			if tea, ok := db.teas[entry.Tea]; ok {
				tea.Add(*entry)
				db.teas[entry.Tea] = tea // TODO: why do I have to do this?
			}
		}
	}

	return db, nil
}
