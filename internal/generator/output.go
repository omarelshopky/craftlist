package generator

import (
	"bufio"
	"fmt"
	"os"
)

type OutputManager struct{}

type OutputWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewOutputManager() *OutputManager {
	return &OutputManager{}
}

func (om *OutputManager) CreateWriter(filename string) (*OutputWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	writer := bufio.NewWriter(file)

	return &OutputWriter{
		file:   file,
		writer: writer,
	}, nil
}

func (ow *OutputWriter) WritePassword(password string) error {
	_, err := fmt.Fprintln(ow.writer, password)
	return err
}

func (ow *OutputWriter) Flush() error {
	return ow.writer.Flush()
}

func (ow *OutputWriter) Close() error {
	if err := ow.writer.Flush(); err != nil {
		ow.file.Close()
		return err
	}
	return ow.file.Close()
}