package engine

type FuncKey *func()

type FuncList struct {
	Functions []FuncKey
}

func NewFuncList(f func()) *FuncList {
	fa := &FuncList{nil}
	fa.Add(f)
	return fa
}

func (this *FuncList) Run() {
	for i, fnc := range this.Functions {
		if fnc != nil && *fnc != nil {
			(*fnc)()
		} else {
			this.Functions[len(this.Functions)-1], this.Functions[i], this.Functions = nil, this.Functions[len(this.Functions)-1], this.Functions[:len(this.Functions)-1]
			if fnc != nil {
				*fnc = nil
			}
		}
	}
}

func (this *FuncList) Add(fnc func()) (key FuncKey) {
	if fnc == nil {
		return
	}
	c := &fnc
	this.Functions = append(this.Functions, c)
	return c
}

func (this *FuncList) Remove(key FuncKey) (deleted bool) {
	if key == nil || *key == nil {
		return false
	}
	for i, c := range this.Functions {
		if c == key {
			this.Functions[len(this.Functions)-1], this.Functions[i], this.Functions = nil, this.Functions[len(this.Functions)-1], this.Functions[:len(this.Functions)-1]
			*key = nil
			return true
		}
	}
	return false
}
