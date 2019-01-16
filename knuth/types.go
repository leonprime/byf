package main

import (
	"fmt"
	"unicode"
)

type Node struct {
	L, R, U, D *Node
	C          *Header
}

type Header struct {
	Node
	S int
	N string
}
