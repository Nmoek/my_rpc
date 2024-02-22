//go:build e2e

package v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_Start(t *testing.T) {

	s := &Server{}

	err := s.Start("tcp", ":8888")
	assert.NoError(t, err)
}
