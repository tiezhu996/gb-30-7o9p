package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/dto"
	"github.com/gbadopt/gbadopt/internal/util"
)

var allowedUploadExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".pdf": true,
}

// UploadHandler exposes MinIO-backed file upload and download.
type UploadHandler struct {
	minio  *util.MinIOClient
	logger *slog.Logger
}

// NewUploadHandler creates an UploadHandler.
func NewUploadHandler(minio *util.MinIOClient, logger *slog.Logger) *UploadHandler {
	return &UploadHandler{minio: minio, logger: logger}
}

// Upload handles POST /uploads (multipart field "file").
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "upload failed: file field required"))
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedUploadExt[ext] {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest,
			fmt.Sprintf("Upload[filename=%s] failed: unsupported extension", file.Filename)))
		return
	}
	if file.Size > 8<<20 {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "upload failed: file too large (>8MB)"))
		return
	}
	src, err := file.Open()
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "upload failed: cannot open file"))
		return
	}
	defer src.Close()
	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	key, err := h.minio.Upload(ctx, objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		h.logger.Error(fmt.Sprintf(constants.LogUploadFailed, file.Filename), "error", err)
		c.Error(util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, "upload failed"))
		return
	}
	h.logger.Info(fmt.Sprintf(constants.LogUploadSuccess, key), "filename", file.Filename)
	c.JSON(http.StatusOK, dto.OK(gin.H{"url": "/api/v1/files/" + key}))
}

// Get handles GET /files/:key streaming from MinIO.
func (h *UploadHandler) Get(c *gin.Context) {
	key := c.Param("key")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	reader, size, contentType, err := h.minio.Get(ctx, key)
	if err != nil {
		c.Error(util.NewAppError(http.StatusNotFound, constants.CodeNotFound, "file not found"))
		return
	}
	defer reader.Close()
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.DataFromReader(http.StatusOK, size, contentType, reader, nil)
}
