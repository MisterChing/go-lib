package confguard

import "sync"

type Guarder struct {
	child interface{}
	rw    sync.RWMutex
}

func NewGuarder(child interface{}) *Guarder {
	obj := &Guarder{
		child: child,
	}
	return obj
}

func (gdr *Guarder) Get() interface{} {
	gdr.rw.RLock()
	defer gdr.rw.RUnlock()
	return gdr.child
}

func (gdr *Guarder) Lock() {
	gdr.rw.Lock()
}

func (gdr *Guarder) Unlock() {
	gdr.rw.Unlock()
}
