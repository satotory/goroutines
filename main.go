package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/nfnt/resize"
)

type Image struct {
	FileName string
	Index    int
	File     *os.File
	Out      *os.File
	Img      image.Image
	ResImg   image.Image
}

var wg sync.WaitGroup
var folder = "./image/"

func TLT() {
	// names := []Image{}
	fn := make(chan string)

	go ScanDir(fn)

	for v := range c {
		fmt.Printf("%+v\n", v)
		wg.Add(1)
		go ResImg(&v)
	}
	wg.Wait()
}

func ScanDir(c chan Image) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Panic(err)
	}
	for i, f := range files {
		img := Image{
			FileName: f.Name(),
			Index:    i,
		}
		fmt.Printf("%+d step complited\n", i)
		// go ResImg(f.Name())
		fn <- img.FileName
	}
	close(c)
}

var fCh = make(chan *os.File)
var iCh = make(chan image.Image)

func ResImg(obj *Image) {
	defer wg.Done()

	go func() {
		var err error
		obj.File, err = os.Open(folder + obj.FileName)
		if err != nil {
			log.Panic(err)
		}

		obj.Out, err = os.Create("./resized2/" + obj.FileName)
		if err != nil {
			log.Panic(err)
		}
		defer obj.Out.Close()
	}()

	go func() {
		var err error
		obj.Img, err = jpeg.Decode(<-msgCh)
		if err != nil {
			log.Panic(err)
		}

		obj.File.Close()
	}()

	go func() {
		if obj.Img != nil {
			obj.ResImg = resize.Thumbnail(300, 200, <-iCh, resize.Lanczos3)
		}
	}()

	select {
	case fn := obj.FileName:
		var err error
		obj.File, err = os.Open(folder + obj.FileName)
		if err != nil {
			log.Panic(err)
		}
		fCh <- obj.File
	case file := <-fCh:
		var err error
		obj.Img, err = jpeg.Decode(file)
		if err != nil {
			log.Panic(err)
		}

		obj.File.Close()
	}

	jpeg.Encode(out, resImg, nil)
}
