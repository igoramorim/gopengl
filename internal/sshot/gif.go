package sshot

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

var outdir = filepath.Join(".", "images")

func MakeGIF() error {
	files, err := os.ReadDir(tmpdir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	fmt.Printf("sshot: generating gif from %d images in %s\n", len(files), tmpdir)

	gifbuf := gif.GIF{LoopCount: len(files)}
	var outFilename string

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if outFilename == "" {
			outFilename = gifFilename(f.Name())
		}

		path := fmt.Sprintf("%s/%s", tmpdir, f.Name())
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Printf("sshot: reading: %s\n", path)

		src, err := gif.Decode(file)
		if err != nil {
			return err
		}

		// FIXME: The quality is terrible and animation not fluid
		dst := image.NewPaletted(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2), src.(*image.Paletted).Palette)
		draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

		gifbuf.Image = append(gifbuf.Image, dst)
		gifbuf.Delay = append(gifbuf.Delay, 15)
	}

	if outFilename == "" {
		return errors.New("sshot: could not build the outfile filename")
	}

	gifFile, err := os.Create(fmt.Sprintf("%s/%s.gif", outdir, outFilename))
	if err != nil {
		return err
	}
	defer gifFile.Close()

	err = gif.EncodeAll(gifFile, &gifbuf)
	if err != nil {
		return err
	}

	return deleteTmpDir()
}

func gifFilename(imgFilename string) string {
	idx := strings.LastIndex(imgFilename, ".")
	imgFilename = imgFilename[:idx]

	idx = strings.LastIndex(imgFilename, "_")
	imgFilename = imgFilename[:idx]

	return imgFilename
}

func deleteTmpDir() error {
	return os.RemoveAll(tmpdir)
}
