package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	sam "github.com/eyedeekay/sam3/helper"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

//go:embed www/index.html
//go:embed www
var content embed.FS

func randId() string {
	// generate a random number between 1000-9999
	id := rand.Intn(9000) + 1000
	return fmt.Sprintf("%d", id)
}

func ServeWWW(db string) {
	listener, err := sam.I2PListener("gitforum-"+randId(), "127.0.0.1:7656", "gitforum")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("defaulting to index.html")
		indexBytes, err := content.Open("www/index.html")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		w.Write(index)

	}))
	http.Handle("/easymde.min.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving easymde.min.js")
		indexBytes, err := content.Open("www/easymde/easymde.min.js")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "application/javascript")
		w.Write(index)

	}))
	http.Handle("/easymde.min.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving easymde.min.css")
		indexBytes, err := content.Open("www/easymde/easymde.min.css")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "text/css")
		w.Write(index)

	}))
	http.Handle("/home.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving home.css")
		indexBytes, err := content.Open("www/home.css")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "text/css")
		w.Write(index)
	}))
	http.Handle("/font-awesome.min.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving easymde.min.css")
		indexBytes, err := content.Open("www/font-awesome.min.css")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "text/css")
		w.Write(index)

	}))
	http.Handle("/fonts/font-awesome.min.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving easymde.min.css")
		indexBytes, err := content.Open("www/font-awesome.min.css")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "text/css")
		w.Write(index)

	}))
	http.Handle("/fonts/fontawesome-webfont.woff2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving fontawesome-webfont.woff2")
		indexBytes, err := content.Open("www/fontawesome-webfont.woff2")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		//w.Header().Add("Content-Type", "text/css")
		w.Header().Add("Content-Type", "application/font-woff2")
		w.Write(index)
	}))
	http.Handle("/en_US.aff", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving en_US.aff")
		indexBytes, err := content.Open("www/en_US.aff")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		//w.Header().Add("Content-Type", "application/x-font-ttf")
		w.Write(index)
	}))
	http.Handle("/en_US.dic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving en_US.dic")
		indexBytes, err := content.Open("www/en_US.dic")
		if err != nil {
			panic(err)
		}
		index, err := ioutil.ReadAll(indexBytes)
		if err != nil {
			panic(err)
		}
		//w.Header().Add("Content-Type", "application/x-font-ttf")
		w.Write(index)
	}))
	http.Handle("/feeds", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feedsPage := createPageFromFeeds(db)
		w.Write(feedsPage)
	}))
	http.Handle("/addpeer", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if *peers != "" {
			if _, err := os.Stat(*peers); err == nil {
				log.Println("adding peer", r)
				//if r.Method == "POST" {

				peerbytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				peer := strings.Split(string(peerbytes), "peer=")[1]
				log.Println("peer", peer)
				if peer != "" {
					// append the peer to the end of the file at *peers
					peersFile, err := os.OpenFile(*peers, os.O_APPEND|os.O_WRONLY, 0600)
					if err != nil {
						panic(err)
					}
					defer peersFile.Close()
					peersFile.WriteString("\n" + peer)
					log.Println("added peer: " + peer)
				}
			}
		}
		w.Write([]byte("ok"))
	}))
	http.Handle("/post", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving post")
		// get the address of the request sender
		remoteAddr := r.RemoteAddr
		log.Println("Remote Addr:" + remoteAddr)

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		newthread := strings.Split(strings.Split(string(bodyBytes), "post=")[0], "thread=")[1]
		body := strings.Split(string(bodyBytes), "post=")[1]
		// hash the body
		if newthread == "true" {
			hash := sha256.Sum256([]byte(body))
			hashString := hex.EncodeToString(hash[:])
			err = os.MkdirAll(filepath.Join(db, hashString), 0755)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(filepath.Join(db, hashString, "op.txt"), WriteThread(hashString, body, remoteAddr), 0644)
			if err != nil {
				panic(err)
			}
			//w.Write(WriteThread(hashString, body))
			feedsPage := createPageFromFeeds(db)
			w.Write(feedsPage)
		} else {
			hash := sha256.Sum256([]byte(body))
			hashString := hex.EncodeToString(hash[:])
			err = ioutil.WriteFile(filepath.Join(db, newthread, hashString+".txt"), WritePost(newthread, hashString, body, remoteAddr), 0644)
			if err != nil {
				panic(err)
			}
			//w.Write(WritePost(newthread, hashString, body))
			feedsPage := createPageFromFeeds(db)
			w.Write(feedsPage)
		}
		AddAllGitDB(db)
	}))

	http.Serve(listener, nil)
}

func WriteThread(threadHash, body, remoteAddr string) []byte {
	unsafe := blackfriday.MarkdownCommon([]byte(body))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	div := `
<div class="thread">
	<div class="thread-hash">
		` + threadHash + `
	</div>
	<div class="thread-op">
		` + remoteAddr + `
	</div>
	<div class="post-body">
		` + string(html) + `
	</div>
</div>
`
	return []byte(div)
}

func WritePost(threadHash, postHash, body, remoteAddr string) []byte {
	unsafe := blackfriday.MarkdownCommon([]byte(body))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	div := `
<div class="thread-post">
	<div class="thread-hash">
		` + threadHash + `
	</div>
	<div class="post-op">
		` + remoteAddr + `
	</div>
	<div class="post-hash">
		` + postHash + `
	</div>
	<div class="post-body">
		` + string(html) + `
	</div>
</div>`
	return []byte(div)
}

var check = []string{"branches", "hooks", "info", "objects", "refs", "config", "description", "HEAD"}

func createPageFromFeeds(db string) []byte {
	// list all the directories in the db directory
	dirs, err := ioutil.ReadDir(db)
	if err != nil {
		panic(err)
	}
	var view []byte
	// ignore all git-related files and directories, like hooks, info, objects, refs, config, description, and HEAD
	for _, dir := range dirs {
		if dir.IsDir() && !strings.Contains(dir.Name(), ".git") {
			skip := false
			for _, checkDir := range check {
				if strings.Contains(dir.Name(), checkDir) {
					skip = true
				}
			}
			if !skip {
				// open the file named "op.txt"
				op, err := ioutil.ReadFile(filepath.Join(db, dir.Name(), "op.txt"))
				if err != nil {
					panic(err)
				}
				view = append(view, op...)
				threadDirs, err := ioutil.ReadDir(filepath.Join(db, dir.Name()))
				if err != nil {
					panic(err)
				}
				for _, threadDir := range threadDirs {
					if threadDir.IsDir() {
						// recurse, call this function again
						view = append(view, createPageFromFeeds(filepath.Join(db, dir.Name(), threadDir.Name()))...)
					}
					skip := false
					for _, checkDir := range check {
						if strings.Contains(dir.Name(), checkDir) {
							skip = true
						}
					}
					if skip {
						continue
					}
					// otherwise, we're looking at a file, so we can just append it
					if threadDir.Name() != "op.txt" {
						post, err := ioutil.ReadFile(filepath.Join(db, dir.Name(), threadDir.Name()))
						if err != nil {
							panic(err)
						}
						view = append(view, post...)
					}
				}
			}

		}
	}
	return view
}
