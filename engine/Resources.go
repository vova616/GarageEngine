package engine

import (
	"fmt"
	//"github.com/go-gl/gl"
)

type ResID interface{}

type Resource interface {
	Release()
}

var (
	ResourceManager = &Resources{make(map[ResID]Resource)}
)

type MemHandle struct {
	Buff []byte
}

func Allocate(size int) *MemHandle {
	return &MemHandle{make([]byte, size)}
}

func (m *MemHandle) Release() {
	m.Buff = nil
}

func (m *MemHandle) Bytes() []byte {
	return m.Buff
}

type Resources struct {
	Resources map[ResID]Resource
}

func (r *Resources) Add(res Resource) error {
	_, exists := r.Resources[res]
	if exists {
		return fmt.Errorf("Cannot add res %d %v", res, res)
	} else {
		r.Resources[res] = res
	}
	return nil
}

func (r *Resources) AddManual(res Resource, key interface{}) error {
	_, exists := r.Resources[key]
	if exists {
		return fmt.Errorf("Cannot add res %d %v", res, res)
	} else {
		r.Resources[res] = res
	}
	return nil
}

func (r *Resources) Release() {
	for _, res := range r.Resources {
		releaseSafe(res)
	}
	r.Resources = make(map[ResID]Resource)
}

func (r *Resources) ReleaseResource(res interface{}) {
	resData, exists := r.Resources[res]
	if exists {
		resData.Release()
		delete(r.Resources, res)
	}
}

func releaseSafe(res Resource) {
	defer recoverFromPanic(res)
	res.Release()
}

func recoverFromPanic(res Resource) {
	if err := recover(); err != nil {
		fmt.Printf("release failed: ID:%d Res:%v Error:%v\n", res, res, err)
	}
}
