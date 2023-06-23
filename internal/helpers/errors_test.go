package helpers

import (
	"errors"
	"testing"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/stretchr/testify/assert"
)

func TestParseClientError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		action   string
		name     string
		err      error
		expected string
	}{
		"openapi": {
			action:   "create",
			name:     "prowlarr_tag",
			err:      &prowlarr.GenericOpenAPIError{},
			expected: "Unable to create prowlarr_tag, got error: \nDetails:\n",
		},
		"generic": {
			action:   "create",
			name:     "prowlarr_tag",
			err:      errors.New("other error"),
			expected: "Unable to create prowlarr_tag, got error: other error",
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ParseClientError(test.action, test.name, test.err))
		})
	}
}

func TestParseNotFoundError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		kind     string
		field    string
		search   string
		expected string
	}{
		"generic": {
			kind:     "prowlarr_tag",
			field:    "label",
			search:   "test",
			expected: "Unable to find prowlarr_tag, got error: data source not found: no prowlarr_tag with label 'test'",
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ParseNotFoundError(test.kind, test.field, test.search))
		})
	}
}
