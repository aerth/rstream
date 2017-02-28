// rstream streams to stdout
package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

var chunk = "1440"
var host string
var defaultsocks = "socks5://127.0.0.1:1080"
var defaultradio = "anonradio.net"
var defaultmount = "anonradio"
var chunksize = 1440
var custom = "FALSE"

func init() {
	var err error
	chunksize, err = strconv.Atoi(chunk)
	if err != nil {
		fmt.Println("Invalid chunk size")
		os.Exit(1)
	}
}

func openStream(socksurl, host, mountpoint string) {
	mountpoint = strings.TrimPrefix(mountpoint, "/")
	p := func() proxy.Dialer {
		if socksurl == "none" {
			return proxy.FromEnvironment()
		}
		if socksurl == "tor" {
			socksurl = "socks5://127.0.0.1:9050"
		}
		u, err := url.Parse(socksurl)
		if err != nil {
			fmt.Println("SOCKS", err)
			os.Exit(2)
		}
		px, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			fmt.Println("SOCKS", err)
			os.Exit(2)
		}
		return px
	}()

	conn, err := p.Dial("tcp", host)
	for err != nil {
		fmt.Fprintln(os.Stderr, err)
		time.Sleep(5 * time.Second)
		conn, err = p.Dial("tcp", host)
	}
	defer conn.Close()
	fmt.Fprintln(os.Stderr, "connected")
	fmt.Fprint(conn, "GET /"+mountpoint+" HTTP/1.0\r\n\r\n")

	for {
		buff := make([]byte, chunksize) // 128kb/s ?
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Fprint(os.Stderr, n, err, " ")
		} else {
			if os.Getenv("DEBUG") != "" {fmt.Fprint(os.Stderr, n, "bytes ")}
		}
		os.Stdout.Write(buff)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "remote mp3 to stdout")
	fmt.Fprintln(os.Stderr, "proxy defaults to socks5://127.0.0.1:1080")
	fmt.Fprintln(os.Stderr, "Usage: ", os.Args[0], " <host.server.com:port> <mountpoint> [<SOCKS>]")
	fmt.Fprintln(os.Stderr, "Example (no proxy)", os.Args[0], "example.com:8000 techno none")
}

func init() {
	if custom != "TRUE" { // CUSTOM="TRUE" in Makefile
		if len(os.Args) != 3 && len(os.Args) != 4 {
			printUsage()
			os.Exit(1)
		}
	}

}
func main() {
	var host, mountpoint, socks string
	if len(os.Args) == 2 {
		host = os.Args[1]
	}

	if len(os.Args) == 3 {
		host, mountpoint = os.Args[1], os.Args[2]
	}

	if len(os.Args) == 4 {
		socks = os.Args[3]
	}
	if host == "" {
		host = defaultradio
	}
	if mountpoint == "" {
		mountpoint = defaultmount
	}
	if socks == "" {
		socks = defaultsocks
	}

	openStream(socks, host, mountpoint)
}
