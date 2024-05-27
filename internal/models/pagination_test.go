package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPagination_Validate_Success(t *testing.T) {
	p := NewPagination()
	err := p.Validate()
	require.NoError(t, err, "unexpected error")
}

func TestPagination_Validate_Fail_Fields(t *testing.T) {
	p := Pagination{
		Rows: -1,
		Page: -1,
	}
	err := p.Validate()
	require.Error(t, err, "unexpected nil error")
}

func TestPagination_Validate_Fail_Exceeded_Max_Rows(t *testing.T) {
	p := Pagination{
		Rows: 1000,
		Page: 1,
	}
	err := p.Validate()
	require.Error(t, err, "unexpected nil error")
}
