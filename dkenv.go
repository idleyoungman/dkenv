package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

const version = "0.0.1"

func main() {
	viper.SetConfigName("config")
	viper.SetDefault("BinDir", "/usr/local/bin")
	viper.AddConfigPath("$HOME/.dkenv")
	viper.ReadInConfig()

	version := kingpin.Flag("version", "Set Docker version").Short('v').String()
	list := kingpin.Flag("list", "List downloaded Docker versions").Short('l').Bool()
	apiVersion := kingpin.Flag("apiVersion", "Set Docker API version").Short('a').String()
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	var ver string

	if *list {
		fmt.Println("Versions downloaded:")
		listDownloadedVersions()
		os.Exit(0)
	}

	if len(*apiVersion) > 0 || len(*version) > 0 {
		fmt.Println("version has value ", *version)
		fmt.Println("apiVersion has value ", *apiVersion)
		if len(*apiVersion) > 0 {
			var err error
			ver, err = apiToVersion(*apiVersion)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("For apiVersion ", *apiVersion, " using version ", ver)
		} else {
			ver = string(*version)
		}

		if versionDownloaded(ver) {

		} else {
			getDocker(ver, viper.GetString("BinDir"))
		}
		switchVersion(ver, viper.GetString("BinDir"))

	} else {
		kingpin.FatalUsage("Must specify a flag")
	}
}
