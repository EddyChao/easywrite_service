package util

import (
	"fmt"
	"io"
)

func CopyReader(reader io.Reader) (io.Reader, io.Reader) {
	reader1, writer1 := io.Pipe()
	reader2, writer2 := io.Pipe()

	go func() {
		// 从输入流中读取数据，并写入两个输出流中
		_, err := io.Copy(io.MultiWriter(writer1, writer2), reader)
		if err != nil {
			fmt.Println("Error copying data:", err)
		}

		writer1.Close()
		writer2.Close()
	}()

	return reader1, reader2
}
