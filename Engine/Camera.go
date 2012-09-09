package Engine

type Camera struct {
	BaseComponent
	Projection *Matrix
}

func NewCamera() *Camera {

	proj := NewIdentity()
	w := float32(Width)
	h := float32(Height)
	//proj.Ortho(-w, w, -h, h, -1000, 1000) 
	proj.Ortho(0, w, 0, h, -1000, 1000)

	return &Camera{NewComponent(), proj}
}

func (c *Camera) Update() {
	//c.BaseComponent2.
	/*
		w := float32(Width)/2
		h := float32(Height)/2
		proj := NewIdentity()
		proj.Ortho(-w, w, -h, h, -1000, 1000) 
		c.Projection = proj
	*/
}

func (c *Camera) UpdateResolution() {
	c.Projection.Ortho(0, float32(Width), 0, float32(Height), -1000, 1000)
}

func (c *Camera) MouseRealPosition(x, y int) Vector {
	d := NewIdentity()
	s := c.Transform().WorldScale()

	d.Mul(c.Transform().Matrix())
	d.Scale(-1, -1, 0)
	d.Translate(float32(x), float32(y), 0)
	d.Scale(1/s.X, 1/s.Y, 0)

	return d.Translation()
}

func (c *Camera) Render() {
	s := GetScene()
	if s != nil {
		tcam := s.SceneBase().Camera
		s.SceneBase().Camera = c
		arr := s.SceneBase().gameObjects
		if arr == nil {
			println("arr")
		}
		if c.GameObject() == nil {
			println("c.GameObject()")
		}

		IterExcept(arr, drawGameObject, c.GameObject())
		s.SceneBase().Camera = tcam
	}
}
