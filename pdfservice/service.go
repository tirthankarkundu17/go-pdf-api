package pdfservice

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/sunshineplan/imgconv"
)

type Service interface {
	CreateFromImage(ctx context.Context, imagePath string) ([]byte, error)
	CreateFromText(ctx context.Context, text string) ([]byte, error)
}

func (s *service) CreateFromImage(ctx context.Context, imagePath string) ([]byte, error) {
	src, err := imgconv.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %v", err)
	}

	var buf bytes.Buffer
	writer := io.Writer(&buf)

	// Write the resulting image as PDF.
	err = imgconv.Write(writer, src, &imgconv.FormatOption{Format: imgconv.PDF})
	if err != nil {
		return nil, fmt.Errorf("failed to write image: %v", err)
	}

	return buf.Bytes(), nil
}

func (s *service) CreateFromText(ctx context.Context, text string) ([]byte, error) {
	// TODO : To be implemented
	return nil, nil
}

type service struct {
}

func New() Service {
	return &service{}
}
