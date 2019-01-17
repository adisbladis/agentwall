package main

import (
	"strings"
)

// ArrayFlags - Accept multiple of a CLI flag
type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return "HACK: string repr"
}

// Set - Set value
func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}
