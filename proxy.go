package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os/exec"
)

type proxykeyring struct {
	backends []agent.Agent
}

// NewBackendAgent - Instantiate a single backend
func NewBackendAgent(backend string) (agent.Agent, error) {
	sock, err := net.Dial("unix", backend)
	if err != nil {
		return nil, err
	}
	backendAgent := agent.NewClient(sock)

	return backendAgent, nil
}

// NewProxyAgent - Aggregate multiple backends to a single agent
func NewProxyAgent(backends []agent.Agent) agent.Agent {
	return &proxykeyring{
		backends: backends,
	}
}

func (r *proxykeyring) List() ([]*agent.Key, error) {
	var keys []*agent.Key

	for _, backend := range r.backends {
		backendList, err := backend.List()
		if err != nil {
			return nil, err
		}

		for _, key := range backendList {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

func (r *proxykeyring) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	// TODO: Match key based on signing request

	cmd := exec.Command("zenity", "--question")
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	for _, backend := range r.backends {
		signResult, err := backend.Sign(key, data)
		if err != nil {
			continue
			// return nil, err
		}
		return signResult, nil
	}

	return nil, fmt.Errorf("No agent could successfully sign")
}

func (r *proxykeyring) Add(key agent.AddedKey) error {
	return fmt.Errorf("Add to each individual agent instead")
}

func (r *proxykeyring) Remove(key ssh.PublicKey) error {
	return fmt.Errorf("Remove to each individual agent instead")
}

func (r *proxykeyring) RemoveAll() error {
	return fmt.Errorf("Remove to each individual agent instead")
}

func (r *proxykeyring) Lock(passphrase []byte) error {
	return fmt.Errorf("Lock/unlock not supported")
}

func (r *proxykeyring) Unlock(passphrase []byte) error {
	return fmt.Errorf("Lock/unlock not supported")
}

func (r *proxykeyring) Signers() ([]ssh.Signer, error) {
	var signers []ssh.Signer

	for _, backend := range r.backends {
		bSigners, err := backend.Signers()
		if err != nil {
			return nil, err
		}

		for _, signer := range bSigners {
			signers = append(signers, signer)
		}
	}

	return signers, nil
}
