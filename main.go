package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

//go:generate go run gen/main.go

var peers *string
var readOnly *bool

func main() {
	peers = flag.String("peers", "peers", "file containing a list of peers to connect to")
	ipeer := flag.String("peer", "http://git.idk.i2p/idk/db", "initial I2P peer to connect to")
	readOnly = flag.Bool("readOnly", false, "read only mode")
	flag.Parse()
	if _, err := os.Stat(*peers); err != nil {
		if err := ioutil.WriteFile(*peers, []byte(*ipeer), 0644); err != nil {
			log.Fatal(err)
		}
	}

	cloner()
	if err := EnsureGitDBInitialized("db/db"); err != nil {
		log.Fatal(err)
	}
	if err := AddAllGitDB("db/db"); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)
	go HostARemote("db")
	go updater()
	ServeWWW("db/db")
}

func cloner() {
	if *peers != "" {
		peersListBytes, err := ioutil.ReadFile(*peers)
		if err != nil {
			panic(err)
		}
		peersList := strings.Split(string(peersListBytes), "\n")
		for _, peer := range peersList {
			if peer != "" {
				CloneARemote(peer, "db/db")
				time.Sleep(time.Second * 1)
				log.Println("cloned in repo from", peer)
				break
			}
		}
	}
}

func updater() {
	if *peers != "" {
		for {
			peersListBytes, err := ioutil.ReadFile(*peers)
			if err != nil {
				panic(err)
			}
			peersList := strings.Split(string(peersListBytes), "\n")
			for _, peer := range peersList {
				if peer != "" {
					PullInCommitsFromRemote(peer, "db/db")
					time.Sleep(time.Second * 1)
					log.Println("pulled in commits from", peer)
				}
			}
			time.Sleep(time.Second * 120)
		}
	}
}
