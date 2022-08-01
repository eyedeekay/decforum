package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eyedeekay/goSam"
	sam "github.com/eyedeekay/sam3/helper"
	git "github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sosedoff/gitkit"

	"github.com/go-git/go-git/v5/plumbing/transport/client"
)

// EnsureGitDBInitialized ensures that the git database is initialized, accepts directory(repository) as argument
func EnsureGitDBInitialized(dir string) {
	// check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// if it doesn't exist, create it
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
	// check if the directory is a git repository
	_, err := os.Stat(dir + "/.git")
	if os.IsNotExist(err) {
		// initialize the git repository
		_, err := git.PlainInit(dir, false)
		if err != nil {
			panic(err)
		}
	}
}

func AddAllGitDB(dir string) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		panic(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		panic(err)
	}
	// add every file in the dir to the git repository, recursively
	e := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		skip := false
		for _, checkDir := range check {
			if strings.Contains(path, checkDir) {
				skip = true
			}
		}
		if !skip {
			if !f.IsDir() && !strings.Contains(path, "/.git") {
				// add path to the repository
				path := strings.Replace(path, dir+"/", "", 1)
				log.Println("adding", path)
				what, err := w.Add(path)
				if err != nil {
					panic(err)
				}
				log.Println(what)
				CommitAllToGitDB(dir)
			}
		}
		return err

	})
	if e != nil {
		panic(e)
	}
}
func CommitAllToGitDB(dir string) {
	// msg is the unix timestamp of the commit
	msg := fmt.Sprintf("%d", time.Now().Unix())
	repo, err := git.PlainOpen(dir)
	if err != nil {
		panic(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		panic(err)
	}
	date := time.Unix(0, 0)
	log.Println("committing", msg, "as if it occurred on", date, "for deterministic purposes")
	// commit the changes to the git repository
	hash, err := w.Commit(msg, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "idk",
			Email: "instance@listener.i2p",
			When:  date,
		},
		Committer: nil,
	})
	if err != nil {
		panic(err)
	}
	log.Println("committed to git", hash)
}

func HostARemote(dir string) error {
	// Configure git service
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	log.Println("configuring git service", dir)
	service := gitkit.New(gitkit.Config{
		Dir:        dir + "/",
		AutoCreate: true,
	})

	// Configure git server. Will create git repos path if it does not exist.
	// If hooks are set, it will also update all repos with new version of hook scripts.
	if err := service.Setup(); err != nil {
		log.Fatal(err)
	}

	httpmux := &http.ServeMux{}

	httpmux.Handle("/", service)

	httpServer := &http.Server{
		Handler: httpmux,
	}

	ln, err := sam.I2PListener("gitforum-db-"+randId(), "127.0.0.1:7656", "gitforum-db")
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server
	if err := httpServer.Serve(ln); err != nil {
		log.Fatal(err)
	}
	return nil
}

func CloneARemote(remote string, dir string) error {
	sam, err := goSam.NewDefaultClient()
	// Create a custom http(s) client with your config
	customClient := &http.Client{
		// accept any certificate (might be useful for testing)
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial:            sam.Dial,
		},

		// 15 second timeout
		Timeout: 15 * time.Second,

		// don't follow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	defer sam.Close()

	// Override http(s) default protocol to use our custom client
	client.InstallProtocol("http", githttp.NewClient(customClient))

	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: remote,
	})
	if err != nil {
		return err
	}
	head, err := r.Head()
	if err != nil {
		return err
	}
	log.Println("cloned", remote, "to", dir, "at", head.Hash())
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	// checkout master
	err = w.Checkout(&git.CheckoutOptions{
		Hash: head.Hash(),
	})
	return nil
}

func PullInCommitsFromRemote(remote, dir string) error {
	sam, err := goSam.NewDefaultClient()
	// Create a custom http(s) client with your config
	customClient := &http.Client{
		// accept any certificate (might be useful for testing)
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial:            sam.Dial,
		},

		// 15 second timeout
		Timeout: 15 * time.Second,

		// don't follow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	defer sam.Close()

	// Override http(s) default protocol to use our custom client
	client.InstallProtocol("http", githttp.NewClient(customClient))

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}
	// get the hostname from the remote, which will be an HTTP URL
	url, err := url.Parse(remote)
	if err != nil {
		return err
	}
	hostname := url.Hostname()
	repo.CreateRemote(&config.RemoteConfig{
		Name: hostname,
		URLs: []string{remote},
	})

	// if the remote isn't added yet, add it, refer to it by it's hostname
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	// pull from the remote
	err = w.Pull(&git.PullOptions{RemoteName: hostname, Force: true})
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: head.Hash(),
	})
	return nil
}
