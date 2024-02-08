package image

import (
	"context"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-filestorage-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-filestorage-service/pkg/logger"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
)

const FolderPath = "storage/image"

type Service struct {
	log *slog.Logger
	cfg *config.Config
}

func New(cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{cfg: cfg, log: logger}
}

func (s Service) Upload(ctx context.Context, image []byte, filename string) (string, error) {
	const op = "service.image.upload"
	log := s.log.With(slog.String("op", op))

	filePath := filepath.Join(FolderPath, filename)

	dir := filepath.Dir(FolderPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Error("failed to create directory", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Write the file
	err := os.WriteFile(filePath, image, 0644)
	if err != nil {
		log.Error("failed to save image", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf("http://%s/%s/%s", s.cfg.HTTP.Address, "image", filename), nil
}

func (s Service) Delete(ctx context.Context, imageUrl string) error {
	const op = "service.image.delete"
	log := s.log.With(slog.String("op", op))

	parsedUrl, err := url.Parse(imageUrl)
	if err != nil {
		log.Error("failed to parse image URL", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	imagePath := parsedUrl.Path

	safePath := filepath.Join(FolderPath, filepath.Base(imagePath))

	// Delete the file
	err = os.Remove(safePath)
	if err != nil {
		log.Error("failed to delete image", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
