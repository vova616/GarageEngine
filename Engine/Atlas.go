package Engine

import (
	"image"
	"github.com/banthar/gl"
	"image/draw"
	"log"
	"os"
	"errors"
	"strings"
	"strconv"
)	

type Atlas interface {
	GLTexture
	Index(id interface{}) image.Rectangle
	Group(id interface{}) []image.Rectangle
} 

type UV struct {
	 U1,V1,U2,V2,Ratio float32
}

func NewUV(u1,v1,u2,v2,ratio float32) UV {
	return UV{u1,v1,u2,v2,ratio}
}

type AnimatedUV []UV

func AnimatedUVs(a Atlas, ids ...interface{}) AnimatedUV {
	uvs := make([]UV, len(ids))
	for i,id := range ids {
		uvs[i] = IndexUV(a, id)
	}
	return uvs
}

func AnimatedGroupUVs(a Atlas, groups ...interface{}) (AnimatedUV,map[interface{}][2]int) {
	uvs := make([]UV, 0, len(groups))
	indecies := make(map[interface{}][2]int)
	last := 0
	for _,id := range groups {
		var ic [2]int
		uvs = append(uvs, IndexGroupUV(a, id)...)
		ic[0] = last
		ic[1] = len(uvs)
		last = ic[1]
		indecies[id] = ic
	}
	return uvs,indecies
}

func AtlasLoadDirectory(path string) (*ManagedAtlas, error) {
	d,e := os.Open(path)
	if e != nil {
		return nil,e
	}
	defer d.Close()
	ds, e := d.Stat()
	if e != nil {
		return nil,e
	}
	if !ds.IsDir() {
		
		return nil,errors.New("The path is not a directory. " + path)
	}
	
	atlas := NewManagedAtlas(2048,512)
	 
	files, er := d.Readdir(0) 
	for _,file := range files {
		fullName := file.Name()
		extIndex := strings.LastIndex(fullName, ".")
		if extIndex == -1 {
			continue
		}
		fName := fullName[:extIndex]
		ext := fullName[extIndex:]
		
		nIndex := strings.LastIndex(fName, "_")
		if nIndex == -1 {
			fload := d.Name() + "/" + fullName
			
			img,e := LoadImage(fload)
			if e != nil {
				log.Println(e)
				continue
			}
			
			atlas.AddImage(img, fName)
			group := make([]interface{}, 1)
			group[0] = fName
			
			fulldir := d.Name() + "/" + fName + "_"
			for i:=0;;i++ {
				is := strconv.FormatInt(int64(i), 10)
				fload = fulldir + is + ext
				f,e := os.Open(fload)
				if e == nil {
					//log.Println(fload)
					f.Close()
					img,e := LoadImage(fload)
					if e != nil {
						log.Println(e)
						continue
					}
					atlas.AddImage(img, fName + is)
					group = append(group, fName + is)
				} else {
					if i > 1 {
						break
					}
				}
			}
			atlas.groups[fName] = group
		}
	}
	//if e != nil {
	//	println(e)
	//}
	
	return atlas,er
}

func IndexUV(a Atlas, id interface{}) UV {
	rect := a.Index(id)
	h := float32(a.Height())
	w := float32(a.Width())
	return NewUV(float32(rect.Min.X) / w, float32(rect.Min.Y) / h, float32(rect.Max.X) / w, float32(rect.Max.Y) / h, float32(rect.Dx())/float32(rect.Dy())) 
}

func IndexGroupUV(a Atlas, group interface{}) AnimatedUV {
	rects := a.Group(group)
	uvs := make([]UV, len(rects))
	for i,r := range rects {
		h := float32(a.Height())
		w := float32(a.Width())
		uvs[i] = NewUV(float32(r.Min.X) / w, float32(r.Min.Y) / h, float32(r.Max.X) / w, float32(r.Max.Y) / h, float32(r.Dx())/float32(r.Dy()))
	}
	return uvs
}


func RenderAtlas(a Atlas) {
	a.Bind()
	xratio := float32(a.Width()) / float32(a.Height())
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1); gl.Vertex3f(-0.5, -0.5, 1) 
	gl.TexCoord2f(1, 1); gl.Vertex3f((xratio)-0.5, -0.5, 1) 
	gl.TexCoord2f(1, 0); gl.Vertex3f((xratio)-0.5, 0.5, 1) 
	gl.TexCoord2f(0, 0); gl.Vertex3f(-0.5, 0.5, 1) 
	gl.End()
}

type ManagedAtlas struct {
	 *Texture
	 image  *image.RGBA
	 lastPoint image.Point
	 images map[interface{}] image.Rectangle
	 groups map[interface{}] []interface{}
}

func NewManagedAtlas(width, height int) *ManagedAtlas {
	return &ManagedAtlas{image: image.NewRGBA(image.Rect(0,0,width,height)), images: make(map[interface{}] image.Rectangle),groups:make(map[interface{}] []interface{}) , lastPoint: image.Pt(1,1)}
}

func (ma *ManagedAtlas) AddImage(img image.Image, id interface{} ) {
	if img == nil {
		panic("img is null")
	}
	_, exist := ma.images[id] 
	if exist {
		panic("id already exists")
	}
	
	bd := img.Bounds().Add(ma.lastPoint)
	
	if bd.Max.X + 1 <= ma.image.Bounds().Max.X && bd.Max.Y + 1 <= ma.image.Bounds().Max.Y{
		draw.Draw(ma.image, bd, img, image.Pt(0,0), draw.Over)
		ma.lastPoint.X = bd.Max.X + 1
			
		ma.images[id] = bd
	} else {
		for i:=0;i<len(ma.images);i++ {
			rect := ma.images[i]
			log.Println(ma.image.Bounds().Max.Y - (rect.Max.Y + 1), ma.image.Bounds().Max.X - (rect.Min.X + 1), img.Bounds().Max.Y, img.Bounds().Max.X)
			if ma.image.Bounds().Max.Y - (rect.Max.Y + 1) >= img.Bounds().Max.Y + 1 && ma.image.Bounds().Max.X - (rect.Min.X + 1) >= img.Bounds().Max.X + 1 {
				ma.lastPoint.Y = rect.Max.Y + 1
				ma.lastPoint.X = rect.Min.X 
				
				ma.AddImage(img, id)
				break
			} 
		}
		panic("Not enough room in atlas")
	}
}

func (ma *ManagedAtlas) Group(id interface{}) []image.Rectangle {
	rects, exist := ma.groups[id] 
	if exist {
		m := make([]image.Rectangle, 0, len(rects))
		for _,id := range rects {
			m = append(m, ma.Index(id))
		}
		return m
	}
	return nil
}

func (ma *ManagedAtlas) Index(id interface{}) image.Rectangle {
	rect, exist := ma.images[id] 
	if exist {
		return rect
	}
	return image.Rectangle{}
}

func (ma *ManagedAtlas) Indexs() []interface{} {
	images := make([]interface{}, 0, len(ma.images))
	for key,_ := range ma.images {
		images = append(images, key)
	}
	return images
}

func (ma *ManagedAtlas) BuildAtlas() {
	//log.Println(len(ma.image.Pix), ma.image.Bounds().Dx()*ma.image.Bounds().Dy())
	//t,err := LoadTextureFromImage(ma.image)
	//if err != nil {
	//	panic(err)
	//}
	//ma.Texture = t
	ma.Texture = NewRGBATexture(ma.image.Pix, ma.image.Bounds().Dx(), ma.image.Bounds().Dy())
	ma.image.Pix = nil
}

