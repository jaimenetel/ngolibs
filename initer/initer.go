package initer

import (
	"fmt"
	"sort"
	"sync"
)

// InitFunc es un tipo que define una función de inicialización.
type InitFunc func()

// prioritizedInitFunc contiene una función de inicialización y su prioridad.
type prioritizedInitFunc struct {
	function InitFunc
	priority int
}

// initManager es una estructura que contiene la lista de funciones de inicialización y la lógica singleton.
type initManager struct {
	initFuncs       []prioritizedInitFunc
	initmanageronce sync.Once
}

// instance es la única instancia de initManager.
var instance *initManager

// getInstance devuelve la única instancia de initManager, creándola si es necesario.
func GetInitManager() *initManager {
	if instance == nil {
		instance = &initManager{}
	}
	return instance
}

// RegisterInitFunc agrega una función a la lista de inicialización con una prioridad específica.
func RegisterInitFunc(f InitFunc, priority int) {
	GetInitManager().initFuncs = append(GetInitManager().initFuncs, prioritizedInitFunc{function: f, priority: priority})
}

// Initialize recorre todas las funciones de inicialización y las ejecuta una sola vez en orden de prioridad.
func Initialize() {
	GetInitManager().initmanageronce.Do(func() {
		// Ordenar las funciones por prioridad.
		sort.SliceStable(GetInitManager().initFuncs, func(i, j int) bool {
			return GetInitManager().initFuncs[i].priority > GetInitManager().initFuncs[j].priority
		})

		// Ejecutar las funciones en orden.
		for _, initFunc := range GetInitManager().initFuncs {
			initFunc.function()
		}
	})
}

// IniterSample es una función de ejemplo que muestra cómo usar el paquete initer.

func IniterSample() {
	fmt.Println("Ejecutando inicialización por primera vez:")
	Initialize()

	fmt.Println("Inicialización completa.")
}
