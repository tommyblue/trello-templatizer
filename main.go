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
	Template struct {
		BoardID   string `yaml:"board_id,omitempty"`
		BoardName string `yaml:"board_name,omitempty"`
		Lists     []struct {
			Name string // Name of the List to create things to
		}
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
		log.SetLevel(log.InfoLevel)
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

	var board *trello.Board
	if conf.Template.BoardID == "" {
		if board, err = api.SearchBoardByName(conf.Template.BoardName); err != nil {
			log.Fatal(err)
		}
	} else {
		if board, err = api.SearchBoardByID(conf.Template.BoardID); err != nil {
			log.Fatal(err)
		}
	}

	log.Infof("Found board \"%s\" (ID: %s)\n", board.Name, board.ID)

	// Look for the list to create the cards into
	for _, l := range conf.Template.Lists {
		var list *trello.List
		list, err = api.SearchListByName(board.ID, l.Name)
		if err != nil {
			switch err.(type) {
			case trello.ErrListNotFound:
				log.Infof("List \"%s\" not found, creating", l.Name)
				list, err = api.CreateList(board.ID, l.Name)
				if err != nil {
					log.Fatal(err)
				}
			default:
				log.Fatal(err)
			}
		}
		log.Infof("Found list \"%s\" (ID: %s)\n", list.Name, list.ID)

		if list.Closed {
			log.Errorf("Cannot write in a closed list: %s (ID: %s)", list.Name, list.ID)
			continue
		}

	}

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
