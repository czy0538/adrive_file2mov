package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	HEAD = `6D6F6F76`                 // MOV file header.
	END  = `0000000049454E44AE426082` // end of png file, just too lazy to delete.
)

func main() {
	head, _ := hex.DecodeString(HEAD)
	end, _ := hex.DecodeString(END)
	fmt.Println("Welcome to use this small tool, which can help you share any types of file by adrive")

	// check filename
	inputFilename := flag.String("i", "", "the path of file you want to convert")
	outputFilename := flag.String("o", "", "output file path")
	flag.Parse()
	if *inputFilename == "" || !fileExists(*inputFilename) {
		log.Fatalln("Please input correct filename you want to convert.")
	}
	fmt.Println(path.Ext(*inputFilename))
	fmt.Println(path.Base(*inputFilename))
	fmt.Println(fileprefix(*inputFilename))
	if path.Ext(*inputFilename) == ".mov" { //decode
		if *outputFilename == "" {
			*outputFilename = fileprefix(*inputFilename)
		}
		if fileExists(*outputFilename) {
			log.Fatalf("file %s exites.\n", *outputFilename)
		}

		inf, err := os.Open(*inputFilename)
		if err != nil {
			log.Fatalln(err)
		}
		defer inf.Close()
		outf, err := os.Create(*outputFilename)
		if err != nil {
			log.Fatalln(err)
		}
		defer outf.Close()

		_, err = inf.Seek(int64(len(head)), 0)
		if err != nil {
			log.Fatalln(err)
		}

		err = readBigFile(inf, func(data []byte) error {
			return fileAppend(outf, data)
		})
		if err != nil {
			log.Fatalln(err)
		}

		size, _ := outf.Seek(0, 1)
		size -= int64(len(end))
		err = outf.Truncate(size)
		if err != nil {
			log.Fatalln(err)
		}
	} else { //encode
		if *outputFilename == "" {
			*outputFilename = fileprefix(*inputFilename) + ".mov"
		}
		if fileExists(*outputFilename) {
			log.Fatalf("file %s exites.\n", *outputFilename)
		}

		inf, err := os.Open(*inputFilename)
		if err != nil {
			log.Fatalln(err)
		}
		defer inf.Close()
		outf, err := os.Create(*outputFilename)
		if err != nil {
			log.Fatalln(err)
		}
		defer outf.Close()

		if fileAppend(outf, head) != nil {
			log.Fatalln(err)
		}

		if readBigFile(inf, func(data []byte) error {
			return fileAppend(outf, data)
		}) != nil {
			log.Fatalln(err)
		}

		if fileAppend(outf, end) != nil {
			log.Fatalln(err)
		}
	}
	fmt.Printf("convert file %s success\n", *outputFilename)
}

func readBigFile(f *os.File, handle func([]byte) error) error {
	s := make([]byte, 4096)
	for {
		switch nr, err := f.Read(s[:]); true {
		case nr < 0:
			return fmt.Errorf("ReadBigFile error:%v", err)
		case nr == 0: // EOF
			return nil
		case nr > 0:
			err = handle(s[0:nr])
			if err != nil {
				return fmt.Errorf("ReadBigFile error:%v", err)
			}
		}
	}
}

func fileAppend(f *os.File, data []byte) error {
	if f == nil {
		return fmt.Errorf("os.File error")
	}
	_, err := f.Seek(0, 2)
	if err != nil {
		return fmt.Errorf("fileAppend seek error:%v", err)
	}
	n, err := f.Write(data)
	if err != nil {
		return fmt.Errorf("fileAppend write error:%v", err)
	}
	if n != len(data) {
		return fmt.Errorf("only write %d bytes into file, %d bytes in total", n, len(data))
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func fileprefix(filename string) string {
	filesuffix := path.Ext(filename)
	fileprefix := filename[0 : len(filename)-len(filesuffix)]
	return fileprefix
}
