package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	svcName 	   = "dns-controller"
	shortDescription = "DNS Controller for managing DNS records"
	longDescription  = `DNS Controller is a tool for managing DNS records through various providers like Cloudflare.`
)

var (
    cfgFile       string
    cloudflareKey string
)

var rootCmd = &cobra.Command{
    Use:   svcName,
    Short: shortDescription,
	Long:  longDescription,
    PersistentPreRun: persistPreRun,
    Run: run,
}

func persistPreRun(cmd *cobra.Command, args []string) {
	// If cloudflareKey wasn't set by flag, try to get it from viper
	if cloudflareKey == "" {
		cloudflareKey = viper.GetString("cloudflare_key")
	}
	
	// Validate required configuration
	if cloudflareKey == "" {
		fmt.Fprintln(os.Stderr, "Error: Cloudflare API key not provided")
		fmt.Fprintln(os.Stderr, "Please set it using --cloudflare-key flag, DNS_CONTROLLER_CLOUDFLARE_KEY environment variable, or in the config file")
		os.Exit(1)
	}
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default same directory as binary)")
    rootCmd.PersistentFlags().StringVar(&cloudflareKey, "cloudflare-key", "", "Cloudflare API key")
}

func initConfig() {
    // Don't forget to bind the flag to viper
    viper.BindPFlag("cloudflare_key", rootCmd.PersistentFlags().Lookup("cloudflare-key"))
    viper.SetDefault("cloudflare_key", "")
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

    // Use config file from the flag if specified
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        // Search config in various locations
        viper.AddConfigPath(".")
    }

    // Read in config
    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    }
}

func run(cmd *cobra.Command, args []string) {
    // Main application logic
    fmt.Println("DNS Controller starting...")

    debugCfg()

    // todo insert logic
}

func debugCfg() {
	// Read the yaml config file
    fmt.Println("Configuration settings:")
    fmt.Println("------------------------")
    
    // Get all settings from viper
    allSettings := viper.AllSettings()
    if len(allSettings) == 0 {
        fmt.Println("No configuration settings found")
		os.Exit(1)
    }
	for key, value := range allSettings {
		if key == "cloudflare_key" {
			continue
		}
		fmt.Printf("%s: %v\n", key, value)
	}
    fmt.Println("------------------------")
    fmt.Printf("Using Cloudflare API key: %s\n", maskAPIKey(cloudflareKey))
}

// Helper function to mask API key for display
func maskAPIKey(key string) string {
    if len(key) <= 4 {
        return "****"
    }
    visible := 4 // Show last 4 characters
    return strings.Repeat("*", len(key)-visible) + key[len(key)-visible:]
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}