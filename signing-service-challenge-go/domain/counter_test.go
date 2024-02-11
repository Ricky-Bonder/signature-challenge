package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSignatureCounter(t *testing.T) {
	assert.Equal(t, int32(0), Increment().counter)
	counter := Increment()
	assert.Equal(t, int32(1), counter.counter)
	assert.Equal(t, int32(1), counter.Get())
}
