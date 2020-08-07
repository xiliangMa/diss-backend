package models

type GoRoutineManager struct {
	GoRoutineMap map[string]interface{}
}

func NewGoRoutineManager() *GoRoutineManager {
	return &GoRoutineManager{
		GoRoutineMap: make(map[string]interface{}),
	}
}
