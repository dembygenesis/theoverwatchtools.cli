package utils_common

import (
	"errors"
	"sync"
)

type ErrorList []error

var m sync.Mutex

func (e *ErrorList) AddErr(err error) {
	m.Lock()
	defer m.Unlock()
	if e == nil {
		*e = make([]error, 0)
	}
	*e = append(*e, err)
}

func (e *ErrorList) Add(err string) {
	m.Lock()
	defer m.Unlock()
	if e == nil {
		*e = make([]error, 0)
	}
	*e = append(*e, errors.New(err))
}

// Single will return the error string value
func (e *ErrorList) Single() error {
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
		bs = append(bs, ',', '\n')
	}

	return errors.New(string(bs))
}
