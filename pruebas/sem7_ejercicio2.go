package main

import (
	"fmt"
	"sync"
	"time"
)

type Student struct {
	name   string
	mutex  *sync.Mutex
	exam   chan bool
}

type Collector struct {
	name   string
	mutex  *sync.Mutex
	exam   chan bool
}

func (s *Student) takeExam() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Printf("%s toma el examen.\n", s.name)
	time.Sleep(time.Second * 5)
	s.exam <- true
}

func (c *Collector) collectExams() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	fmt.Printf("%s recoge el examen.\n", c.name)
	time.Sleep(time.Second * 5)
	<-c.exam
}

func main() {
	exam := make(chan bool, 5)

	students := make([]*Student, 5)
	collectors := make([]*Collector, 3)

	for i := 0; i < 5; i++ {
		students[i] = &Student{
			name:  fmt.Sprintf("Estudiante %d", i+1),
			mutex: &sync.Mutex{},
			exam:  exam,
		}
	}

	for i := 0; i < 3; i++ {
		collectors[i] = &Collector{
			name:  fmt.Sprintf("Encargado %d", i+1),
			mutex: &sync.Mutex{},
			exam:  exam,
		}
	}

	var wg sync.WaitGroup
	wg.Add(5)
	for _, student := range students {
		go func(s *Student) {
			defer wg.Done()
			s.takeExam()
		}(student)
	}

	wg.Add(3)
	for _, collector := range collectors {
		go func(c *Collector) {
			defer wg.Done()
			c.collectExams()
		}(collector)
	}

	wg.Wait()
}