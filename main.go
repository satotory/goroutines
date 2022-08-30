package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
)

const (
	FOLDER    = "./images/"
	RESFOLDER = "./resized_images/"
)

const (
	ERR_READDIR       = "Read dir error"
	ERR_OPENINGFILE   = "os Opening file error"
	ERR_DECODINGFILE  = "png Decoding file error"
	ERR_CREATING_FILE = "os Create file error"
	ERR_ENCODING_PIC  = "png Encode pic error"
)

// chanels
var (
	firstCh   = make(chan imgs)
	secondCh  = make(chan img)
	thirdCh   = make(chan img)
	fourthCh  = make(chan img)
	endCh     = make(chan img)
	endMainCh = make(chan string)
)

type img struct {
	file      fs.DirEntry
	index     int64
	openedImg *os.File
	saveFile  *os.File
	decoImg   image.Image
	resImg    image.Image
	finished  bool
}

type imgs struct {
	Files []fs.DirEntry
	count int64
}

// files handler
func selectCases() {
	var (
		lenFs    = int64(-1)
		counterF int64
	)

	for {
		if lenFs == counterF {
			endMainCh <- "end of main goroutine"
			break
		}

		select {
		case imgs := <-firstCh:
			lenFs = imgs.count
			for i, file := range imgs.Files {
				img := img{file: file, index: int64(i)}
				go img.openingGorou()
			}
		case f := <-secondCh:
			go f.decodingImgGorou()
		case f := <-thirdCh:
			go f.createFileToSave()
		case f := <-fourthCh:
			go f.resizeAndEncodeImg()
		case f := <-endCh:
			f.finished = true
			counterF++
			fmt.Println(f.file.Name(), f.index, counterF)
		}
	}
}

func (f img) openingGorou() {
	var err error

	f.openedImg, err = os.Open(FOLDER + f.file.Name())
	if err != nil {
		log.Println(ERR_OPENINGFILE)
		log.Fatal(err)
	}

	secondCh <- f
}

func (f img) decodingImgGorou() {
	var err error

	f.decoImg, err = jpeg.Decode(f.openedImg)
	if err != nil {
		log.Println(ERR_DECODINGFILE)
		log.Fatal(err)
	}
	f.openedImg.Close()

	thirdCh <- f
}

func (f img) createFileToSave() {
	var err error

	f.saveFile, err = os.Create(RESFOLDER + f.file.Name())
	if err != nil {
		log.Println(ERR_CREATING_FILE)
		log.Fatal(err)
	}

	fourthCh <- f
}

func (f img) resizeAndEncodeImg() {
	defer f.saveFile.Close()
	f.resImg = resize.Thumbnail(300, 200, f.decoImg, resize.Lanczos3)

	err := jpeg.Encode(f.saveFile, f.resImg, nil)
	if err != nil {
		log.Println()
		log.Fatal(err)
	}
	endCh <- f
}

func run() {
	start := time.Now()
	go selectCases()

	fs, err := os.ReadDir(FOLDER)
	if err != nil {
		log.Println(ERR_READDIR)
		log.Fatal(err)
	}

	firstCh <- imgs{Files: fs, count: int64(len(fs))}

	elapsed := time.Since(start)
	log.Printf("\nExecution time %s\n%s", elapsed, <-endMainCh)
}

func main() {
	run()
}
