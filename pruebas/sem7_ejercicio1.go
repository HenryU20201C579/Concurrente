package main

import (
	"fmt"
	"sync"
	"time"
)

type Student struct {
	name       string
	hours      int
	becario    *Becario
	mutex      *sync.Mutex
}

type Becario struct {
	name       string
	hours      int
	student    *Student
	mutex      *sync.Mutex
}

func (s *Student) enterRoom() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Printf("%s entra a la sala de lectura.\n", s.name)
	time.Sleep(time.Duration(s.hours) * time.Second)

	s.becario.mutex.Lock()
	defer s.becario.mutex.Unlock()

	fmt.Printf("%s espera su turno con el becario.\n", s.name)
	s.becario.student = s
	s.becario.mutex.Unlock()

	s.mutex.Unlock()
	s.becario.mutex.Lock()
	s.mutex.Lock()
}

func (b *Becario) helpStudent() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.student == nil {
		fmt.Printf("%s está en la sala de lectura.\n", b.name)
		time.Sleep(time.Duration(b.hours) * time.Second)
		return
	}

	fmt.Printf("%s ayuda a %s.\n", b.name, b.student.name)
	time.Sleep(time.Duration(b.hours) * time.Second)

	b.student.mutex.Unlock()
	b.student = nil
}

func main() {
	student := &Student{
		name:    "Estudiante",
		hours:   10,
		becario: &Becario{
			name:  "Becario",
			hours: 10,
		},
		mutex: &sync.Mutex{},
	}

	student.becario.mutex = &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		student.enterRoom()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		student.becario.helpStudent()
	}()

	wg.Wait()
}