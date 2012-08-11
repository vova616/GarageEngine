package Engine

type Camera struct {
	BaseComponent
	Projection *Matrix
}

func NewCamera() *Camera {
			
	proj := NewIdentity()
	proj.Ortho(0, float32(Width), 0, float32(Height), -1000, 1000) 
		
	return &Camera{NewComponent(), proj}
}

func (c *Camera) UpdateResolution() {
	c.Projection.Ortho(0, float32(Width), 0, float32(Height), -1000, 1000) 
}

func (c *Camera) MouseRealPosition(x,y int) Vector {
	d := NewIdentity()
	s := c.Transform().WorldScale()
	
	d.Mul(c.Transform().Matrix())
	d.Scale(-1,-1,0)
	d.Translate(float32(x),float32(y),0) 
	d.Scale(1/s.X,1/s.Y,0)
	
	return d.Translation()
}