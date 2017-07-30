package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/jmcvetta/randutil"
)

func copyloop(dir string, wait time.Duration, name func(string) string) {
	imgs, err := AssetDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case <-time.Tick(wait):
			img, err := randutil.ChoiceString(imgs)
			if err != nil {
				fmt.Println(err)
				continue
			}
			imgpath := dir + "/" + img
			data, err := Asset(imgpath)
			if err != nil {
				// Asset was not found.
				fmt.Println(err)
				continue
			}

			go copy(bytes.NewReader(data), name(imgpath))
		}
	}
}

func copy(src *bytes.Reader, dn string) error {
	dst, err := os.Create(dn)
	defer dst.Close()
	if err != nil {
		switch e := err.(type) {
		case *os.PathError:
			d, _ := path.Split(e.Path)
			err := os.MkdirAll(d, 0777)
			if err != nil {
				return err
			}
			dst, err = os.Create(dn)
		default:
			log.Fatal(err)
			if dst != nil {
				dst.Close()
			}
			return err
		}
	}
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
	return nil
}
