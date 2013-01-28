package Engine

import (
	"errors"
	"github.com/vova616/gl"
	"image"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Padding = 1

type AnimatedUV []UV
type ID interface{}

type Atlas interface {
	GLTexture
	Index(id ID) image.Rectangle
	Group(id ID) []image.Rectangle
}

type UV struct {
	U1, V1, U2, V2, Ratio float32
}

func NewUV(u1, v1, u2, v2, ratio float32) UV {
	return UV{u1, v1, u2, v2, ratio}
}

func AnimatedUVs(a Atlas, ids ...ID) AnimatedUV {
	uvs := make([]UV, len(ids))
	for i, id := range ids {
		uvs[i] = IndexUV(a, id)
	}
	return uvs
}

func AnimatedGroupUVs(a Atlas, groups ...ID) (AnimatedUV, map[ID][2]int) {
	uvs := make([]UV, 0, len(groups))
	indecies := make(map[ID][2]int)
	last := 0
	for _, id := range groups {
		var ic [2]int
		uvs = append(uvs, IndexGroupUV(a, id)...)
		ic[0] = last
		ic[1] = len(uvs)
		last = ic[1]
		indecies[id] = ic
	}
	return uvs, indecies
}

func AtlasLoadDirectory(path string) (*ManagedAtlas, error) {
	d, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer d.Close()
	ds, e := d.Stat()
	if e != nil {
		return nil, e
	}
	if !ds.IsDir() {

		return nil, errors.New("The path is not a directory. " + path)
	}

	atlas := NewManagedAtlas(1024, 512)

	return atlas, atlas.LoadGroup(path)
}

func IndexUV(a Atlas, id ID) UV {
	rect := a.Index(id)
	h := float64(a.Height())
	w := float64(a.Width())
	return NewUV(float32(float64(rect.Min.X)/w),
		float32(float64(rect.Min.Y)/h),
		float32(float64(rect.Max.X)/w),
		float32(float64(rect.Max.Y)/h),
		float32(float64(rect.Dx())/float64(rect.Dy())))
}

func IndexGroupUV(a Atlas, group ID) AnimatedUV {
	rects := a.Group(group)
	uvs := make([]UV, len(rects))
	for i, r := range rects {
		h := float64(a.Height())
		w := float64(a.Width())
		uvs[i] = NewUV(float32(float64(r.Min.X)/w),
			float32(float64(r.Min.Y)/h),
			float32(float64(r.Max.X)/w),
			float32(float64(r.Max.Y)/h),
			float32(float64(r.Dx())/float64(r.Dy())))
	}
	return uvs
}

func RenderAtlas(a Atlas) {
	a.Bind()
	xratio := float32(a.Width()) / float32(a.Height())
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-0.5, -0.5, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f((xratio)-0.5, -0.5, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f((xratio)-0.5, 0.5, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-0.5, 0.5, 1)
	gl.End()
}

type ManagedAtlas struct {
	*Texture
	image  *image.RGBA
	uvs    map[ID]image.Rectangle
	groups map[ID][]ID
	images map[ID]image.Image
	Tree   *AtlasNode
}

type AtlasNode struct {
	Child   [2]*AtlasNode
	Rect    image.Rectangle
	ImageID ID
}

func NewAtlasNode(width, height int) *AtlasNode {
	return &AtlasNode{Rect: image.Rect(0, 0, width, height)}
}

func NewManagedAtlas(width, height int) *ManagedAtlas {
	m := &ManagedAtlas{
		image:  image.NewRGBA(image.Rect(0, 0, width, height)),
		uvs:    make(map[ID]image.Rectangle),
		groups: make(map[ID][]ID),
		images: make(map[ID]image.Image),
		Tree:   NewAtlasNode(width, height)}
	ResourceManager.Add(m)
	return m
}

/*
func NewRGBA(r image.Rectangle) (*image.RGBA, *MemHandle) {
	w, h := r.Dx(), r.Dy()
	memHandle := Allocate(4 * w * h)
	buf := memHandle.Bytes()
	ResourceManager.Add(memHandle)
	return &image.RGBA{buf, 4 * w, r}, memHandle
}
*/
func AtlasFromSheet(path string, width, height, frames int) (atlas *ManagedAtlas, groupID ID, err error) {
	fName := filepath.Base(path)
	extIndex := strings.LastIndex(fName, ".")
	if extIndex != -1 {
		fName = fName[:extIndex]
	}
	fName = atlas.nextGroupID(fName)

	file, e := os.Open(path)
	if e != nil {
		return nil, nil, e
	}
	defer file.Close()
	ds, e := file.Stat()
	if e != nil {
		return nil, nil, e
	}
	if ds.IsDir() {
		return nil, nil, errors.New("The path is not a file. " + path)
	}

	file.Close()

	img, e := LoadImage(path)
	if e != nil {
		return nil, nil, e
	}

	atlas = &ManagedAtlas{
		image:  image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy())),
		uvs:    make(map[ID]image.Rectangle),
		groups: make(map[ID][]ID),
		images: make(map[ID]image.Image),
		Tree:   NewAtlasNode(img.Bounds().Dx(), img.Bounds().Dy())}
	ResourceManager.Add(atlas)

	draw.Draw(atlas.image, atlas.image.Bounds(), img, image.Point{0, 0}, draw.Src)

	group := make([]ID, 0)

	maxx, maxy := img.Bounds().Dx()/width, img.Bounds().Dy()/height

	i := 0
	for y := 0; y < maxy; y++ {
		for x := 0; x < maxx; x++ {
			is := strconv.FormatInt(int64(i), 10)
			if i == 0 {
				is = ""
			}
			group = append(group, fName+is)
			atlas.uvs[fName+is] = image.Rect(x*width, y*width, (x+1)*width, (y+1)*width)
			i++
		}
	}
	atlas.groups[fName] = group
	return atlas, fName, nil
}

func (atlas *ManagedAtlas) Release() {
	if atlas.Texture != nil {
		atlas.Texture.Release()
	}
	if atlas.image != nil {
		if atlas.image.Pix != nil {
			atlas.image.Pix = nil
		}
		atlas.image = nil
	}
	atlas.uvs = nil
	atlas.groups = nil
	atlas.images = nil
	atlas.Tree = nil
}

func (node *AtlasNode) Insert(img image.Image, id ID) *AtlasNode {
	if node.Child[0] != nil {
		newNode := node.Child[1].Insert(img, id)
		if newNode != nil {
			return newNode
		}
		newNode = node.Child[0].Insert(img, id)
		if newNode != nil {
			return newNode
		}
	} else {
		if node.ImageID != nil {
			return nil
		}

		dw := node.Rect.Dx() - (img.Bounds().Dx() + Padding)
		dh := node.Rect.Dy() - (img.Bounds().Dy() + Padding)

		if dw < 0 ||
			dh < 0 {
			return nil
		}

		if dw == 0 && dh == 0 {
			return node
		}

		node.Child[0] = &AtlasNode{}
		node.Child[1] = &AtlasNode{}

		if dw > dh {
			node.Child[0].Rect = image.Rect(
				node.Rect.Min.X, node.Rect.Min.Y,
				node.Rect.Min.X+dw, node.Rect.Max.Y)

			node.Child[1].Rect = image.Rect(
				node.Rect.Min.X+dw, node.Rect.Min.Y,
				node.Rect.Max.X, node.Rect.Max.Y)
		} else {
			node.Child[0].Rect = image.Rect(
				node.Rect.Min.X, node.Rect.Min.Y,
				node.Rect.Max.X, node.Rect.Min.Y+dh)

			node.Child[1].Rect = image.Rect(
				node.Rect.Min.X, node.Rect.Min.Y+dh,
				node.Rect.Max.X, node.Rect.Max.Y)
		}
		//log.Println(node.Child[0].Rect, node.Child[1].Rect, id)
		return node.Child[1].Insert(img, id)
	}
	return nil
}

func (atlas *ManagedAtlas) LoadGIF(path string) (err error, groupID ID) {
	fName := filepath.Base(path)
	extIndex := strings.LastIndex(fName, ".")
	if extIndex != -1 {
		fName = fName[:extIndex]
	}
	fName = atlas.nextGroupID(fName)

	file, e := os.Open(path)
	if e != nil {
		return e, nil
	}
	defer file.Close()
	ds, e := file.Stat()
	if e != nil {
		return e, nil
	}
	if ds.IsDir() {
		return errors.New("The path is not a file. " + path), nil
	}

	file.Close()

	imgs, e := LoadGIF(path)
	if e != nil {
		return e, nil
	}

	group := make([]ID, 0)

	for i, img := range imgs {
		is := strconv.FormatInt(int64(i), 10)
		if i == 0 {
			is = ""
		}
		group = append(group, fName+is)
		atlas.AddImage(img, fName+is)
		i++
	}
	atlas.groups[fName] = group
	return nil, fName
}

func (atlas *ManagedAtlas) LoadGroupSheet(path string, width, height, frames int) (err error, groupID ID) {
	return atlas.LoadGroupSheetOffset(path, image.Pt(0, 0), width, height, frames)
}

func (atlas *ManagedAtlas) LoadGroupSheetOffset(path string, pt image.Point, width, height, frames int) (err error, groupID ID) {
	fName := filepath.Base(path)
	extIndex := strings.LastIndex(fName, ".")
	if extIndex != -1 {
		fName = fName[:extIndex]
	}
	fName = atlas.nextGroupID(fName)

	file, e := os.Open(path)
	if e != nil {
		return e, nil
	}
	defer file.Close()
	ds, e := file.Stat()
	if e != nil {
		return e, nil
	}
	if ds.IsDir() {
		return errors.New("The path is not a file. " + path), nil
	}

	file.Close()

	img, e := LoadImage(path)
	if e != nil {
		return e, nil
	}

	group := make([]ID, 0)

	x, y := 0, 0
	for i := 0; i < frames; i++ {
		is := strconv.FormatInt(int64(i), 10)
		if i == 0 {
			is = ""
		}
		sprite := image.NewRGBA(image.Rect(0, 0, width, height))
		point := image.Point{(x * width) + pt.X, (y * height) + pt.Y}
		if sprite.Rect.Add(point).In(img.Bounds()) {
			draw.Draw(sprite, sprite.Rect, img, point, draw.Src)
			atlas.AddImage(sprite, fName+is)
			group = append(group, fName+is)
			x++
		} else {
			x = 0
			y++
		}
	}
	atlas.groups[fName] = group
	return nil, fName
}

func (atlas *ManagedAtlas) LoadGroup(path string) error {
	d, e := os.Open(path)
	if e != nil {
		return e
	}
	defer d.Close()
	ds, e := d.Stat()
	if e != nil {
		return e
	}
	if !ds.IsDir() {

		return errors.New("The path is not a directory. " + path)
	}

	files, er := d.Readdir(0)
	for _, file := range files {
		fullName := file.Name()
		fullName = filepath.Base(fullName)
		extIndex := strings.LastIndex(fullName, ".")

		fName := fullName
		ext := ""
		if extIndex != -1 {
			fName = fullName[:extIndex]
			ext = fullName[extIndex:]
		}

		nIndex := strings.LastIndex(fName, "_")
		if nIndex == -1 {
			fload := d.Name() + "/" + fullName

			img, e := LoadImage(fload)
			if e != nil {
				log.Println(e)
				continue
			}

			atlas.AddImage(img, fName)
			group := make([]ID, 1)
			group[0] = fName

			fulldir := d.Name() + "/" + fName + "_"
			for i := 0; ; i++ {
				is := strconv.FormatInt(int64(i), 10)
				fload = fulldir + is + ext
				f, e := os.Open(fload)
				if e == nil {
					//log.Println(fload)
					f.Close()
					img, e := LoadImage(fload)
					if e != nil {
						log.Println(e)
						continue
					}
					atlas.AddImage(img, fName+is)
					group = append(group, fName+is)
				} else {
					if i > 1 {
						break
					}
				}
			}
			atlas.groups[fName] = group
		}
	}
	return er
}

func (ma *ManagedAtlas) LoadImage(path string) (err error, id ID) {

	img, e := LoadImage(path)

	if e != nil {
		return e, nil
	}

	fName := filepath.Base(path)
	extIndex := strings.LastIndex(fName, ".")
	if extIndex != -1 {
		fName = fName[:extIndex]
	}

	id = ma.nextImageID(fName)
	ma.images[id] = img

	return nil, id
}

func (ma *ManagedAtlas) nextImageID(name string) string {
	for i := -1; i < 9999; i++ {
		is := strconv.FormatInt(int64(i), 10)
		if i == -1 {
			is = ""
		}

		id := name + is
		_, exist := ma.images[id]
		if !exist {
			return id
		}
	}
	panic("Cannot choose id.")
}

func (ma *ManagedAtlas) nextGroupID(name string) string {
	for i := -1; i < 9999; i++ {
		is := strconv.FormatInt(int64(i), 10)
		if i == -1 {
			is = ""
		}

		id := name + is
		_, exist := ma.groups[id]
		if !exist {
			return id
		}
	}
	panic("Cannot choose id.")
}

func (ma *ManagedAtlas) LoadImageID(path string, id ID) error {
	_, exist := ma.images[id]
	if exist {
		errors.New("id already exists")
	}

	img, e := LoadImage(path)

	if e != nil {
		return e
	}

	ma.images[id] = img

	return nil
}

func (ma *ManagedAtlas) AddImage(img image.Image, id ID) error {
	if img == nil {
		return errors.New("image is nil")
	}
	_, exist := ma.images[id]
	if exist {
		return errors.New("id already exists")
	}

	ma.images[id] = img

	return nil
}

func (ma *ManagedAtlas) Group(id ID) []image.Rectangle {
	rects, exist := ma.groups[id]
	if exist {
		m := make([]image.Rectangle, 0, len(rects))
		for _, id := range rects {
			m = append(m, ma.Index(id))
		}
		return m
	}
	return nil
}

func (ma *ManagedAtlas) Index(id ID) image.Rectangle {
	rect, exist := ma.uvs[id]
	if exist {
		return rect
	}
	return image.Rectangle{}
}

func (ma *ManagedAtlas) Indexs() []ID {
	images := make([]ID, 0, len(ma.uvs))
	for key, _ := range ma.uvs {
		images = append(images, key)
	}
	return images
}

func (ma *ManagedAtlas) BuildAtlas() error {
	for {
		maxArea := 0
		var bigImage image.Image = nil
		var bigID ID

		for id, img := range ma.images {
			if img != nil {
				area := img.Bounds().Dx() * img.Bounds().Dy()
				if area > maxArea {
					maxArea = area
					bigImage = img
					bigID = id
				}
			}
		}

		if bigImage == nil {
			break
		}

		node := ma.Tree.Insert(bigImage, bigID)
		var rect image.Rectangle
		if node != nil {
			rect = node.Rect
			//rect.Min.X += Padding
			//rect.Min.Y += Padding
			rect.Max.X -= Padding
			rect.Max.Y -= Padding
			node.ImageID = bigID
			draw.Draw(ma.image, rect, bigImage, image.ZP, draw.Src)
		} else {
			return errors.New("not enough space in atlas")
		}

		ma.uvs[bigID] = rect
		ma.images[bigID] = nil
	}
	ma.Texture = NewRGBATexture(ma.image.Pix, ma.image.Bounds().Dx(), ma.image.Bounds().Dy())
	ma.image.Pix = nil
	ma.image = nil
	return nil
}
