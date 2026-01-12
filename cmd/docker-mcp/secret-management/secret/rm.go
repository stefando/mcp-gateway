package secret

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/mcp-gateway/pkg/desktop"
)

type RmOpts struct {
	All      bool
	Provider string
}

func Remove(ctx context.Context, names []string, opts RmOpts) error {
	// Handle keychain provider
	if opts.Provider == Keychain {
		p := NewCredStoreProvider()
		if opts.All && len(names) == 0 {
			var err error
			names, err = p.ListSecrets()
			if err != nil {
				return err
			}
		}
		var errs []error
		for _, name := range names {
			if err := p.DeleteSecret(name); err != nil {
				errs = append(errs, err)
				fmt.Printf("failed removing secret %s from keychain\n", name)
				continue
			}
			fmt.Printf("removed secret %s from keychain\n", name)
		}
		return errors.Join(errs...)
	}

	// Default: Docker Desktop
	c := desktop.NewSecretsClient()
	if opts.All && len(names) == 0 {
		l, err := c.ListJfsSecrets(ctx)
		if err != nil {
			return err
		}
		for _, secret := range l {
			names = append(names, secret.Name)
		}
	}
	var errs []error
	for _, name := range names {
		if err := c.DeleteJfsSecret(ctx, name); err != nil {
			errs = append(errs, err)
			fmt.Printf("failed removing secret %s\n", name)
			continue
		}
		fmt.Printf("removed secret %s\n", name)
	}
	return errors.Join(errs...)
}
