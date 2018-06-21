package main

import (
	"errors"
	"github.com/dist-ribut-us/log"
	gl "log"
)

func main() {
	log.Go()
	log.Info("this is a test")
	log.Info(log.Line())
	gl.Print("Hi Go")
	log.Error(nil)
	log.PathDepth = 0
	log.Error(errors.New("just an error"))
	A()

	log.PathDepth = 1
	A()
}

func A() {
	B()
}

func B() {
	C()
}

func C() {
	log.Info(log.Lbl("Trace -->"), log.Line())
	log.Info(log.Lbl("One line up"), log.Trace(0, 1))
}
