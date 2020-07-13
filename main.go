package main

import (
	"flag"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/tommyblue/trello-templatizer/trello"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Auth struct {
		Key   string
		Token string
	}
}

var flagVerbose = flag.Bool("v", false, "verbose logging (debug level)")

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)
}

func main() {
	flag.Parse()
	if *flagVerbose {
		log.SetLevel(log.DebugLevel)
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal("errors opening config: ", err)
	}

	conf := parseConfig(f)

	api, err := trello.New(conf.Auth.Key, conf.Auth.Token)
	if err != nil {
		log.Fatalf("Can't initialize: %v", err)
	}
	log.Debugf("%+v", api)
}

func parseConfig(f io.Reader) *Config {
	conf := Config{}
	d := yaml.NewDecoder(f)
	if err := d.Decode(&conf); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Debugf("%+v", conf)
	return &conf
}
