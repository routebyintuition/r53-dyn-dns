package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/viper"
)

// Config hold the configuration settings passed from viper
type Config struct {
	Interval     int
	Profile      string
	Hostname     string
	ZoneID       string
	LogDirectory string
	DNSURL       string
	Route53      *route53.Route53
	Done         chan os.Signal
}

var configFile = flag.String("config", "config.toml", "Path and file name of configuration TOML file")

func main() {

	var conf Config

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	viper.SetConfigFile(*configFile)
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath.Dir(*configFile))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))

	}

	conf.Hostname = viper.GetString("hostname")
	conf.DNSURL = viper.GetString("dns_url")
	conf.LogDirectory = viper.GetString("server.log_directory")
	conf.Interval = viper.GetInt("server.refresh_interval")
	conf.Profile = viper.GetString("aws.aws_profile")
	conf.ZoneID = viper.GetString("aws.hosted_zone_id")

	conf.logInit()

	Info.Println("Initializing service...")

	err = conf.initAWS()

	_, err = conf.GetPublicIPService()
	if err != nil {
		Error.Println("Could not get public IP: ", err)
		os.Exit(1)
	}

	_, err = conf.getRecord()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	_, err = conf.GetRecords()
	if err != nil {
		Error.Println("Error in Route 53 query: ", err)
	}

	ticker := time.NewTicker(time.Duration(conf.Interval) * time.Second)

	conf.Done = make(chan os.Signal)
	signal.Notify(conf.Done, os.Interrupt)

	err = conf.Process(time.Now())
	if err != nil {
		Error.Println("Error during startup: ", err)
		os.Exit(1)
	}

	Info.Println("Initialization complete")

	for {
		select {
		case <-conf.Done:
			Info.Println("Server stopping due to interrupt signal...")
			return
		case changeDate := <-ticker.C:
			Info.Println("Running update... ")
			conf.Process(changeDate)
		}
	}
}
