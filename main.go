package main

import (
	"flag"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nfnt/resize"

	_ "image/png"
)

const Delay = time.Second * 1

var Src string
var Dst string
var W uint
var H uint
var Q int

func init() {
	flag.UintVar(&W, "w", 1024, "max thumbnail width")
	flag.UintVar(&H, "h", 1024, "max thumbnail height")
	flag.IntVar(&Q, "q", 95, "jpeg quality")
	flag.StringVar(&Src, "src", ".", "directory to watch for images")
	flag.StringVar(&Dst, "dst", "thumbs", "directory to place thumbnails")
}

func loadImage(name string) (image.Image, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	return im, err
}

func saveImage(name string, im image.Image) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, im, &jpeg.Options{Q})
}

func createThumbnail(src, dst string) error {
	im, err := loadImage(src)
	if err != nil {
		return err
	}
	im = resize.Thumbnail(W, H, im, resize.Lanczos3)
	return saveImage(dst, im)
}

func main() {
	flag.Parse()
	os.MkdirAll(Dst, 0755)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(Src)
	if err != nil {
		log.Fatal(err)
	}

	timers := make(map[string]*time.Timer)
	ch := make(chan string)

	for {
		select {
		case name := <-ch:
			delete(timers, name)
			src := path.Join(Src, name)
			dst := path.Join(Dst, name+".jpg")
			log.Printf("%s -> %s\n", src, dst)
			err := createThumbnail(src, dst)
			if err != nil {
				log.Println(err)
			}
		case event := <-watcher.Events:
			timer, ok := timers[event.Name]
			if ok {
				timer.Reset(Delay)
			} else {
				timers[event.Name] = time.AfterFunc(Delay, func() {
					ch <- event.Name
				})
			}
		case err := <-watcher.Errors:
			log.Println(err)
		}
	}
}
