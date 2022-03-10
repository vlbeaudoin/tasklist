package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

func declareFlags() {
	// db.path
	rootCmd.PersistentFlags().String(
		"db-path", "",
		"Database path (config: 'db.path')")
	viper.BindPFlag(
		"db.path",
		rootCmd.PersistentFlags().Lookup("db-path"))
	rootCmd.MarkPersistentFlagRequired("db.path")

	// db.type
	rootCmd.PersistentFlags().String(
		"db-type", "",
		"Database type (config: 'db.type')")
	viper.BindPFlag(
		"db.type",
		rootCmd.PersistentFlags().Lookup("db-type"))
	rootCmd.MarkPersistentFlagRequired("db.type")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tasklist",
	Short: "Maintain a list of tasks.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tasklist.yaml)")

	declareFlags()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tasklist" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tasklist")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
