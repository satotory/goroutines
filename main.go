package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
)

type Iob struct {
	Path       string
	ImageName  string
	osFile     *os.File
	decoded    image.Image
	resIm      image.Image
	pathToSave string
	saveTo     *os.File
}

var opfch, decch, encch, resch = make(chan Iob), make(chan Iob), make(chan Iob), make(chan Iob)
var quit, sizze, closed = make(chan int, 1), make(chan int, 1), make(chan int, 1)

func main() {
	fmt.Println("Start")
	f := "./image/"
	go func() {
		for o := range opfch {
			go ope(o)
		}
	}()
	go func() {
		for o := range encch {
			go resIm(o)
			select {
			case op := <-resch:
				fmt.Printf("LAST STAGE WITH %s\n", op.pathToSave)
			}
		}
		quit <- 0
	}()
	sel()
	getFiles(f)

	fmt.Println("End")
}

func sel() {
	for i := 0; i < <-sizze; i++ {
		select {
		case opF := <-decch:
			go decIm(opF)
		case p := <-encch:
			go resIm(p)
		case op := <-resch:
			fmt.Printf("LAST STAGE WITH %s\n", op.pathToSave)
		}
	}
}

func decIm(opF Iob) {
	fmt.Printf("FILE opened - %s\n", opF.ImageName)
	var err error
	opF.decoded, err = jpeg.Decode(opF.osFile)
	if err != nil {
		log.Print(err)
	}
	encch <- opF
}

func resIm(o Iob) {
	fmt.Printf("resizing - %s\n", o.ImageName)
	resch <- o
}

func encodeJp(o Iob) {
	fmt.Printf("%s in encodeJp\n", o.ImageName)
	// var err error
	// o.saveTo, err = os.Create(o.pathToSave)
	// if err != nil {
	// 	log.Print(err)
	// }
	// o.saveTo.Close()

	// o.resIm, err = resize.Thumbnail(300, 200, o.decoded, resize.Lanczos3)
}

func ope(o Iob) {
	var err error
	o.osFile, err = os.Open(o.Path)
	if err != nil {
		log.Print(err)
	}
	decch <- o
}

func getFiles(fold string) {
	files, err := ioutil.ReadDir(fold)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("SIZE - %d\n", len(files))
	sizze <- len(files)
	for i, f := range files {
		select {
		case opfch <- Iob{Path: fold + f.Name(), ImageName: f.Name(), pathToSave: "./resized/" + f.Name()}:
			fmt.Printf("+INDEX - %d %s\n", i, fold+f.Name())
		case <-quit:
			fmt.Print("quit\n")
			return
		}
	}

}
