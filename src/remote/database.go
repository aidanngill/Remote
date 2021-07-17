package main

type Package struct {
	ID       string
	Path     string
	Row      uint16
	Size     uint64
	Progress uint64
	Task     Task
}

var (
	packages []Package
)
