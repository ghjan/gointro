package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ghjan/gointro/pipeline"
)

var fileName = "large.in"
var fileNameOutput = "large.out"

const FILE_SIZE = 800000000
const CHUNK_COUNT = 4

func main() {
	p := createPipeline(fileName, FILE_SIZE, CHUNK_COUNT)
	writeToFile(p, fileNameOutput)
	printFile(fileNameOutput, 100)
}

//createPipeline 创建流水线
// @TODO return File *, which should be closed after process
func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	pipeline.Init()
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
func printFile(fileName string, count int) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	p := pipeline.ReaderSource(file, -1)
	i := 0
	for v := range p {
		i++
		fmt.Println(v)
		if i >= count {
			break
		}
	}
}
