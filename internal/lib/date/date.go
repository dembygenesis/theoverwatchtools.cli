package date

import (
	"github.com/dembygenesis/local.tools/internal/utils/error_util"
	"github.com/golang-module/carbon/v2"
)

func isValidStartAndEndDates(start, end string) error {
	var errs error_util.List

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
