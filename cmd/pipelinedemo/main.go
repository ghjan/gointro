package main

import (
	"github.com/ghjan/gointro/pipeline"
	"fmt"
)

func main() {
	//demo1()
	demo2()
}
func demo1() {
	p := pipeline.ArraySource(3, 2, 6, 7, 4)
	//for {
	//	if num, ok := <-p; ok {
	//		fmt.Println(num)
	//	}else{
	//		break
	//	}
	//}
	//下面是比较简便的写法 使用range，但是发送方里面的p里面必须close，否则就不知道什么时候结束了
	for v := range p {
		fmt.Println(v)
	}
}

func demo2() {
	p := pipeline.Merge(pipeline.InMemSort(pipeline.ArraySource(3, 2, 6, 7, 4)),
		pipeline.InMemSort(pipeline.ArraySource(7, 4, 0, 3, 2, 13, 8)))
	for v := range p {
		fmt.Println(v)
	}
}
