package models

import (
	"errors"
	"time"
)

var monthYearLayout = "01/2006"
var dayLayout = "02"

var expirationDateParseError = errors.New(`ExpirationDateParseError: should be a string formatted as "01/2006"`)
var dueDateParseError = errors.New(`DueDateParseError: should be a string formatted as "02"`)

type expirationDate struct {
	*time.Time
}

func (ed expirationDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ed.Format(monthYearLayout) + `"`), nil
}

func (ed *expirationDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) != 9 {
		return expirationDateParseError
	}
	ret, err := time.Parse(monthYearLayout, s[1:8])
	if err != nil {
		return err
	}
	ed.Time = &ret
	return nil
}

type dueDay struct {
	*time.Time
}

func (dd dueDay) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dd.Format(dayLayout) + `"`), nil
}

func (dd *dueDay) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) != 4 {
		return dueDateParseError
	}
	ret, err := time.Parse(dayLayout, s[1:3])
	if err != nil {
		return err
	}
	dd.Time = &ret
	return nil
}

func (dd dueDay) NextDueDate() time.Time {
	nextDueDateYear, nextDueDateMonth, nowDay := time.Now().Date()
	nextDueDateDay := dd.Day()
	
	if nextDueDateDay <= nowDay {
		if nextDueDateMonth == time.December {
			nextDueDateMonth = time.January
			nextDueDateYear += 1
		} else {
			nextDueDateMonth += 1
		}
	}

	return time.Date(nextDueDateYear, nextDueDateMonth, nextDueDateDay, 0, 0, 0, 0, time.UTC)
}