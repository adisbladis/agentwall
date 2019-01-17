package main // import "github.com/adisbladis/agentwall"

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
	"os"
)

func main() {
	var backendPaths ArrayFlags
	flag.Var(&backendPaths, "backend", "UNIX socket path to agent backend")
	flag.Parse()

	if len(backendPaths) == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Aggregate multiple agents into a single interface
	var backends []agent.Agent
	for _, sockPath := range backendPaths {
		backend, err := NewBackendAgent(sockPath)
		if err != nil {
			panic(err)
		}
		backends = append(backends, backend)
	}
	proxyAgent := NewProxyAgent(backends)

	ln, err := net.Listen("unix", "./go.sock")
	if err != nil {
		panic(err)
	}

	for {
		fd, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			err = agent.ServeAgent(proxyAgent, fd)
			if err != nil && err != io.EOF {
				panic(err)
			}
		}()
	}

	fmt.Println(proxyAgent)
}
