package inmem

import (
	"testing"

	"github.com/matryer/vice"
	"github.com/matryer/vice/vicetest"
)

func TestTransport(t *testing.T) {
	new := func() vice.Transport {
		return New()
	}
	vicetest.Transport(t, new)
}
