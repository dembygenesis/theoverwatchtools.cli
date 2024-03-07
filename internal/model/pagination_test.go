package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPagination_Validate_Success(t *testing.T) {
	p := NewPagination()
	err := p.ValidatePagination()
	require.NoError(t, err, "unexpected error")
}

func TestPagination_Validate_Success_Empty_Values(t *testing.T) {
	p := Pagination{
		Rows: 0,
		Page: 0,
	}
	err := p.ValidatePagination()
	require.NoError(t, err, "unexpected nil error")
}

func TestPagination_Validate_Fail_Fields(t *testing.T) {
	p := Pagination{
		Rows: -1,
		Page: -1,
	}
	err := p.ValidatePagination()
	require.Error(t, err, "unexpected nil error")
}

func TestPagination_Validate_Fail_Exceeded_Max_Rows(t *testing.T) {
	p := Pagination{
		Rows: 1000,
		Page: 1,
	}
	err := p.ValidatePagination()
	require.Error(t, err, "unexpected nil error")
}
