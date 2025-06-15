package texture

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func New(imgPath string, texType int, slotType uint32, sourceFormat, destFormat, pixelType int) (*Texture, error) {
	var id uint32

	gl.GenTextures(1, &id)
	gl.BindTexture(uint32(texType), id)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	imageData, err := loadImage(imgPath)
	if err != nil {
		return nil, err
	}

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(imageData.Rect.Size().X),
		int32(imageData.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(imageData.Pix),
	)
	// gl.GenerateTextureMipmap(texture) // FIXME: Está gerando panic

	// Unbind
	gl.BindTexture(uint32(texType), 0)

	return &Texture{
		id:    id,
		xtype: texType,
		slot:  slotType,
	}, nil
}

func loadImage(path string) (*image.RGBA, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", path, err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}

	// FIXME: Images are being loaded upside-down
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
}

type Texture struct {
	id    uint32
	xtype int
	slot  uint32
}

func (t *Texture) ActiveAndBind() {
	gl.ActiveTexture(t.slot)
	t.Bind()
}

func (t *Texture) Bind() {
	gl.BindTexture(uint32(t.xtype), t.id)
}

func (t *Texture) Unbind() {
	gl.BindTexture(uint32(t.xtype), 0)
}

func (t *Texture) Delete() {
	gl.DeleteTextures(1, &t.id)
}
