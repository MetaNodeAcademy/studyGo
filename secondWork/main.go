package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 指针
func pointNum(num *int) {
	*num = *num + 10
}
func pointCut(num []int) {
	for i := range num {
		num[i] = num[i] * 2
	}
}

// Goroutine
func goroutine() {
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 != 0 {
				fmt.Println(i)
				time.Sleep(time.Millisecond * 1)
			}
		}
	}()
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Println(i)
				time.Sleep(time.Millisecond * 1)
			}
		}
	}()
}

// 面向对象
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}
func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.radius
}

type Person struct {
	name string
	age  int
}

type Employee struct {
	Person
	employeeId int
}

func (e Employee) PrintInfo() {
	fmt.Println(e)
}

// channel
func inputOnly(ch chan<- int) {
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
}
func outputOnly(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

func inputOnly2(ch chan<- int) {
	for i := 1; i <= 100; i++ {
		ch <- i
	}
	close(ch)
}
func outputOnly2(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

// 锁机制
func testLock(num *int, sy *sync.Mutex) {
	sy.Lock()
	defer sy.Unlock()
	for i := 0; i < 1000; i++ {
		*num++
	}
}

func noLock(num *int64) {
	for i := 0; i < 1000; i++ {
		atomic.AddInt64(num, 1)
	}
}

func main() {
	//指针
	num := 10
	pointNum(&num)
	fmt.Println(num)
	arr := []int{1, 2, 3, 4, 5}
	pointCut(arr)
	fmt.Println(arr)
	//goroutine
	goroutine()
	time.Sleep(time.Millisecond * 100)
	//面向对象
	r := Rectangle{width: 10, height: 5}
	fmt.Println(r.Area())
	fmt.Println(r.Perimeter())
	c := Circle{radius: 5}
	fmt.Println(c.Area())
	fmt.Println(c.Perimeter())
	e := Employee{Person: Person{name: "张三", age: 18}, employeeId: 1001}
	e.PrintInfo()
	//channel
	ch := make(chan int)
	go func() {
		outputOnly(ch)
	}()
	go func() {
		inputOnly(ch)
	}()

	ch1 := make(chan int, 5)
	go func() {
		outputOnly2(ch1)
	}()
	go func() {
		inputOnly2(ch1)
	}()
	time.Sleep(time.Millisecond * 100)
	num1 := 0
	var mu sync.Mutex
	var wg sync.WaitGroup // 添加等待组
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testLock(&num1, &mu)
		}()
	}
	wg.Wait() // 等待所有协程完成
	fmt.Println(num1)
	var noLockNum int64 = 0
	var wg2 sync.WaitGroup // 添加等待组
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			noLock(&noLockNum)
		}()
	}
	wg2.Wait()
	fmt.Println(noLockNum)
	time.Sleep(time.Millisecond * 100)
}
