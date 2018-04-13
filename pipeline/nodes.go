package pipeline

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

var startTime time.Time

// channel加上buffer_size, 可以提高性能
const BUFFER_SIZE = 1024

//初始化开始时间（用于计时）
func Init() {
	startTime = time.Now()
}

//ArraySource 一个数据源-数列
func ArraySource(a ...int) <-chan int { //用户只能拿东西
	out := make(chan int, BUFFER_SIZE)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

//InMemSort 内存中排序
func InMemSort(in <-chan int) <-chan int {
	out := make(chan int, BUFFER_SIZE)
	go func() {
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read done:", time.Since(startTime))
		//sort
		sort.Ints(a)
		fmt.Println("InMemSort done:", time.Since(startTime))

		// Output
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

//Merge 归并节点
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, BUFFER_SIZE)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 < v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		fmt.Println("Merge done:", time.Since(startTime))

	}()
	return out
}

//ReaderSource 读取数据源  因为分块，保证不能超过chunkSize
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {

	// channel加上buffer_size, 可以提高性能
	out := make(chan int, BUFFER_SIZE)
	go func() {
		buffer := make([]byte, 8) //64bit int, 8个int
		bytesRead := 0
		for {
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

//WriterSink 把结果写入到某个地方
func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

//RandomSource 一个数据源(输出count个随机数)
func RandomSource(count int) <-chan int {
	out := make(chan int, BUFFER_SIZE)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

//MergeN 合并N个输入channel
func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}
