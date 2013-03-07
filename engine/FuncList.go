package engine

type FuncList struct {
	Functions []func()
}

func NewFuncList(f func()) *FuncList {
	fa := &FuncList{nil}
	fa.Add(f)
	return fa
}

func (this *FuncList) Run() {
	for _, fnc := range this.Functions {
		fnc()
	}
}

func (this *FuncList) Add(fnc func()) (index int) {
	if fnc == nil {
		return
	}
	this.Functions = append(this.Functions, fnc)
	return len(this.Functions) - 1
}

func (this *FuncList) Remove(index int) {
	this.Functions[len(this.Functions)-1], this.Functions[index], this.Functions = nil, this.Functions[len(this.Functions)-1], this.Functions[:len(this.Functions)-1]
}
