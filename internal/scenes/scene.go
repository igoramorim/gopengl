package scenes

type Scene interface {
	Name() string
	Width() int
	Height() int
	Show()
}
