package dateutil

import (
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/golang-module/carbon/v2"
	"time"
)

const (
	layoutDate = "2006-01-02"
	layoutTS   = "2006-01-02T15:04:05.000Z"
)

func IsValid(s string) bool {
	_, err := time.Parse(layoutDate, s)
	return err != nil
}

func isValidStartAndEndDates(start, end string) error {
	var errs errs.List

	t1 := carbon.Parse(start)
	t2 := carbon.Parse(end)

	if t1.IsInvalid() {
		errs.Add("start is invalid")
	}
	if t2.IsInvalid() {
		errs.Add("end is invalid")
	}
	if t1.Compare(">", t2) && (t1.IsValid() && t2.IsValid()) {
		errs.Add("start is greater than end")
	}

	return errs.Single()
}

func NewDateRanges(start, end string) ([]*DateRange, error) {
	if err := isValidStartAndEndDates(start, end); err != nil {
		return nil, err
	}
	return getCarbonRanges(carbon.Parse(start), carbon.Parse(end)), nil
}

type DateRange struct {
	Start string
	End   string
}

func getCarbonRanges(t1, t2 carbon.Carbon) []*DateRange {
	diff := int(t1.DiffInMonths(t2)) + 1

	var curr carbon.Carbon
	var startTime string
	var endTime string

	dateRanges := make([]*DateRange, 0)
	for i := 0; i < diff; i++ {
		if i == 0 {
			curr = t1.StartOfMonth()
		} else {
			curr = curr.AddMonth()
		}

		startTime = curr.String()

		if i == diff-1 {
			endTime = t2.String()
		} else {
			endTime = curr.AddMonth().SubDays(1).String()
		}

		dateRanges = append(dateRanges, &DateRange{
			Start: startTime[0:10],
			End:   endTime[0:10],
		})
	}
	return dateRanges
}
