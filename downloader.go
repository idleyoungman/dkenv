package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
)

var (
	errResource          = errors.New("invalid webfinger resource")
	errNotYetImplemented = errors.New("Not yet implemented")
	errTooManyRedirects  = errors.New("Too many redirects")
	errHTTPRedirect      = errors.New("Redirect to non-https server")
	errHTTPCode          = errors.New("Received unexpected http code")
	errSubjectMissmatch  = errors.New("Subject doesn't match resource")
)

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.total += int64(n)
		percentage := float64(pt.total) / float64(pt.length) * float64(100)
		// i := int(percentage / float64(10))
		// is := fmt.Sprintf("Transferred %v bytes\n", i)
		is := fmt.Sprintf("Transferred %d percent\n", int(percentage))
		if percentage-pt.progress > 2 {
			fmt.Fprintf(os.Stderr, is)
			pt.progress = percentage
		}
	}

	return n, err
}

// PassThru wraps an existing io.Reader.
//
// It simply forwards the Read() call, while displaying
// the results from individual calls to it.
type PassThru struct {
	io.Reader
	total    int64 // Total # of bytes transferred
	length   int64 // Expected length
	progress float64
}

func createDotDKEnvDirectory() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	_err := os.Mkdir(usr.HomeDir+"/.dkenv/", 0700)
	if _err != nil {
	}
}

func createVersionFile(version string) (out *os.File) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	out, _err := os.Create(usr.HomeDir + "/.dkenv/docker-" + version)

	if _err != nil {
		log.Fatal(_err)
	}
	return out
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	if len(via) > 10 {
		return errTooManyRedirects
	}

	if req.URL.Scheme != "https" {
		return errHTTPRedirect
	}

	return nil
}

func getHTTP(version string) *http.Response {
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	system := "Darwin"

	if runtime.GOOS == "windows" {
		system = "Windows"
	}

	if runtime.GOOS == "linux" {
		system = "Linux"
	}

	if runtime.GOOS == "darwin" {
		system = "Darwin"
	}

	resp, _err := client.Get("https://get.docker.com/builds/" + system + "/x86_64/docker-" + version)

	if _err != nil {
		log.Fatal(_err)
	}

	return resp

}

func getDocker(version string, binDir string) {
	createDotDKEnvDirectory()
	// Create the docker binary file
	// out := createVersionFile(version)

	// defer out.Close()
	// Do the htp get
	resp := getHTTP(version)

	defer resp.Body.Close()

	readerpt := &PassThru{Reader: resp.Body, length: resp.ContentLength}

	body, _err := ioutil.ReadAll(readerpt)

	if _err != nil {
		log.Fatal(_err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if http.DetectContentType(body) != "application/octet-stream" {
		log.Fatal("download not Content-Type application/octet-stream; check the version exists")
	}

	_err = ioutil.WriteFile(usr.HomeDir+"/.dkenv/docker-"+version, body, 0777)

	switchVersion(version, binDir)

	if _err != nil {
		log.Fatal(_err)
	}

}
