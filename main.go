package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/freman/scantp/driver"
	"goftp.io/server"
)

func main() {
	configFile := flag.String("config", "config.toml", "path to configuration file")
	flag.Parse()

	if *configFile == "" {
		log.Println("Configuration file is required")
		flag.Usage()
		os.Exit(1)
	}

	md, err := toml.DecodeFile(*configFile, &config)
	if err != nil {
		log.Println("Unable to parse configuration file")
		log.Fatal(err)
	}

	mdf := &driver.MultipleDriverFactory{}

	for pathName, prim := range config.Paths {
		var tmp struct {
			Type string `toml:"type"`
		}

		if err := md.PrimitiveDecode(prim, &tmp); err != nil {
			log.Printf("Unable to parse configuration file for %s", pathName)
			log.Fatal(err)
		}

		fmt.Println("Configuring", tmp.Type, pathName)

		if err := mdf.AddPath(pathName, tmp.Type, md, prim); err != nil {
			log.Printf("Unable to start %s for %s", tmp.Type, pathName)
			log.Fatal(err)
		}
	}

	ftpServer := server.NewServer(&server.ServerOpts{
		Factory:  mdf,
		Hostname: config.Host,
		Port:     config.Port,
		Auth:     config,
	})

	if err := ftpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
