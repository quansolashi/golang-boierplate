package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/internal/web/response"
	"github.com/quansolashi/golang-boierplate/backend/pkg/storage"
)

func (c *controller) uploadRoutes(rg *gin.RouterGroup) {
	rg.POST("/single", c.uploadFile)
}

// FIXME: test with S3 storage.
func (c *controller) uploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		c.badRequest(ctx, err)
		return
	}

	src, err := file.Open()
	if err != nil {
		c.httpError(ctx, fmt.Errorf("controller: could not open file %s", err.Error()))
		return
	}
	defer src.Close()

	path := fmt.Sprintf("upload/private/%s", time.Now().Format("2006-01-02-15-04"))
	err = c.bucket.Upload(ctx, path, src, storage.WithACL(types.ObjectCannedACLPrivate))
	if err != nil {
		c.httpError(ctx, fmt.Errorf("controller: could not upload file %s", err.Error()))
	}

	url, err := c.bucket.GetPresignedURL(ctx, path, 3600*time.Second)
	if err != nil {
		c.httpError(ctx, fmt.Errorf("controller: could not get file url %s", err.Error()))
	}

	res := &response.UploadFileResponse{
		FileURL: url,
	}
	ctx.JSON(http.StatusOK, res)
}
