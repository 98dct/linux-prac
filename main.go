package main

import (
	"errors"
	"fmt"
	xerr "github.com/pkg/errors"
	"math"
	"os"
	"runtime"
	"time"
)

func test1() {
	_, err := os.OpenFile("aaa", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(math.MinInt)
	fmt.Println(math.MaxInt)

	dst := append([]int(nil), []int{1, 2, 3}...)
	fmt.Println(dst)
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d KB\n", m.Alloc/1024)
}

func getMessage(msg []byte) []byte {
	return msg[:5]
}

func test2() {
	// 打印初始内存分配
	fmt.Print("Initial memory allocation: ")
	printAlloc()

	// 创建大量临时数据以增加内存分配

	data := make([]byte, 1000000)
	data1 := getMessage(data)
	fmt.Println(len(data1))

	fmt.Print("Memory allocation after creating data: ")
	printAlloc()
	// 强制触发垃圾回收
	runtime.GC()

	// 打印垃圾回收后的内存分配
	fmt.Print("Memory allocation after GC: ")
	printAlloc()

	//runtime.KeepAlive(data1)
}

type Foo struct {
	v []byte
}

func test3() {
	foos := make([]Foo, 1_000)
	printAlloc()

	for i := 0; i < len(foos); i++ {
		foos[i] = Foo{v: make([]byte, 1024*1024)}
	}

	printAlloc()

	two := keepFirstTwoElementsOnly(foos)
	runtime.GC()
	printAlloc()
	runtime.KeepAlive(two)
}

func test4() {

	ch := make(chan string, 3)
	ch <- "aaa"
	ch <- "bbb"
	ch <- "ccc"
	close(ch)

	for v := range ch {
		func() {
			fmt.Println(v)
		}()
	}

}

func test5() {

	sli := []string{"aaa", "bbb", "ccc"}
	for _, v := range sli {
		fmt.Println(&v)
		go func() {
			fmt.Println(v)
		}()
	}
	time.Sleep(1 * time.Second)

}

func test6() {

	err1 := errors.New("错误1")
	err2 := xerr.Wrap(err1, "附加信息")
	fmt.Println(err2)
	flag := errors.Is(err2, err1)
	fmt.Println(flag)
	err3 := fmt.Errorf("错误包装：%w", err1)
	fmt.Println(err3)
	flag2 := errors.Is(err3, err1)
	fmt.Println(flag2)
}

func test7() {

	i := 0
	ch := make(chan struct{}, 1)
	go func() {
		i = 1
		<-ch
	}()

	ch <- struct{}{}
	fmt.Println(i)

}

func test8() {

	i := 0
	ch := make(chan struct{})
	go func() {
		i = 1
		<-ch
	}()

	ch <- struct{}{}
	fmt.Println(i)

}

func main() {
	//test2()
	//test3()
	//test4()
	//test5()
	//test6()
	test7()
	test8()
}

func keepFirstTwoElementsOnly(foos []Foo) []Foo {
	return foos[:2]
}
