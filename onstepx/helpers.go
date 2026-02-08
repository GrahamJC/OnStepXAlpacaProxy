package onstepx

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

func zeroPad(value int, width int, incSign bool) string {
	sign := ""
	if value >= 0 {
		if incSign {
			sign = "+"
		}
	} else {
		if incSign {
			sign = "-"
		}
		value = -value
	}
	return fmt.Sprintf("%s%0*d", sign, width, value)
}

var reHHMMSS = regexp.MustCompile(`(\d+):(\d+):(\d+)`)

func parseHHMMSS(s string) (float32, error) {
	var err error
	match := reHHMMSS.FindStringSubmatch(s)
	if match == nil {
		return 0.0, errors.New("Invalid format - should be HH:MM:SS")
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
	secs, err := strconv.Atoi(match[3])
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing seconds - %w", err)
	}
	if secs < 0 || secs > 59 {
		return 0.0, errors.New("Invalid seconds - should be 0-59")
	}
	return float32(hrs) + float32(mins)/60 + float32(secs)/3600, nil
}

var reHHMM = regexp.MustCompile(`(\d+):(\d+)`)

func parseHHMM(s string) (float32, error) {
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
	if mins < 0 || mins > 59 {
		return 0.0, errors.New("Invalid minutes - should be 0-59")
	}
	return float32(hrs) + float32(mins)/60, nil
}

var reDDMMSS = regexp.MustCompile(`^([\+-]{0,1})(\d+)\*(\d+)\:(\d+)$`)

func parseDDMMSS(s string) (float32, error) {
	var err error
	match := reDDMMSS.FindStringSubmatch(s)
	if match == nil {
		return 0.0, errors.New("Invalid format - should be DDD*MM:SS")
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
	secs, err := strconv.Atoi(match[4])
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing seconds - %w", err)
	}
	if secs < 0 || secs > 59 {
		return 0.0, errors.New("Invalid seconds - should be 0-59")
	}
	result := float32(degs) + float32(mins)/60 + float32(secs)/3600
	if match[1] == "-" {
		result = -result
	}
	return result, nil
}

var reDDMM = regexp.MustCompile(`^([\+-]{0,1})(\d+)\*(\d+)$`)

func parseDDMM(s string) (float32, error) {
	var err error
	match := reDDMM.FindStringSubmatch(s)
	if match == nil {
		return 0.0, errors.New("invalid format - should be DDD*MM:SS")
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
	if mins < 0 || mins > 59 {
		return 0.0, errors.New("invalid minutes - should be 0-59")
	}
	result := float32(degs) + float32(mins)/60
	if match[1] == "-" {
		result = -result
	}
	return result, nil
}

func formatDDMM(value float32, incSign bool) string {
	dd := int(math.Trunc(float64(value)))
	mm := int(math.Trunc(float64(value-float32(dd)) * 60))
	return zeroPad(dd, 2, incSign) + "*" + zeroPad(mm, 2, false)
}

func formatDDDMM(value float32) string {
	dd := int(math.Trunc(float64(value)))
	mm := int(math.Trunc(float64(value-float32(dd)) * 60))
	return zeroPad(dd, 3, false) + "*" + zeroPad(mm, 2, false)
}

var reHH = regexp.MustCompile(`^([\+-]{0,1})(\d+)$`)

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

func formatHH(value int, incSign bool) string {
	return zeroPad(value, 2, incSign)
}
