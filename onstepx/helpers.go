package onstepx

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

var reHH = regexp.MustCompile(`^([\+-]?)(\d+)$`)

func parseHH(s string) (int, error) {
	var err error
	match := reHH.FindStringSubmatch(s)
	if match == nil {
		return 0, errors.New("invalid format - should be [s]HH")
	}
	hh, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, fmt.Errorf("error parsing hours - %w", err)
	}
	result := hh
	if match[1] == "-" {
		result = -result
	}
	return result, nil
}

var reHHMMSS = regexp.MustCompile(`^(\d+):(\d+):(\d+(?:\.\d+)?)$`)

func parseHHMMSS(s string) (float64, error) {
	var err error
	match := reHHMMSS.FindStringSubmatch(s)
	if match == nil {
		return 0.0, errors.New("Invalid format - should be HH:MM:SS.SSS")
	}
	hrs, err := strconv.Atoi(match[1])
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing hours - %w", err)
	}
	if hrs < 0 || hrs > 23 {
		return 0.0, errors.New("Invalid hours - should be 0-23")
	}
	mins, err := strconv.Atoi(match[2])
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing minutes - %w", err)
	}
	if mins < 0 || mins > 59 {
		return 0.0, errors.New("Invalid minutes - should be 0-59")
	}
	secs, err := strconv.ParseFloat(match[3], 64)
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing seconds - %w", err)
	}
	if secs < 0 || secs >= 60 {
		return 0.0, errors.New("Invalid seconds - should be 0-59")
	}
	return float64(hrs) + float64(mins)/60 + secs/3600, nil
}

var reHHMM = regexp.MustCompile(`^(\d+):(\d+)$`)

func parseHHMM(s string) (float64, error) {
	var err error
	match := reHHMM.FindStringSubmatch(s)
	if match == nil {
		return 0.0, errors.New("Invalid format - should be HH:MM")
	}
	hrs, err := strconv.Atoi(match[1])
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing hours - %w", err)
	}
	if hrs < 0 || hrs > 23 {
		return 0.0, errors.New("Invalid hours - should be 0-23")
	}
	mins, err := strconv.Atoi(match[2])
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing minutes - %w", err)
	}
	if mins < 0 || mins >= 60 {
		return 0.0, errors.New("Invalid minutes - should be 0-60")
	}
	return float64(hrs) + float64(mins)/60, nil
}

var reDDMMSS = regexp.MustCompile(`^([\+-]?)(\d+)\*(\d+):(\d+(?:\.\d+)?)$`)

func parseDDMMSS(s string) (float64, error) {
	var err error
	match := reDDMMSS.FindStringSubmatch(s)
	if match == nil {
		return 0.0, errors.New("Invalid format - should be sDD*MM:SS.SSS")
	}
	degs, err := strconv.Atoi(match[2])
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing degrees - %w", err)
	}
	if degs < 0 || degs > 359 {
		return 0.0, errors.New("Invalid degrees - should be 0-359")
	}
	mins, err := strconv.Atoi(match[3])
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing minutes - %w", err)
	}
	if mins < 0 || mins > 59 {
		return 0.0, errors.New("Invalid minutes - should be 0-59")
	}
	secs, err := strconv.ParseFloat(match[4], 64)
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing seconds - %w", err)
	}
	if secs < 0 || secs >= 60 {
		return 0.0, errors.New("Invalid seconds - should be 0-60")
	}
	result := float64(degs) + float64(mins)/60 + secs/3600
	if match[1] == "-" {
		result = -result
	}
	return result, nil
}

var reDDMM = regexp.MustCompile(`^([\+-]?)(\d+)\*(\d+)$`)

func parseDDMM(s string) (float64, error) {
	var err error
	match := reDDMM.FindStringSubmatch(s)
	if match == nil {
		return 0.0, errors.New("invalid format - should be [s]DDD*MM")
	}
	degs, err := strconv.Atoi(match[2])
	if err != nil {
		return 0.0, fmt.Errorf("error parsing degrees - %w", err)
	}
	if degs < 0 || degs > 359 {
		return 0.0, errors.New("invalid degrees - should be 0-359")
	}
	mins, err := strconv.Atoi(match[3])
	if err != nil {
		return 0.0, fmt.Errorf("error parsing minutes - %w", err)
	}
	if mins < 0 || mins >= 60 {
		return 0.0, errors.New("invalid minutes - should be 0-59")
	}
	result := float64(degs) + float64(mins)/60
	if match[1] == "-" {
		result = -result
	}
	return result, nil
}

func formatHH(value int) string {
	sign := ""
	if value < 0 {
		value = -value
		sign = "-"
	}
	return fmt.Sprintf("%s%02d", sign, value)
}

func formatHHMMSS(value float64) string {
	sign := ""
	if value < 0 {
		value = -value
		sign = "-"
	}
	hh := int(math.Floor(value))
	mm := int(math.Floor(value*60)) - hh*60
	ss := value*3600 - float64(hh*3600+mm*60.0)
	ss = math.Round(ss*1000) / 1000
	if ss >= 60 {
		ss -= 60
		mm += 1
	}
	if mm >= 60 {
		mm -= 60
		hh += 1
	}
	return fmt.Sprintf("%s%02d:%02d:%06.3f", sign, hh, mm, ss)
}

func formatDDMM(value float64) string {
	sign := ""
	if value < 0 {
		value = -value
		sign = "-"
	}
	dd := int(math.Floor(value))
	mm := int(math.Floor(value*60)) - dd*60
	return fmt.Sprintf("%s%02d*%02d", sign, dd, mm)
}

func formatDDMMSS(value float64, incSign bool, deg3 bool) string {
	sign := ""
	if value < 0 {
		value = -value
		sign = "-"
	} else if incSign {
		sign = "+"
	}
	dd := int(math.Floor(value))
	mm := int(math.Floor(value*60)) - dd*60
	ss := value*3600 - float64(dd*3600+mm*60.0)
	ss = math.Round(ss*1000) / 1000
	if ss >= 60 {
		ss -= 60
		mm += 1
	}
	if mm >= 60 {
		mm -= 60
		dd += 1
	}
	if deg3 {
		return fmt.Sprintf("%s%03d:%02d:%06.3f", sign, dd, mm, ss)
	} else {
		return fmt.Sprintf("%s%02d:%02d:%06.3f", sign, dd, mm, ss)
	}
}
