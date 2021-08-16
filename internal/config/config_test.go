package config

import (
	"context"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		desc      string
		filename  string
		wantError bool
	}{
		{
			desc:      "success case",
			filename:  "../../deployment/go-surf/values.yaml",
			wantError: false,
		},
		{
			desc:      "wrong path",
			filename:  "wrong/path.yaml",
			wantError: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			filename, err := filepath.Abs(tC.filename)
			if err != nil {
				t.Fatalf("cannot find config filepath, error: %v", err)
			}
			_, err = Parse(context.Background(), filename)
			if err != nil && !tC.wantError {
				t.Fatalf("got error: %v", err)
			}
		})
	}
}
