package errutil

import (
	"errors"
	"sync"
)

type List []error

var m sync.Mutex

// ToArr wraps error as string array,
// primarily for formatting issues.
func ToArr(err error) []string {
	if err == nil {
		return make([]string, 0)
	}
	return []string{err.Error()}
}

func (e *List) AddErr(err error) {
	m.Lock()
	defer m.Unlock()
	if e == nil {
		*e = make([]error, 0)
	}
	*e = append(*e, err)
}

func (e *List) Add(err string) {
	m.Lock()
	defer m.Unlock()
	if e == nil {
		*e = make([]error, 0)
	}
	*e = append(*e, errors.New(err))
}

// Single will return the error string value
func (e *List) Single() error {
	m.Lock()
	defer m.Unlock()

	if e == nil {
		return nil
	}

	deref := *e

	switch len(deref) {
	case 0:
		return nil

	case 1:
		return deref[0]
	}

	var bs []byte
	for _, err := range deref {
		if err == nil {
			continue
		}
		bs = append(bs, []byte(err.Error())...)
		bs = append(bs, ',', ' ')
	}

	return errors.New(string(bs))
}

func (e *List) HasErrors() bool {
	if e == nil {
		return false
	}
	return len(*e) != 0
}
