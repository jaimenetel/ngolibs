package nhttp

import (
	"net/http"
	"sync"
)

// CustomMux es un multiplexor HTTP que permite agregar y eliminar rutas dinámicamente
type CustomMux struct {
	Handlers map[string]http.Handler
	mu       sync.RWMutex
}

var instancemux *CustomMux
var oncemux sync.Once

// GetCustomMux retorna la instancia única del CustomMux
func GetCustomMux() *CustomMux {
	oncemux.Do(func() {
		instancemux = &CustomMux{
			Handlers: make(map[string]http.Handler),
		}
	})
	return instancemux
}

// NewCustomMux crea un nuevo CustomMux
func __NewCustomMux() *CustomMux {
	return &CustomMux{
		Handlers: make(map[string]http.Handler),
	}
}
func (c *CustomMux) ToEndPointList() []string {
	var list []string
	for k := range c.Handlers {
		list = append(list, k)
	}
	return list
}

// ServeHTTP maneja las solicitudes HTTP
func (c *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mu.RLock()
	handler, ok := c.Handlers[r.URL.Path]
	c.mu.RUnlock()

	if ok {
		handler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// Handle agrega una nueva ruta al CustomMux
func (c *CustomMux) Handle(pattern string, handler http.Handler) {
	c.mu.Lock()
	c.Handlers[pattern] = handler
	c.mu.Unlock()
}

// Remove elimina una ruta del CustomMux
func (c *CustomMux) Remove(pattern string) {
	c.mu.Lock()
	delete(c.Handlers, pattern)
	c.mu.Unlock()
}
