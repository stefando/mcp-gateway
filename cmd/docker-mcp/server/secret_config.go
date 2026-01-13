package server

import (
	"context"

	"github.com/docker/mcp-gateway/cmd/docker-mcp/secret-management/secret"
	"github.com/docker/mcp-gateway/pkg/config"
	"github.com/docker/mcp-gateway/pkg/desktop"
	"github.com/docker/mcp-gateway/pkg/gateway"
)

// getConfiguredSecretNames returns a map of configured secret names for quick lookup.
// This is a shared helper used by both ls.go and enable.go.
func getConfiguredSecretNames(ctx context.Context) (map[string]struct{}, error) {
	configuredSecretNames := make(map[string]struct{})

	// Check if keychain is the default provider
	if data, err := config.ReadGatewayDefaults(); err == nil && len(data) > 0 {
		if defaults, err := gateway.ParseGatewayDefaults(data); err == nil && defaults.DefaultSecretProvider == secret.Keychain {
			// Use keychain provider
			provider := secret.NewCredStoreProvider()
			secrets, err := provider.ListSecrets()
			if err != nil {
				return nil, err
			}
			for _, name := range secrets {
				configuredSecretNames[name] = struct{}{}
			}
			return configuredSecretNames, nil
		}
	}

	// Default: use Docker Desktop secrets
	secretsClient := desktop.NewSecretsClient()
	configuredSecrets, err := secretsClient.ListJfsSecrets(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range configuredSecrets {
		configuredSecretNames[s.Name] = struct{}{}
	}

	return configuredSecretNames, nil
}
