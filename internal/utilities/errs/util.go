package errs

import (
	"errors"
	"sync"
)

type Util struct {
	StatusCode int
	List       List
}

func (e *Util) Error() string {
	return e.List.Single().Error()
}

type Cfg struct {
	StatusCode int
	Err        error
	Errs       List
}

// New returns a new error util instance.
func New(cfg *Cfg) error {
	u := &Util{
		StatusCode: cfg.StatusCode,
	}
	if cfg.Err != nil {
		u.List = []error{cfg.Err}
	} else if cfg.Errs != nil {
		u.List = cfg.Errs
	}
	return u

}

type List []error

var m sync.Mutex

// ToArr wraps error as string array,
// primarily for formatting issues.
func ToArr(err error) []string {
	arrs := make([]string, 0)
	if err == nil {
		return arrs
	}

	var util *Util
	if errors.As(err, &util) {
		for _, v := range util.List {
			arrs = append(arrs, v.Error())
		}
		return arrs
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

// ErrAsUtil checks if the error is an underlying
// *Util object.
func ErrAsUtil(i interface{}) (*Util, bool) {
	_, ok := i.(error)
	if !ok {
		return nil, false
	}

	var util *Util
	ok = errors.As(i.(error), &util)
	if !ok {
		return nil, false
	}

	return util, true
}

// AsErrUtil returns the List alongside a statusCode,
// very useful if we want to add more context on what API
func (e *List) AsErrUtil(statusCode int) error {
	return &Util{
		List:       *e,
		StatusCode: statusCode,
	}
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

func (e *List) ErrsAsStrArr() []string {
	if e == nil {
		return nil
	}
	arr := make([]string, 0)
	for _, v := range *e {
		arr = append(arr, v.Error())
	}
	return arr
}
