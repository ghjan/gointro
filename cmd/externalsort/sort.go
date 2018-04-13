package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ghjan/gointro/pipeline"
)

var fileName = "small.in"
var fileNameOutput = "small.out"

func main() {
	p := createPipeline(fileName, 512, 4)
	writeToFile(p, fileNameOutput)
	printFile(fileNameOutput)
}

//createPipeline 创建流水线
// @TODO return File *, which should be closed after process
func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	sortResults := []<-chan int{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortResults = append(sortResults, pipeline.InMemSort(source))

	}
	return pipeline.MergeN(sortResults...)
}

//writeToFile 写入文件
func writeToFile(p <-chan int, filename string) {
	//不同的defer， 先进后出，后面的先执行
	file, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, p)
}

//printFile 打印文件内容
func printFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)
	for v := range p {
		fmt.Println(v)
	}
}
