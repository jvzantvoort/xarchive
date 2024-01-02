package main

import (
	"os"

	"github.com/jvzantvoort/xarchive/database"
	"github.com/jvzantvoort/xarchive/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool

	Username string
	Password string
	Database string
	Hostname string
	Port     int

	DB *database.Database
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xarchive",
	Short: messages.GetUsage("root"),
	Long:  messages.GetLong("root"),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Setup logging
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	cobra.OnInitialize(initConfig)
	initConfig()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xarchive.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")

	DB = database.NewDatabase(Username, Password, Hostname, Database, Port)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.Debugf("initConfig start")
	defer log.Debugf("initConfig end")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".xarchive" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".xarchive")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Using config file: %s", viper.ConfigFileUsed())
		log.Error(err)
	} else {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}

	Username = viper.GetString("database.username")
	Password = viper.GetString("database.password")
	Hostname = viper.GetString("database.hostname")
	Database = viper.GetString("database.database")
	Port = viper.GetInt("database.port")
}
