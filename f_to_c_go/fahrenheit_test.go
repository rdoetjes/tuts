package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConformLine(t *testing.T) {
	// Test valid Fahrenheit input
	fahrenheit := "98.6 F"
	expected := "37.0 C"
	actual, err := conform_line(fahrenheit)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Test invalid Fahrenheit input
	fahrenheit = "abc F"
	expected = "abc F"
	actual, err = conform_line(fahrenheit)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Test missing Fahrenheit unit
	fahrenheit = "98.6"
	actual, err = conform_line(fahrenheit)
	assert.NoError(t, err)
	assert.Equal(t, fahrenheit, actual)
}
