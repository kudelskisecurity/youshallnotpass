package datetime

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

func TestNewDateTimeCheck(t *testing.T) {
	newDateTimeCheckTests := []struct {
		name        string
		config      config.CheckConfig
		jobName     string
		timeScale   uint
		intervals   []int
		tolerance   int
		abortOnFail bool
		mfaOnFail   bool
		hours       int
		minutes     int
	}{
		{
			name: "new date time check with daily scale",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"scale": "daily",
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			tolerance:   300,
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "new date time check with weekly scale",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"scale": "weekly",
				},
			},
			jobName:     "testJob",
			timeScale:   weekly,
			intervals:   []int{0},
			tolerance:   300,
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "new date time check with monthly scale",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"scale": "monthly",
				},
			},
			jobName:     "testJob",
			timeScale:   monthly,
			intervals:   []int{0},
			tolerance:   300,
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "new date time check with yearly scale",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"scale": "yearly",
				},
			},
			jobName:     "testJob",
			timeScale:   yearly,
			intervals:   []int{0},
			tolerance:   300,
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "new date time check with interval",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"scale":     "monthly",
					"intervals": []int{1, 5, 9},
				},
			},
			jobName:     "testJob",
			timeScale:   monthly,
			intervals:   []int{1, 5, 9},
			tolerance:   300,
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "test new date time check tolerance",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"tolerance": 10,
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			tolerance:   10,
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "test new date time check abortOnFail=false",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"abortOnFail": false,
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			tolerance:   300,
			abortOnFail: false,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "test new date time check mfaOnFail=true",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"abortOnFail": false,
					"mfaOnFail":   true,
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			abortOnFail: false,
			mfaOnFail:   true,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "test new date time check with time given",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"time": "04:10:30",
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       4,
			minutes:     10,
		},
		{
			name: "test new date time check with time given (pm)",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"time": "14:30:00",
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       14,
			minutes:     30,
		},
		{
			name: "test new date time bad time fall-back",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"time": "04:04",
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "test new date time (bad time letters instead on numbers)",
			config: config.CheckConfig{
				Name: "datetimeCheck",
				Options: map[string]interface{}{
					"time": "HH:MM:SS",
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
		{
			name: "test new date time check (unknown scale defaults to daily)",
			config: config.CheckConfig{
				Name: "datetimecheck",
				Options: map[string]interface{}{
					"scale": "unknown",
				},
			},
			jobName:     "testJob",
			timeScale:   daily,
			intervals:   []int{0},
			abortOnFail: true,
			mfaOnFail:   false,
			hours:       time.Now().Hour(),
			minutes:     time.Now().Minute(),
		},
	}

	for testNum, test := range newDateTimeCheckTests {
		check := NewDateTimeCheck(test.config, test.jobName)

		if check.jobName != test.jobName {
			t.Errorf("\n%d) unexpected job name -\nexpected: (%s)\ngot: (%s)", testNum, test.jobName, check.jobName)
		}

		if check.timeScale != test.timeScale {
			t.Errorf("\n%d) unexpected time scale -\nexpected: (%d)\ngot: (%d)", testNum, test.timeScale, check.timeScale)
		}

		if !reflect.DeepEqual(check.intervals, test.intervals) {
			t.Errorf("\n%d) unexpected interval -\nexpected: (%+v)\ngot: (%+v)", testNum, test.intervals, check.intervals)
		}

		if check.abortOnFail != test.abortOnFail {
			t.Errorf("\n%d) unexpected abortOnFail -\nexpected: (%t)\ngot: (%t)", testNum, test.abortOnFail, check.abortOnFail)
		}

		if check.mfaOnFail != test.mfaOnFail {
			t.Errorf("\n%d) unexpected mfaOnFail -\nexpected: (%t)\ngot: (%t)", testNum, test.mfaOnFail, check.mfaOnFail)
		}

		if check.hours != test.hours {
			t.Errorf("\n%d) unexpected hours -\nexpected: (%d)\ngot: (%d)", testNum, test.hours, check.hours)
		}

		if check.minutes != test.minutes {
			t.Errorf("\n%d) unexpected minutes -\nexpected: (%d)\ngot: (%d)", testNum, test.minutes, check.minutes)
		}
	}
}
func TestDateTimeCheckIsValidForPlatform(t *testing.T) {
	dateTimeCheckValidForPlatformTests := []struct {
		name     string
		platform string
		valid    bool
	}{
		{
			name:     "date time check should be valid for github",
			platform: "github",
			valid:    true,
		},
		{
			name:     "date time check should be valid for gitlab",
			platform: "gitlab",
			valid:    true,
		},
	}

	dateTimeCheck := DateTimeCheck{}
	for testNum, test := range dateTimeCheckValidForPlatformTests {
		valid := dateTimeCheck.IsValidForPlatform(test.platform)
		if valid != test.valid {
			t.Errorf("\n%d) unexpected valid for platform -\nexpected: (%t)\ngot: (%t)", testNum, test.valid, valid)
		}
	}
}
func TestDateTimeCheckIsValidForCheckTypes(t *testing.T) {
	dateTimeCheckValidForCheckTypeTests := []struct {
		name      string
		checkType uint
		valid     bool
	}{
		{
			name:      "date time check should be valid for image checks",
			checkType: checks.ImageCheck,
			valid:     true,
		},
		{
			name:      "date time check should be valid for script checks",
			checkType: checks.ScriptCheck,
			valid:     true,
		},
		{
			name:      "date time check should be valid for all checks",
			checkType: checks.All,
			valid:     true,
		},
	}

	dateTimeCheck := DateTimeCheck{}
	for testNum, test := range dateTimeCheckValidForCheckTypeTests {
		valid := dateTimeCheck.IsValidForCheckType(test.checkType)
		if valid != test.valid {
			t.Errorf("\n%d) unexpected valid for check type -\nexpected (%t)\ngot: (%t)", testNum, test.valid, valid)
		}
	}
}
func TestCheckWithinTime(t *testing.T) {
	checkWithinTimeTests := []struct {
		name          string
		dateTimeCheck DateTimeCheck
		expected      bool
	}{
		{
			name: "check within time (within time)",
			dateTimeCheck: DateTimeCheck{
				jobName:   "testJob",
				timeScale: daily,
				intervals: []int{0},
				tolerance: 300,
				hours:     time.Now().Hour(),
				minutes:   time.Now().Minute(),
				seconds:   time.Now().Second(),
			},
			expected: true,
		},
		{
			name: "check within time (before time)",
			dateTimeCheck: DateTimeCheck{
				jobName:   "testJob",
				timeScale: daily,
				intervals: []int{0},
				tolerance: 300,
				hours:     time.Now().Hour() + 1,
				minutes:   time.Now().Minute(),
				seconds:   time.Now().Second(),
			},
			expected: false,
		},
		{
			name: "test check with time (after time)",
			dateTimeCheck: DateTimeCheck{
				jobName:   "testJob",
				timeScale: daily,
				intervals: []int{0},
				tolerance: 10,
				hours:     time.Now().Hour(),
				minutes:   time.Now().Minute() - 1,
				seconds:   time.Now().Second(),
			},
			expected: false,
		},
	}

	for testNum, test := range checkWithinTimeTests {
		withinTime := test.dateTimeCheck.checkWithinTime()
		if withinTime != test.expected {
			t.Errorf("\n%d) unexpected within time -\nexpected: (%t)\ngot: (%t)", testNum, test.expected, withinTime)
		}
	}
}
func TestCheckOnRightDay(t *testing.T) {
	getCurrentDayYear := func() int {
		totalDays := 0
		currentMonth := int(time.Now().Month())
		for month := 1; month < currentMonth; month++ {
			totalDays += time.Date(time.Now().Year(), time.Month(month), 0, 0, 0, 0, 0, time.UTC).Day()
		}
		totalDays += time.Now().Day()
		return totalDays
	}

	checkOnRightDayTests := []struct {
		name          string
		dateTimeCheck DateTimeCheck
		rightDay      bool
	}{
		{
			name: "test on right day (daily)",
			dateTimeCheck: DateTimeCheck{
				intervals: []int{0},
				timeScale: daily,
			},
			rightDay: true,
		},
		{
			name: "test on right day (weekly)",
			dateTimeCheck: DateTimeCheck{
				intervals: []int{int(time.Now().Weekday())},
				timeScale: weekly,
			},
			rightDay: true,
		},
		{
			name: "test on wrong day (weekly)",
			dateTimeCheck: DateTimeCheck{
				intervals: []int{int(time.Now().Weekday()) + 1},
				timeScale: weekly,
			},
			rightDay: false,
		},
		{
			name: "test on right day (monthly)",
			dateTimeCheck: DateTimeCheck{
				intervals: []int{int(time.Now().Day())},
				timeScale: monthly,
			},
			rightDay: true,
		},
		{
			name: "test on wrong day (monthly)",
			dateTimeCheck: DateTimeCheck{
				intervals: []int{int(time.Now().Day() + 1)},
				timeScale: monthly,
			},
			rightDay: false,
		},
		{
			name: "test on right day (yearly)",
			dateTimeCheck: DateTimeCheck{
				intervals: []int{getCurrentDayYear()},
				timeScale: yearly,
			},
			rightDay: true,
		},
		{
			name: "test on wrong day (yearly)",
			dateTimeCheck: DateTimeCheck{
				intervals: []int{getCurrentDayYear() + 1},
				timeScale: yearly,
			},
			rightDay: false,
		},
	}

	for testNum, test := range checkOnRightDayTests {
		rightDay := test.dateTimeCheck.checkOnRightDay()
		if rightDay != test.rightDay {
			t.Errorf("\n%d) unexpected rightDay -\nexpected: (%t)\ngot: (%t)", testNum, test.rightDay, rightDay)
		}
	}
}
func TestDateTimeCheck(t *testing.T) {
	dateTimeCheckTests := []struct {
		name           string
		dateTimeCheck  DateTimeCheck
		expectedResult checks.CheckResult
	}{
		{
			name: "test successful date time check",
			dateTimeCheck: DateTimeCheck{
				jobName:     "testJob",
				timeScale:   weekly,
				intervals:   []int{int(time.Now().Weekday())},
				tolerance:   300,
				abortOnFail: true,
				mfaOnFail:   true,
				hours:       time.Now().Hour(),
				minutes:     time.Now().Minute(),
				seconds:     time.Now().Second(),
			},
			expectedResult: checks.CheckResult{
				Name:    DateTimeCheckName,
				Version: DateTimeCheckVersion,
				Error:   nil,
				Abort:   false,
				Mfa:     false,
				Details: DateTimeCheckSuccess,
			},
		},
		{
			name: "test unsuccessful date time check (wrong time) - abort",
			dateTimeCheck: DateTimeCheck{
				jobName:     "testJob",
				timeScale:   weekly,
				intervals:   []int{int(time.Now().Weekday())},
				tolerance:   300,
				abortOnFail: true,
				mfaOnFail:   true,
				hours:       time.Now().Hour() + 1,
				minutes:     time.Now().Minute(),
				seconds:     time.Now().Second(),
			},
			expectedResult: checks.CheckResult{
				Name:    DateTimeCheckName,
				Version: DateTimeCheckVersion,
				Error:   nil,
				Abort:   true,
				Mfa:     false,
				Details: DateTimeCheckTimeNotAllowed,
			},
		},
		{
			name: "test unsuccessful date time check (wrong time) - mfa",
			dateTimeCheck: DateTimeCheck{
				jobName:     "testJob",
				timeScale:   weekly,
				intervals:   []int{int(time.Now().Weekday())},
				tolerance:   300,
				abortOnFail: false,
				mfaOnFail:   true,
				hours:       time.Now().Hour() - 1,
				minutes:     time.Now().Minute(),
				seconds:     time.Now().Second(),
			},
			expectedResult: checks.CheckResult{
				Name:    DateTimeCheckName,
				Version: DateTimeCheckVersion,
				Error:   nil,
				Abort:   false,
				Mfa:     true,
				Details: DateTimeCheckTimeNotAllowed,
			},
		},
		{
			name: "test unsuccessful date time check (wrong date) - abort",
			dateTimeCheck: DateTimeCheck{
				jobName:     "testJob",
				timeScale:   weekly,
				intervals:   []int{int(time.Now().Weekday() + 1)},
				tolerance:   300,
				abortOnFail: true,
				mfaOnFail:   true,
				hours:       time.Now().Hour(),
				minutes:     time.Now().Minute(),
				seconds:     time.Now().Second(),
			},
			expectedResult: checks.CheckResult{
				Name:    DateTimeCheckName,
				Version: DateTimeCheckVersion,
				Error:   nil,
				Abort:   true,
				Mfa:     false,
				Details: DateTimeCheckDateNotAllowed,
			},
		},
		{
			name: "test unsuccessful date time check (wrong date) - mfa",
			dateTimeCheck: DateTimeCheck{
				jobName:     "testJob",
				timeScale:   weekly,
				intervals:   []int{int(time.Now().Weekday() - 1)},
				tolerance:   300,
				abortOnFail: false,
				mfaOnFail:   true,
				hours:       time.Now().Hour(),
				minutes:     time.Now().Minute(),
				seconds:     time.Now().Second(),
			},
			expectedResult: checks.CheckResult{
				Name:    DateTimeCheckName,
				Version: DateTimeCheckVersion,
				Error:   nil,
				Abort:   false,
				Mfa:     true,
				Details: DateTimeCheckDateNotAllowed,
			},
		},
	}

	whitelist := whitelist.Whitelist{}
	for testNum, test := range dateTimeCheckTests {
		wg := sync.WaitGroup{}
		channel := make(chan checks.CheckResult, 1)

		wg.Add(1)
		go test.dateTimeCheck.Check(channel, &wg, whitelist)

		wg.Wait()
		close(channel)

		result := <-channel

		if !result.CompareCheckResult(test.expectedResult) {
			t.Errorf("\n%d) Results not equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedResult, result)
		}
	}
}
