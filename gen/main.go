package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/c4milo/unpackit"
	cp "github.com/otiai10/copy"
)

func DownloadAndSave(url string, filename string) (string, error) {
	// Download the file using HTTP from the given URL
	// Save it to the given filename
	// Return an error if something went wrong
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func UnzipTarGzFile(filename, tempDir string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	dir, err := unpackit.Unpack(file, ".")
	if err != nil {
		return "", err
	}
	// move dir to tempDir
	err = os.Rename(dir, tempDir)
	if err != nil {
		return "", err
	}
	return tempDir, nil
}

func main() {
	// download and extract: https://github.com/Ionaru/easy-markdown-editor/tarball/master
	filename, err := DownloadAndSave("https://github.com/Ionaru/easy-markdown-editor/tarball/master", "easymde.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll("easymde")
	if err != nil {
		log.Fatal(err)
	}
	dirname, err := UnzipTarGzFile(filename, "easymde")
	if err != nil {
		log.Fatal(err)
	}
	cmdNpmCi := exec.Command("npm", "ci")
	cmdNpmCi.Dir = dirname
	cmdNpmCi.Stdout = os.Stdout
	cmdNpmCi.Stderr = os.Stderr
	err = cmdNpmCi.Run()
	if err != nil {
		log.Fatal(err)
	}
	cmdGulp := exec.Command("gulp")
	cmdGulp.Dir = dirname
	cmdGulp.Stdout = os.Stdout
	cmdGulp.Stderr = os.Stderr
	err = cmdGulp.Run()
	if err != nil {
		log.Fatal(err)
	}
	// remove www/easymde/
	err = os.RemoveAll("www/easymde/")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Unzipped to:", dirname)
	// copy easymde/dist/ to www/easymde/
	err = cp.Copy("easymde/dist", "www/easymde")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Copied to:", "www/easymde")
	// https://maxcdn.bootstrapcdn.com/font-awesome/latest/css/font-awesome.min.css
	facss, err := DownloadAndSave("https://maxcdn.bootstrapcdn.com/font-awesome/latest/css/font-awesome.min.css", "www/font-awesome.min.css")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Downloaded:", facss)
	// https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff2?v=4.7.0
	fawoff2, err := DownloadAndSave("https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff2?v=4.7.0", "www/fontawesome-webfont.woff2")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Downloaded:", fawoff2)
	// replace https://maxcdn.bootstrapcdn.com/font-awesome/latest/css/font-awesome.min.css in www/easymde/easymde.min.js with www/font-awesome.min.css
	err = ReplaceInFile("www/easymde/easymde.min.js", "https://maxcdn.bootstrapcdn.com/font-awesome/latest/css/font-awesome.min.css", "font-awesome.min.css")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Replaced:", "www/easymde/easymde.min.js")
	// replace https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff2?v=4.7.0 in www/easymde/easymde.min.js with www/fontawesome-webfont.woff2
	err = ReplaceInFile("www/easymde/easymde.min.js", "https://maxcdn.bootstrapcdn.com/font-awesome/latest/fonts/fontawesome-webfont.woff2?v=4.7.0", "fontawesome-webfont.woff2")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Replaced:", "www/easymde/easymde.min.js")
	// replace maxcdn.bootstrapcdn.com/font-awesome with nothing
	err = ReplaceInFile("www/easymde/easymde.min.js", "maxcdn.bootstrapcdn.com/font-awesome", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Replaced:", "www/easymde/easymde.min.js")
	// https://cdn.jsdelivr.net/codemirror.spell-checker/latest/en_US.aff
	aff, err := DownloadAndSave("https://cdn.jsdelivr.net/codemirror.spell-checker/latest/en_US.aff", "www/en_US.aff")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Downloaded:", aff)
	// https://cdn.jsdelivr.net/codemirror.spell-checker/latest/en_US.dic
	dic, err := DownloadAndSave("https://cdn.jsdelivr.net/codemirror.spell-checker/latest/en_US.dic", "www/en_US.dic")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Downloaded:", dic)
	// replace https://cdn.jsdelivr.net/codemirror.spell-checker/latest/en_US.aff in www/easymde/easymde.min.js with en_US.aff
	err = ReplaceInFile("www/easymde/easymde.min.js", "https://cdn.jsdelivr.net/codemirror.spell-checker/latest/en_US.aff", "en_US.aff")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Replaced:", "www/easymde/easymde.min.js")
	// replace https://cdn.jsdelivr.net/codemirror.spell-checker/latest/en_US.dic in www/easymde/easymde.min.js with en_US.dic
	err = ReplaceInFile("www/easymde/easymde.min.js", "https://cdn.jsdelivr.net/codemirror.spell-checker/latest/en_US.dic", "en_US.dic")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Replaced:", "www/easymde/easymde.min.js")

}

func ReplaceInFile(filename, search, replace string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// Read the file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	// Replace the target string
	replaced := bytes.Replace(data, []byte(search), []byte(replace), -1)
	// Write the file
	err = ioutil.WriteFile(filename, replaced, 0644)
	if err != nil {
		return err
	}
	return nil
}
