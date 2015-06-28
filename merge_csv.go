package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

func merge_and_generate(f1, f2, out *os.File) {
	r1 := bufio.NewReader(f1)
	r2 := bufio.NewReader(f2)
	w := csv.NewWriter(out)
	defer func() {
		w.Flush()
	}()
	for {
		buf1, partp1, err1 := r1.ReadLine()
		buf2, partp2, err2 := r2.ReadLine()
		if err1 == err2 && err1 == io.EOF {
			break
		} else {
			if err1 != nil {
				log.Fatal(err1)
			}
			if err2 != nil {
				log.Fatal(err2)
			}
		}
		if partp1 || partp2 {
			log.Fatal("line too long")
		}
		if len(buf1) != len(buf2) {
			log.Fatal("has different line length")
		}
		base := 0
		row := make([]string, 0)
		for i := 0; i < len(buf1); i++ {
			if buf1[i] != buf2[i] {
				row = append(row, string(buf1[base:i]))
				base = i + 1
			}
		}
		row = append(row, string(buf1[base:]))
		w.Write(row)
	}
}

func main() {
	fn1 := os.Args[1]
	f1, err := os.Open(fn1)
	if err != nil {
		log.Fatal(err)
	}
	fn2 := os.Args[2]
	f2, err := os.Open(fn2)
	if err != nil {
		log.Fatal(err)
	}
	merge_and_generate(f1, f2, os.Stdout)
}
