package sshot

import (
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func NewScreenShoter(filename string, width, height int) *ScreenShoter {
	return &ScreenShoter{
		filename: filename,
		width:    width,
		height:   height,
		max:      24,
	}
}

type ScreenShoter struct {
	filename  string
	width     int
	height    int
	data      [][]uint8
	max       int
	deltaTime float64
	lastFrame float64
}

func (ss *ScreenShoter) TakeOne() {
	pixels := make([]uint8, 4*ss.width*ss.height) // 4 = R G B A
	gl.ReadPixels(0, 0, int32(ss.width), int32(ss.height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))

	f, err := os.Create(fmt.Sprintf("%s/%s%s", filepath.Join(".", "images"), ss.filename, ".png"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, ss.width, ss.height))
	img.Pix = pixels

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func (ss *ScreenShoter) Take() {
	if len(ss.data) >= ss.max {
		return
	}

	currentFrame := glfw.GetTime()
	ss.deltaTime = currentFrame - ss.lastFrame
	// fmt.Printf("sshot: current frame: %.4f delta frame: %.4f last frame: %.4f\n", currentFrame, ss.deltaTime, ss.lastFrame)

	// Kinda sketchy
	if ss.deltaTime >= 0.3 {
		ss.lastFrame = currentFrame

		fmt.Printf("sshot: taking shot\n")

		pixels := make([]uint8, 4*ss.width*ss.height) // 4 = R G B A
		gl.ReadPixels(0, 0, int32(ss.width), int32(ss.height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))
		ss.data = append(ss.data, pixels)
	}
}

func (ss *ScreenShoter) Save() {
	fmt.Printf("sshot: %d total tmp pixels buffers\n", len(ss.data))
	for i := range ss.data {
		ss.save(i)
	}
}

const tpmImgExtension = ".gif"

var tmpdir = filepath.Join(".", "images", "tmp")

// FIXME: The iamges are being saved upside down
func (ss *ScreenShoter) save(idx int) {
	fmt.Printf("sshot: saving tmp img idx %d\n", idx)

	err := os.MkdirAll(tmpdir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fmt.Sprintf("%s/%s_%06d%s", tmpdir, ss.filename, idx, tpmImgExtension))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, ss.width, ss.height))
	img.Pix = ss.data[idx]

	if err := gif.Encode(f, img, nil); err != nil {
		panic(err)
	}
}
