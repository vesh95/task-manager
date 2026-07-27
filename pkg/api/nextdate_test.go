package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextDate(t *testing.T) {
	values := []struct {
		dstart string
		args   string
		want   string
	}{
		{"20240229", "y", "20250301"},
		{"20240113", "d 7", "20240127"},
		{"20240116", "m 16,5", "20240205"},
		{"20240201", "m -1,18", "20240218"},
		{"20240201", "w 5", "20240202"},
		{"20240125", "w 1", "20240129"},
	}

	now, err := time.Parse(DateFormat, "20240126")

	require.NoError(t, err)

	for i, v := range values {
		t.Run(fmt.Sprintf("Test case %d (%s)", i+1, v.args), func(t *testing.T) {
			date, err := NextDate(now, v.dstart, v.args)
			require.NoError(t, err)
			assert.Equal(t, v.want, date)
		})
	}
}
