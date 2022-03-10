package metric

import (
	"errors"
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"regexp"
	"strconv"
	"strings"
)

var (
	reThreshold         = regexp.MustCompile(`(?i)^(\d+)\s*(%|[TGMK]i?B)?$`)
	ErrThresholdInvalid = errors.New("threshold invalid")
)

const (
	KibiByte uint64 = 1024
	MebiByte        = 1024 * KibiByte
	GibiByte        = 1024 * MebiByte
	TebiByte        = 1024 * GibiByte
)

// Metric allows to check a metric for levels specified by its thresholds.
//
// Is currently implemented for byte values in IEC - which is what Nagios plugins use.
type Metric struct {
	Value    uint64
	Warning  uint64
	Critical uint64
	Total    uint64
	Type     string
}

// NewMetric returns a new Metric.
//  TODO Value can not be greater than Total
func NewMetric(value, total uint64) *Metric {
	return &Metric{Value: value, Total: total}
}

// SetWarning parses warning level for free OR used and remembers the absolute value.
//
// Used:
// Total 100 MB; Threshold 80% => If 80 MB used. Returns warning
//
// Free:
// Total 100MB; Threshold 80% => If 20 MB free. Returns warning
func (m *Metric) SetWarning(threshold string) (err error) {
	//  TODO Refactor
	var thresh uint64

	switch m.Type {
	case "used":
		thresh, err = ThresholdUsed(threshold, m.Total)
		if err != nil {
			return fmt.Errorf("warning: %w", err)
		}
	case "free":
		thresh, err = ThresholdFree(threshold, m.Total)
		if err != nil {
			return fmt.Errorf("warning: %w", err)
		}
	default:
		return fmt.Errorf("wrong type, please specify 'used' OR 'free'")
	}

	m.Warning = thresh

	return nil
}

// SetCritical parses critical level for free OR used and remembers the absolute value.
//
// Used:
// Total 100 MB; Threshold 90% => If 90 MB used. Returns critical
//
// Free:
// Total 100MB; Threshold 90% => If 10 MB free. Returns critical
func (m *Metric) SetCritical(threshold string) (err error) {
	//  TODO Refactor
	var thresh uint64

	switch m.Type {
	case "used":
		thresh, err = ThresholdUsed(threshold, m.Total)
		if err != nil {
			return fmt.Errorf("critical: %w", err)
		}
	case "free":
		thresh, err = ThresholdFree(threshold, m.Total)
		if err != nil {
			return fmt.Errorf("critical: %w", err)
		}
	default:
		return fmt.Errorf("wrong type, please specify 'used' OR 'free'")
	}

	m.Critical = thresh

	return nil
}

// Status  returns the Icinga status in perspective to the current value and thresholds.
//
// Free:
// Value <= Warning will result in warning state
// Value <= Critical will result in critical state
//
// Used:
// Value >= Warning will result in warning state
// Value >= Critical will result in critical state
func (m *Metric) Status() int {
	switch m.Type {
	case "used":
		if m.Critical > 0 && m.Value >= m.Critical {
			return check.Critical
		} else if m.Warning > 0 && m.Value >= m.Warning {
			return check.Warning
		} else {
			return check.OK
		}
	case "free":
		if m.Critical > 0 && m.Value <= m.Critical {
			return check.Critical
		} else if m.Warning > 0 && m.Value <= m.Warning {
			return check.Warning
		} else {
			return check.OK
		}
	default:
		return check.Unknown
	}
}

// Perfdata returns a perfdata.Perfdata object for output with a plugin.
//
// Values are scaled down to MB, so they are more readable. And we won't need that much precision.
func (m *Metric) Perfdata(label string) perfdata.Perfdata {
	// TODO If the values are to small, the value will be evaluated as 0
	return perfdata.Perfdata{
		Label: label,
		Value: m.Value / MebiByte,
		//Value: float64(m.Value) / float64(MebiByte),
		Uom:  "MB", // Warning: This should be IEC, but Nagios plugins won't know that.
		Warn: &check.Threshold{Upper: float64(m.Warning / MebiByte)},
		//Warn:  &check.Threshold{Upper: float64(m.Warning) / float64(MebiByte)},
		Crit: &check.Threshold{Upper: float64(m.Critical / MebiByte)},
		//Crit:  &check.Threshold{Upper: float64(m.Critical) / float64(MebiByte)},
		Min: 0,
		Max: m.Total / MebiByte,
		//Max:   float64(m.Total) / float64(MebiByte),
	}
}

// ThresholdFree returns the threshold level relative to the total.
//
//  10% free from 100MB = 90MB used
//  15MB free from 100MB = 85MB used
func ThresholdFree(threshold string, total uint64) (uint64, error) {
	level, err := ParseThreshold(threshold, total)
	if err != nil {
		return 0, err
	}

	return total - level, nil
}

// ThresholdUsed returns the threshold level relative to the total.
//
// 90% used from 100MB = 10MB free (available)
// 85MB used from 100MB = 15MB free (available)
func ThresholdUsed(threshold string, total uint64) (uint64, error) {
	level, err := ParseThreshold(threshold, total)
	if err != nil {
		return 0, err
	}

	return level, nil
}

// ParseThreshold returns the parsed unit(UOM) from threshold
//
// Total = 100MB; Threshold = 10% => 10MB
// Total = 100MB; Threshold = 30MB => 30MB
// Viable units are: kb, kib, mb, mib, gb, gib, tb, tib, %
func ParseThreshold(threshold string, total uint64) (uint64, error) {
	match := reThreshold.FindStringSubmatch(threshold)
	if match == nil {
		return 0, ErrThresholdInvalid
	}

	value, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return 0, err
	}

	if match[2] == "%" {
		if value > 100 {
			return 0, fmt.Errorf("percentage can't be larger than 100")
		}

		level := (float64(value) / 100) * float64(total)

		return uint64(level), nil
	}

	var level uint64

	switch u := strings.ToLower(match[2]); u {
	case "kb", "kib":
		level = value * KibiByte
	case "", "mb", "mib":
		level = value * MebiByte
	case "gb", "gib":
		level = value * GibiByte
	case "tb", "tib":
		level = value * TebiByte
	default:
		return 0, fmt.Errorf("invalid unit")
	}

	return level, nil
}
