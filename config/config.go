package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// ProxyConfig stores configuration specific to the Go proxy itself
type ProxyConfig struct {
	ComPort       string     `json:"comPort"`
	BaudRate      int        `json:"baudRate"`
	NetworkPort   int        `json:"networkPort"`
	ListenAddress string     `json:"listenAddress"`
	LogLevel      slog.Level `json:"logLevel"`
	SiteLatitude  float64    `json:"sitelatitude"`
	SiteLongitude float64    `json:"sitelongitude"`
	SiteElevation float64    `json:"siteelevation"`
	SiteUTCOffset int        `json:"siteutcoffset"`
}

var (
	proxyConfig     *ProxyConfig // Singleton instance
	proxyConfigFile string       // Full path to the config file
)

// init sets up the path to the configuration file.
func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// This is a critical failure at startup. We can't proceed without a config path.
		// Using log.Fatalf here is acceptable as it's a pre-flight check.
		slog.Error("FATAL: Could not get user config directory", "error", err)
		os.Exit(1)
	}
	appConfigDir := filepath.Join(configDir, "OnStepXAlpacaProxy")
	// The logger setup will create this dir, but it's safe to do it here too.
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		slog.Error("FATAL: Could not create application config directory", "directory", appConfigDir, "error", err)
		os.Exit(1)
	}
	proxyConfigFile = filepath.Join(appConfigDir, "proxy_config.json")
}

// Load reads the configuration from the JSON file into the singleton instance.
// If the file doesn't exist, it initializes a default configuration and saves it.
func Load() error {
	file, err := os.ReadFile(proxyConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("Proxy config file not found. Using default settings.", "configFile", proxyConfigFile)
			// Initialize with default values
			proxyConfig = &ProxyConfig{
				ComPort:       "",
				BaudRate:      9600,
				NetworkPort:   32241,
				ListenAddress: "127.0.0.1",
				LogLevel:      slog.LevelInfo,
			}
			// Attempt to save the initial default config
			return Save() // File not found is not an error, just means defaults apply
		}
		return fmt.Errorf("failed to read proxy config file: %w", err)
	}

	// First, unmarshal into a temporary instance.
	var tempConfig ProxyConfig
	if err := json.Unmarshal(file, &tempConfig); err != nil {
		// Don't overwrite the global config if unmarshalling fails.
		return fmt.Errorf("failed to unmarshal proxy config: %w", err)
	}
	proxyConfig = &tempConfig

	// --- Validate and set defaults for missing fields ---
	if proxyConfig.BaudRate == 0 {
		proxyConfig.BaudRate = 9600
	}
	if proxyConfig.NetworkPort == 0 {
		proxyConfig.NetworkPort = 32241
	}
	if proxyConfig.ListenAddress == "" {
		slog.Warn("Configuration key 'ListenAddress' not found, using default '127.0.0.1'.")
		proxyConfig.ListenAddress = "127.0.0.1"
	}

	// Apply the loaded log level immediately.
	//slog.SetLevelFromString(proxyConfig.LogLevel)
	slog.Info("Loaded proxy config", "configFile", proxyConfigFile)
	return nil
}

// Save writes the current configuration to the JSON file.
func Save() error {
	if proxyConfig == nil {
		return fmt.Errorf("cannot save nil config")
	}
	slog.Debug("Attempting to save proxy config to file", "configFile", proxyConfigFile)
	data, err := json.MarshalIndent(proxyConfig, "", "  ")
	if err != nil {
		slog.Error("saveProxyConfig: failed to marshal proxy config", "error", err)
		return fmt.Errorf("failed to marshal proxy config: %w", err)
	}

	if err := os.WriteFile(proxyConfigFile, data, 0644); err != nil {
		slog.Error("saveProxyConfig: failed to write proxy config file", "configFile", proxyConfigFile, "error", err)
		return fmt.Errorf("failed to write proxy config file: %w", err)
	}
	slog.Info("Successfully saved proxy config to file", "configFile", proxyConfigFile)
	return nil
}

// Get returns a pointer to the singleton ProxyConfig instance.
func Get() *ProxyConfig {
	if proxyConfig == nil {
		// This should not happen in the normal flow, as Load() is called on startup.
		// But as a safeguard, we initialize a default config.
		if err := Load(); err != nil {
			slog.Error("Failed to load configuration on demand", "error", err)
			os.Exit(1)
		}
	}
	return proxyConfig
}

// GetSetupURL builds the full URL for the web setup page based on the current config.
func GetSetupURL() string {
	conf := Get()
	host := conf.ListenAddress
	if host == "0.0.0.0" || host == "" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s:%d/setup", host, conf.NetworkPort)
}

// GetSetupURLFromFile reads the configuration file directly to build the setup URL.
// This is a special case for the single-instance check, which runs before the main
// configuration and logging are initialized. It ensures that a second instance
// opens the correct URL based on the saved listenAddress.
func GetSetupURLFromFile() string {
	const defaultHost = "127.0.0.1"
	const defaultPort = 32241

	file, err := os.ReadFile(proxyConfigFile)
	if err != nil {
		// File not found or other error, use failsafe defaults.
		return fmt.Sprintf("http://%s:%d/setup", defaultHost, defaultPort)
	}

	var config struct {
		NetworkPort   int    `json:"networkPort"`
		ListenAddress string `json:"listenAddress"`
	}
	if err := json.Unmarshal(file, &config); err != nil {
		// JSON is corrupt, use failsafe defaults.
		return fmt.Sprintf("http://%s:%d/setup", defaultHost, defaultPort)
	}

	host := config.ListenAddress
	port := config.NetworkPort

	if host == "0.0.0.0" || host == "" {
		host = defaultHost
	}
	if port == 0 {
		port = defaultPort
	}

	return fmt.Sprintf("http://%s:%d/setup", host, port)
}
