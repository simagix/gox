// Copyright 2019 Kuei-chun Chen. All rights reserved.

package gox

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/snappy"
)

func TestNewFileReader(t *testing.T) {
	var err error
	var file *os.File
	filename := "/tmp/count.file"
	if file, err = os.Create(filename); err != nil {
		t.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	str := "keyhole"
	b := []byte(str)
	writer.Write(b)
	writer.Flush()
	file.Close()

	reader, _ := NewFileReader(filename)
	buf, _, _ := reader.ReadLine()

	if str != string(buf) {
		t.Fatal(string(buf))
	}
}

func TestNewReader(t *testing.T) {
	var err error
	// var file *os.File
	filename := "keyhole.gz"
	str := "keyhole"
	OutputGzipped([]byte(str), filename)
	gzreader, err := NewFileReader(filename)
	if err != nil {
		t.Fatal(err)
	}
	buf, _, err := gzreader.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	if str != string(buf) {
		t.Fatal(string(buf))
	}
	os.Remove(filename)

	filename = "keyhole.sz"
	OutputSnappyZipped([]byte(str), filename)
	snzreader, err := NewFileReader(filename)
	if err != nil {
		t.Fatal(err)
	}
	buf, _, err = snzreader.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	if str != string(buf) {
		t.Fatal(str, string(buf))
	}
	os.Remove(filename)
}

func TestCountLines(t *testing.T) {
	var err error
	var file *os.File
	filename := "/tmp/count.file"
	if file, err = os.Create(filename); err != nil {
		t.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	total := 10

	for i := 0; i < total; i++ {
		writer.WriteString("value\n")
	}
	writer.Flush()
	file.Close()
	file, _ = os.Open(filename)
	defer file.Close()
	reader := bufio.NewReader(file)
	count, _ := CountLines(reader)

	if count != total {
		t.Fatal(count)
	}
}

func TestOutputGzipped(t *testing.T) {
	var err error
	var b []byte
	var fz *gzip.Reader
	var file *os.File
	filename := "/tmp/filename.gz"
	str := "This is a test line! "
	for len(str) < (10 * 1024 * 1024) {
		str += str
	}
	if err = OutputGzipped([]byte(str), filename); err != nil {
		t.Fatal(err)
	}
	if file, err = os.Open(filename); err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if fz, err = gzip.NewReader(file); err != nil {
		t.Fatal(err)
	}
	defer fz.Close()

	if b, err = ioutil.ReadAll(fz); err != nil {
		t.Fatal(err)
	}

	if string(b) != str {
		t.Fatal(err)
	}

	if err = os.Remove(filename); err != nil {
		t.Fatal(err)
	}
}

func TestOutputSnappyZipped(t *testing.T) {
	var err error
	var b []byte
	var fz *snappy.Reader
	var file *os.File
	filename := "/tmp/filename.gz"
	str := "This is a test line! "
	for len(str) < (10 * 1024 * 1024) {
		str += str
	}
	if err = OutputSnappyZipped([]byte(str), filename); err != nil {
		t.Fatal(err)
	}
	if file, err = os.Open(filename); err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if fz = snappy.NewReader(file); err != nil {
		t.Fatal(err)
	}

	if b, err = ioutil.ReadAll(fz); err != nil {
		t.Fatal(err)
	}

	if string(b) != str {
		t.Fatal(err)
	}

	if err = os.Remove(filename); err != nil {
		t.Fatal(err)
	}
}

func TestReadAll(t *testing.T) {
	var file *os.File
	var zfile *os.File
	filename := "/tmp/count.file"
	zfilename := "/tmp/count.file.gz"
	str := "keyhole"
	ioutil.WriteFile(filename, []byte(str), 0644)
	OutputGzipped([]byte(str), zfilename)
	file, _ = os.Open(filename)
	zfile, _ = os.Open(zfilename)

	b, _ := ReadAll(file)
	bz, _ := ReadAll(zfile)

	if string(b) != string(bz) {
		t.Fatal("ReadAll failed")
	}
}

func TestZipFiles(t *testing.T) {
	zipInto := "ioutil.zip"
	toBeZipped := []string{"ioutil.go", "ioutil_test.go"}
	if err := ZipFiles(zipInto, toBeZipped); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(zipInto); err != nil {
		t.Fatal(err)
	}
}
