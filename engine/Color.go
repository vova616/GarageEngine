package engine

type Color struct {
	R, G, B, A float32
}

var (
	Color_White       = Color{1, 1, 1, 1}
	Color_Black       = Color{0, 0, 0, 1}
	Color_Red         = Color{1, 0, 0, 1}
	Color_Green       = Color{0, 1, 0, 1}
	Color_Blue        = Color{0, 0, 1, 1}
	Color_Transparent = Color{1, 1, 1, 0}
)
