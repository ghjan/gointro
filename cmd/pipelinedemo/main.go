package main

import (
	"github.com/ghjan/gointro/pipeline"
	"fmt"
	"os"
	"log"
	"bufio"
)

const filename = "small.in"
const n = 100000000

func main() {
	//demo1()
	//demo2()
	demo3()
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

func demo3() {
	file, err := os.Create(filename)
	if err != nil {
		//不知道怎么办 非产品 我们不需要处理他
		log.Panic(err)
	}
	//最后执行，就算异常退出也执行。有点类似java的finally
	defer file.Close()

	p := pipeline.RandomSource(n)
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, p)
	writer.Flush()
	file, err = os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(bufio.NewReader(file),-1)
	count := 0
	for v := range p {
		fmt.Println(v)

		count++
		if count > 100 {
			break
		}
	}

}
