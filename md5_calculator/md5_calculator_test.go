package md5_calculator

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculator_CalcMD5Hash(t *testing.T) {
	respStr := "Hello, this is a test server!"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, respStr)
	}))
	defer testServer.Close()

	calculator := Calculator{Client: &http.Client{Timeout: 1 * time.Minute}}

	hash := md5.Sum([]byte(respStr))
	want := hex.EncodeToString(hash[:])

	ctx := context.Background()
	got, err := calculator.CalcMD5Hash(ctx, testServer.URL)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// TODO: test errors
