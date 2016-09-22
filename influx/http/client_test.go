package http

import (
	"context"
	"testing"

	"github.com/stellar/go/support/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestClient_Get(t *testing.T) {
	// store, mock := influxtest.NewService()

	client := &Client{
		Raw: httptest.NewClient(),
	}

	_, err := client.Get(context.Background(), "https://google.com")
	assert.Error(t, err)
}
