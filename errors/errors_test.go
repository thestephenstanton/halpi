package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"
)

func TestErrorsCompatibility(t *testing.T) {
	testCases := []struct {
		desc     string
		actual   error
		expected error
	}{
		{
			desc:     "New",
			actual:   New("some error"),
			expected: errors.New("some error"),
		},
		{
			desc:     "Newf",
			actual:   Newf("some error %s", "with args"),
			expected: fmt.Errorf("some error %s", "with args"),
		},
		{
			desc:     "Wrap",
			actual:   Wrap(New("some error"), "wrapped in some message"),
			expected: errors.Wrap(errors.New("some error"), "wrapped in some message"),
		},
		{
			desc:     "Wrapf",
			actual:   Wrapf(New("some error"), "wrapped in some message %s", "with args"),
			expected: errors.Wrapf(errors.New("some error"), "wrapped in some message %s", "with args"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.EqualError(t, tc.actual, tc.expected.Error())
		})
	}
}
