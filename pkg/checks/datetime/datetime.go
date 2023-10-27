package datetime

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

var (
	DateTimeCheckName           = "Date Time Check"
	DateTimeCheckVersion        = "1.0.0"
	DateTimeCheckTimeNotAllowed = "Current Time Not Within Allowed Time"
	DateTimeCheckDateNotAllowed = "Current Date Not Within Allowed Date"
	DateTimeCheckSuccess        = "Successful Datetime Check"
)

const (
	daily   = 0
	weekly  = 1
	monthly = 2
	yearly  = 3
)

type DateTimeCheck struct {
	jobName     string
	timeScale   uint
	intervals   []int
	tolerance   int
	abortOnFail bool
	mfaOnFail   bool
	hours       int
	minutes     int
	seconds     int
}

func NewDateTimeCheck(config config.CheckConfig, jobName string) DateTimeCheck {
	timeScale := uint(daily)
	scale, exists := config.Options["scale"].(string)
	if exists {
		switch scale {
		case "daily":
			timeScale = daily
		case "weekly":
			timeScale = weekly
		case "monthly":
			timeScale = monthly
		case "yearly":
			timeScale = yearly
		default:
			timeScale = daily
		}
	}

	intervals, exists := config.Options["intervals"].([]int)
	if !exists {
		intervals = []int{0}
	}

	tolerance, exists := config.Options["tolerance"].(int)
	if !exists {
		tolerance = 300
	}

	abortOnFail, exists := config.Options["abortOnFail"].(bool)
	if !exists {
		abortOnFail = true
	}

	mfaOnFail, exists := config.Options["mfaOnFail"].(bool)
	if !exists {
		mfaOnFail = false
	}

	hours := time.Now().Hour()
	minutes := time.Now().Minute()
	seconds := time.Now().Second()
	var err error

	timeString, exists := config.Options["time"].(string)
	if exists {
		timeIntervals := strings.Split(timeString, ":")
		if len(timeIntervals) != 3 {
			hours = time.Now().Hour()
			minutes = time.Now().Minute()
			seconds = time.Now().Second()
		} else {
			hours, err = strconv.Atoi(timeIntervals[0])
			if err != nil {
				hours = time.Now().Hour()
			}

			minutes, err = strconv.Atoi(timeIntervals[1])
			if err != nil {
				minutes = time.Now().Minute()
			}

			seconds, err = strconv.Atoi(timeIntervals[2])
			if err != nil {
				seconds = time.Now().Minute()
			}
		}
	}

	dateTimeCheck := DateTimeCheck{
		jobName:     jobName,
		timeScale:   timeScale,
		intervals:   intervals,
		tolerance:   tolerance,
		abortOnFail: abortOnFail,
		mfaOnFail:   mfaOnFail,
		hours:       hours,
		minutes:     minutes,
		seconds:     seconds,
	}

	return dateTimeCheck
}

func (check *DateTimeCheck) Check(channel chan<- checks.CheckResult, wg *sync.WaitGroup, w whitelist.Whitelist) {
	defer wg.Done()

	result := checks.CheckResult{Name: DateTimeCheckName, Version: DateTimeCheckVersion, Details: ""}

	if !check.checkWithinTime() {
		result.Details = DateTimeCheckTimeNotAllowed

		if check.abortOnFail {
			result.Abort = true
			channel <- result
			return
		} else if check.mfaOnFail {
			result.Mfa = true
			channel <- result
			return
		}
	}

	if !check.checkOnRightDay() {
		result.Details += DateTimeCheckDateNotAllowed

		if check.abortOnFail {
			result.Abort = true
			channel <- result
			return
		} else if check.mfaOnFail {
			result.Mfa = true
			channel <- result
			return
		}
	}

	if len(result.Details) == 0 {
		result.Details = DateTimeCheckSuccess
	}

	channel <- result
}

func (check *DateTimeCheck) checkWithinTime() bool {
	now := time.Now()

	startTime := time.Date(now.Year(), now.Month(), now.Day(), check.hours, check.minutes, check.seconds, 0, time.Local)

	toleranceSeconds := check.seconds + check.tolerance
	addedMintues := toleranceSeconds / 60
	endSeconds := toleranceSeconds % 60

	toleranceMinutes := check.minutes + addedMintues
	addedHours := toleranceMinutes / 60
	endMinutes := toleranceMinutes % 60

	toleranceHours := check.hours + addedHours
	addedDays := toleranceHours / 24
	endHours := toleranceHours % 60

	endTime := time.Date(now.Year(), now.Month(), now.Day()+addedDays, endHours, endMinutes, endSeconds, 0, time.Local)

	checkTime := time.Date(now.Year(), now.Month(), now.Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)

	beforeStart := startTime.Before(checkTime)
	afterEnd := endTime.After(checkTime)

	return (beforeStart && afterEnd) || (startTime.Equal(checkTime) || endTime.Equal(checkTime))
}

func (check *DateTimeCheck) checkOnRightDay() bool {
	for _, interval := range check.intervals {
		success := false
		switch check.timeScale {
		case daily:
			success = true
		case weekly:
			success = int(time.Now().Weekday()) == interval
		case monthly:
			success = time.Now().Day() == interval
		case yearly:
			totalDays := 0
			currentMonth := int(time.Now().Month())
			for month := 1; month < currentMonth; month++ {
				totalDays += time.Date(time.Now().Year(), time.Month(month), 0, 0, 0, 0, 0, time.UTC).Day()
			}
			totalDays += time.Now().Day()
			success = totalDays == interval
		default:
			success = false
		}

		if success {
			return true
		}
	}

	return false
}

func (check *DateTimeCheck) IsValidForCheckType(checkType uint) bool {
	return true
}

func (check *DateTimeCheck) IsValidForPlatform(ciPlatform string) bool {
	return true
}
