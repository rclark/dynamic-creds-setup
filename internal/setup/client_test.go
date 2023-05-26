package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFingerprint(t *testing.T) {
	t.Run("app.terraform.io", func(t *testing.T) {
		fingerprint, err := getFingerprint(context.Background(), "app.terraform.io")
		require.NoError(t, err, "should not error")
		assert.Equal(t, "026CC95D81A19F8E4B3A7C15E2D4A9A283A97FF2", fingerprint, "should produce expected fingerprint")
	})
}
