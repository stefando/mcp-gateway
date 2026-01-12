package secret

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/mcp-gateway/cmd/docker-mcp/secret-management/formatting"
	"github.com/docker/mcp-gateway/pkg/desktop"
)

type ListOptions struct {
	JSON     bool
	Provider string
}

func List(ctx context.Context, opts ListOptions) error {
	// Handle keychain provider
	if opts.Provider == Keychain {
		p := NewCredStoreProvider()
		names, err := p.ListSecrets()
		if err != nil {
			return err
		}

		if opts.JSON {
			jsonData, err := json.MarshalIndent(names, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonData))
			return nil
		}

		var rows [][]string
		for _, name := range names {
			rows = append(rows, []string{name, "keychain"})
		}
		formatting.PrettyPrintTable(rows, []int{40, 120})
		return nil
	}

	// Default: Docker Desktop
	l, err := desktop.NewSecretsClient().ListJfsSecrets(ctx)
	if err != nil {
		return err
	}

	if opts.JSON {
		if len(l) == 0 {
			l = []desktop.StoredSecret{} // Guarantee empty list (instead of displaying null)
		}
		jsonData, err := json.MarshalIndent(l, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	}
	var rows [][]string
	for _, v := range l {
		rows = append(rows, []string{v.Name, v.Provider})
	}
	formatting.PrettyPrintTable(rows, []int{40, 120})
	return nil
}
