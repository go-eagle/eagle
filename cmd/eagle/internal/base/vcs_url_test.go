package base

import (
	"net"
	"strings"
	"testing"
)

func TestParseVCSUrl(t *testing.T) {
	repos := []string{
		// ssh://[user@]host.xz[:port]/path/to/repo.git/
		"ssh://git@github.com:7875/go-eagle/eagle.git",
		// git://host.xz[:port]/path/to/repo.git/
		"git://github.com:7875/go-eagle/eagle.git",
		// http[s]://host.xz[:port]/path/to/repo.git/
		"https://github.com:7875/go-eagle/eagle.git",
		// ftp[s]://host.xz[:port]/path/to/repo.git/
		"ftps://github.com:7875/go-eagle/eagle.git",
		//[user@]host.xz:path/to/repo.git/
		"git@github.com:go-eagle/eagle.git",
		// ssh://[user@]host.xz[:port]/~[user]/path/to/repo.git/
		"ssh://git@github.com:7875/go-eagle/eagle.git",
		// git://host.xz[:port]/~[user]/path/to/repo.git/
		"git://github.com:7875/go-eagle/eagle.git",
		//[user@]host.xz:/~[user]/path/to/repo.git/
		"git@github.com:go-eagle/eagle.git",
		///path/to/repo.git/
		"~/go-eagle/eagle.git",
		// file:///path/to/repo.git/
		"file://~/go-eagle/eagle.git",
	}
	for _, repo := range repos {
		url, err := ParseVCSUrl(repo)
		if err != nil {
			t.Fatal(repo, err)
		}
		urlPath := strings.TrimLeft(url.Path, "/")
		if urlPath != "go-eagle/eagle.git" {
			t.Fatal(repo, "parse url failed", urlPath)
		}
	}
}

func TestParseSsh(t *testing.T) {
	repo := "ssh://git@github.com:7875/go-eagle/eagle.git"
	url, err := ParseVCSUrl(repo)
	if err != nil {
		t.Fatal(err)
	}
	host, _, err := net.SplitHostPort(url.Host)
	if err != nil {
		host = url.Host
	}
	t.Log(host, url.Path)
}
