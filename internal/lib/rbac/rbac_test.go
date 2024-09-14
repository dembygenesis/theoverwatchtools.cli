package rbac_test

import (
	"embed"
	"github.com/dembygenesis/local.tools/internal/lib/rbac"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type Sources struct {
	Model    *rbac.Source
	Policies *rbac.Source
}

//go:embed test_model.conf
var model embed.FS

//go:embed test_policy.conf
var policy embed.FS

func testSources(t *testing.T) *Sources {
	return &Sources{
		Model: &rbac.Source{
			Fs:   &model,
			Name: "test_model.conf",
		},
		Policies: &rbac.Source{
			Fs:   &policy,
			Name: "test_policy.conf",
		},
	}
}

func TestNew(t *testing.T) {
	s := testSources(t)

	r, err := rbac.New(s.Model, s.Policies)
	require.NoError(t, err, "unexpected new error")
	require.NotNil(t, r, "unexpected nil enforcer")
}

type testCaseAuthZ struct {
	name       string
	assertions func(t *testing.T, ok bool, err error)
	sub        string
	act        string
}

func getTestCasesAuthZ() []testCaseAuthZ {
	return []testCaseAuthZ{
		{
			name: "ok",
			assertions: func(t *testing.T, ok bool, err error) {
				assert.NoError(t, err, "unexpected error")
				assert.True(t, ok, "unexpected false")
			},
			sub: "admin",
			act: "create-capture-page",
		},
		{
			name: "not-ok",
			assertions: func(t *testing.T, ok bool, err error) {
				assert.NoError(t, err, "unexpected error")
				assert.False(t, ok, "unexpected true")
			},
			sub: "dummy-subject",
			act: "dummy-action",
		},
	}
}

func Test_AuthZ(t *testing.T) {
	s := testSources(t)

	r, err := rbac.New(s.Model, s.Policies)
	require.NoError(t, err, "unexpected new error")
	require.NotNil(t, r, "unexpected nil enforcer")

	for _, tt := range getTestCasesAuthZ() {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := r.AuthZ(tt.sub, tt.act)
			tt.assertions(t, ok, err)
		})
	}
}
