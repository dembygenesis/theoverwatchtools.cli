package rbac

import (
	"embed"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"os"
)

type Rbac struct {
	enforcer *casbin.Enforcer
}

type Source struct {
	Fs   *embed.FS `validate:"required"`
	Name string    `validate:"required"`
}

func (r *Rbac) Load() error {
	return nil
}

func (r *Rbac) AuthZ(sub, act string) (bool, error) {
	if ok, err := r.enforcer.Enforce(sub, act); err != nil {
		return false, fmt.Errorf("enforce: %w", err)
	} else {
		return ok, nil
	}
}

func fsEmbedToModelString(f *embed.FS, name string) (*model.Model, error) {
	modelContent, err := f.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read model file: %w", err)
	}

	m, err := model.NewModelFromString(string(modelContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	return &m, nil
}

func fsEmbedToPolicyAdapter(f *embed.FS, name string) (persist.Adapter, error) {
	policyContent, err := f.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %w", err)
	}

	tmpFile, err := os.CreateTemp("/tmp", "policy_*.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		err = tmpFile.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	if _, err := tmpFile.Write(policyContent); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}

	adapter := fileadapter.NewAdapter(tmpFile.Name())
	return adapter, nil
}

func New(model *Source, policy *Source) (*Rbac, error) {
	if model == nil {
		return nil, fmt.Errorf("model is nil")
	}

	if policy == nil {
		return nil, fmt.Errorf("policy is nil")
	}

	m, err := fsEmbedToModelString(model.Fs, model.Name)
	if err != nil {
		return nil, fmt.Errorf("parse model: %w", err)
	}

	a, err := fsEmbedToPolicyAdapter(policy.Fs, policy.Name)
	if err != nil {
		return nil, fmt.Errorf("parse adapter: %w", err)
	}

	e, err := casbin.NewEnforcer(*m, a)
	if err != nil {
		return nil, fmt.Errorf("new casbin: %w", err)
	}

	return &Rbac{
		enforcer: e,
	}, nil
}
