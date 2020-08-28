package main

import (
	"bufio"
	"image"
	"image/gif"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func openAndDecode(path string) ([]*image.Paletted, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	g, err := gif.DecodeAll(r)
	if err != nil {
		return nil, err
	}

	return g.Image, nil
}

func main() {

	app := &cli.App{
		Name:  "concat gifs",
		Usage: "you can concat gifs",
		Action: func(c *cli.Context) error {

			var allFrames []*image.Paletted = make([]*image.Paletted, 0)

			readErr := filepath.Walk("gifs", func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}

				log.Println("concatinating", path)

				frames, err := openAndDecode(path)
				if err != nil {
					return err
				}

				allFrames = append(allFrames, frames...)

				return nil
			})

			if readErr != nil {
				return readErr
			}

			delays := make([]int, len(allFrames))
			for i := range allFrames {
				delays[i] = 0
			}

			f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
			defer f.Close()
			return gif.EncodeAll(f, &gif.GIF{
				Image: allFrames,
				Delay: delays,
			})
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
