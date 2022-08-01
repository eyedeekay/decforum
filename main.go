package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

//go:generate go run gen/main.go

var peers *string

func main() {
	peers = flag.String("peers", "", "file containing a list of peers to connect to")
	flag.Parse()
	if *peers == "" {
		EnsureGitDBInitialized("db/db")
		AddAllGitDB("db/db")
	} else {
		cloner()
	}
	time.Sleep(time.Second * 5)
	go HostARemote("db")
	go updater()
	ServeWWW("db/db")
}

func cloner() {
	if *peers != "" {
		for {
			peersListBytes, err := ioutil.ReadFile(*peers)
			if err != nil {
				panic(err)
			}
			peersList := strings.Split(string(peersListBytes), "\n")
			for _, peer := range peersList {
				CloneARemote(peer, "db/db")
				time.Sleep(time.Second * 1)
				log.Println("cloned in repo from", peer)
				break
			}
			break
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
