package main

import (
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFoo(t *testing.T) {
	// res, err := foo()

	//require.NoError(t, err)
	//assert.Contains(t, 2, res)

	require.JSONEq(t, `{"a":1,"b":2}`, `{"b":2,"a":1}`)
}
