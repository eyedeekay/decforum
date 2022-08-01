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

func main() {
	peers = flag.String("peers", "peers", "file containing a list of peers to connect to")
	flag.Parse()
	if _, err := os.Stat(*peers); err != nil {
		if err := ioutil.WriteFile(*peers, []byte(""), 0644); err != nil {
			log.Fatal(err)
		}
	}

	cloner()
	EnsureGitDBInitialized("db/db")
	AddAllGitDB("db/db")

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
			time.Sleep(time.Second * 300)
		}
	}
}
