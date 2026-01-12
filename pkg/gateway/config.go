package gateway

import (
	"github.com/docker/mcp-gateway/pkg/catalog"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Options
	WorkingSet         string
	ServerNames        []string
	CatalogPath        []string
	ConfigPath         []string
	RegistryPath       []string
	ToolsPath          []string
	SecretsPath        string
	MCPRegistryServers []catalog.Server // catalog.Server objects from MCP registries
}

type Options struct {
	Port                    int
	Transport               string
	ToolNames               []string
	Interceptors            []string
	OciRef                  []string
	Verbose                 bool
	LongLived               bool
	DebugDNS                bool
	LogCalls                bool
	BlockSecrets            bool
	BlockNetwork            bool
	VerifySignatures        bool
	DryRun                  bool
	Watch                   bool
	Cpus                    int
	Memory                  string
	Static                  bool
	OAuthInterceptorEnabled bool
	McpOAuthDcrEnabled      bool
	DynamicTools            bool
	ToolNamePrefix          bool
	LogFilePath             string
	UseEmbeddings           bool
	UseProfiles             bool
}

// GatewayDefaults holds persistent configuration for Colima/Docker CE users
// Stored in ~/.docker/mcp/gateway.yaml
type GatewayDefaults struct {
	// SkipDesktopCheck bypasses Docker Desktop checks (for Colima users)
	SkipDesktopCheck bool `yaml:"skipDesktopCheck"`
	// Secrets is the default secrets provider (keychain, docker-desktop, or file path)
	Secrets string `yaml:"secrets"`
	// DefaultSecretProvider is the default provider for secret CLI commands
	DefaultSecretProvider string `yaml:"defaultSecretProvider"`
}

// ParseGatewayDefaults parses gateway defaults from YAML bytes
func ParseGatewayDefaults(data []byte) (*GatewayDefaults, error) {
	if len(data) == 0 {
		return &GatewayDefaults{}, nil
	}

	var defaults GatewayDefaults
	if err := yaml.Unmarshal(data, &defaults); err != nil {
		return nil, err
	}

	return &defaults, nil
}
