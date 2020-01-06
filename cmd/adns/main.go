package main

import (
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "adns",
		Short: "adns is a DoH proxy that only allows the use of Adobe domain names",
		Long:  `adns is a DoH proxy that only allows the use of Adobe domain names`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	cobra.MousetrapHelpText = "" // remove trap as we use serve as default
}

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	// default to serve
	if len(os.Args[1:]) == 0 {
		args := append([]string{"serve"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	flag.Parse()
	err := rootCmd.Execute()
	if err != nil {
		glog.Error(err)
	}
}
