package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReserveTimes(t *testing.T) {

	start := 11
	duration := 5

	startHr, endHr := reserveTimes(start, duration)

	assert.Equal(t, start, startHr)
	assert.Equal(t, 16, endHr)
}
