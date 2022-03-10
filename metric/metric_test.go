package metric

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetric(t *testing.T) {
	m := NewMetric(10*MebiByte, 100*MebiByte)

	assert.Equal(t, m.Value, 10*MebiByte)
	assert.Equal(t, m.Total, 100*MebiByte)
}

func TestSimpleMetric_SetWarning(t *testing.T) {
	m := &Metric{Value: 10 * MebiByte, Total: 100 * MebiByte, Type: "free"}

	err := m.SetWarning("10%")
	assert.NoError(t, err)
	assert.Equal(t, 90*MebiByte, m.Warning)

	m = &Metric{Value: 90 * MebiByte, Total: 100 * MebiByte, Type: "used"}

	err = m.SetWarning("90%")
	assert.NoError(t, err)
	assert.Equal(t, 90*MebiByte, m.Warning)
}

func TestSimpleMetric_SetCritical(t *testing.T) {
	m := &Metric{Value: 10 * MebiByte, Total: 100 * MebiByte, Type: "free"}

	err := m.SetCritical("20%")
	assert.NoError(t, err)
	assert.Equal(t, 80*MebiByte, m.Critical)

	m = &Metric{Value: 80 * MebiByte, Total: 100 * MebiByte, Type: "used"}

	err = m.SetCritical("80%")
	assert.NoError(t, err)
	assert.Equal(t, 80*MebiByte, m.Critical)
}

// nolint: dupl
func TestThresholdFree(t *testing.T) {
	threshold, err := ThresholdFree("10%", 100*MebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 90*MebiByte, threshold)

	threshold, err = ThresholdFree("20MB", 100*MebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 80*MebiByte, threshold)

	threshold, err = ThresholdFree("10", 100*MebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 90*MebiByte, threshold)

	threshold, err = ThresholdFree("25MiB", 100*MebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 75*MebiByte, threshold)

	threshold, err = ThresholdFree("25GiB", 100*GibiByte)
	assert.NoError(t, err)
	assert.Equal(t, 75*GibiByte, threshold)

	threshold, err = ThresholdFree("25TiB", 100*TebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 75*TebiByte, threshold)

	_, err = ThresholdFree("101%", 100*MebiByte)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("percentage can't be larger than 100"), err)
	}

	_, err = ThresholdFree("50Exmaple", 100*MebiByte)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("threshold invalid"), err)
	}
}

// nolint: dupl
func TestThresholdUsed(t *testing.T) {
	threshold, err := ThresholdUsed("90%", 100*MebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 90*MebiByte, threshold)

	threshold, err = ThresholdUsed("80MB", 100*MebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 80*MebiByte, threshold)

	threshold, err = ThresholdUsed("90", 100*MebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 90*MebiByte, threshold)

	threshold, err = ThresholdUsed("75MiB", 100*MebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 75*MebiByte, threshold)

	threshold, err = ThresholdUsed("75GiB", 100*GibiByte)
	assert.NoError(t, err)
	assert.Equal(t, 75*GibiByte, threshold)

	threshold, err = ThresholdUsed("75TiB", 100*TebiByte)
	assert.NoError(t, err)
	assert.Equal(t, 75*TebiByte, threshold)

	_, err = ThresholdUsed("101%", 100*MebiByte)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("percentage can't be larger than 100"), err)
	}

	_, err = ThresholdUsed("50Exmaple", 100*MebiByte)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("threshold invalid"), err)
	}
}

func TestSimpleMetric_StatusFree(t *testing.T) {
	m := &Metric{Value: 30 * MebiByte, Total: 100 * MebiByte, Type: "free"}
	err := m.SetCritical("90%") // 10 MB available space
	assert.NoError(t, err)

	err = m.SetWarning("80%") // 20 MB available space
	assert.NoError(t, err)
	assert.Equal(t, check.OK, m.Status())

	m.Value = 20 * MebiByte
	assert.Equal(t, check.Warning, m.Status())

	m.Value = 10 * MebiByte
	assert.Equal(t, check.Critical, m.Status())
}

func TestSimpleMetric_StatusUsed(t *testing.T) {
	m := &Metric{Value: 79 * MebiByte, Total: 100 * MebiByte, Type: "used"}
	err := m.SetCritical("90%") // 10 MB available space
	assert.NoError(t, err)

	err = m.SetWarning("80%") // 20 MB available space
	assert.NoError(t, err)

	assert.Equal(t, check.OK, m.Status())

	m.Value = 80 * MebiByte
	assert.Equal(t, check.Warning, m.Status())

	m.Value = 90 * MebiByte
	assert.Equal(t, check.Critical, m.Status())
}

func TestSimpleMetric_Perfdata(t *testing.T) {
	m := &Metric{Value: 10 * MebiByte, Total: 100 * MebiByte, Type: "free"}
	err := m.SetCritical("10%") // 10 MB
	assert.NoError(t, err)

	err = m.SetWarning("20%") // 20 MB
	assert.NoError(t, err)

	assert.Equal(t, "/=10MB;80;90;0;100", m.Perfdata("/").String())

	m = &Metric{Value: 90 * MebiByte, Total: 100 * MebiByte, Type: "used"}
	err = m.SetCritical("90%")
	assert.NoError(t, err)

	err = m.SetWarning("80%")
	assert.NoError(t, err)

	assert.Equal(t, "/=90MB;80;90;0;100", m.Perfdata("/").String())
}
