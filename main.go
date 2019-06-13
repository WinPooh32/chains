package main

import (
	"net"
	"crypto/md5"
	"fmt"
	"log"

	"github.com/WinPooh32/awesome-merkle-tree/merkle"
)

func md5Wrap(data []byte) string {
	sum := md5.Sum(data)
	return fmt.Sprintf("%X", sum)
}

func main() {
	t := merkle.MakeTree(md5Wrap)

	data := [][]byte{
		[]byte{6, 6, 6},
		[]byte{3, 2, 1},
		[]byte{1, 2},
		[]byte{4, 4},
		[]byte{5, 5},
		[]byte{7, 7},
	}

	if err := t.Insert(data); err != nil {
		log.Fatalln(err)
	}

	log.Println("Tree: \n", t)
}
