package api

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	apiError "github.com/xyedo/blindate/pkg/common/error"
	"github.com/xyedo/blindate/pkg/common/util"
	attachmentEntity "github.com/xyedo/blindate/pkg/domain/attachment"
)

func jsonBindingErrResp(err error, c *gin.Context, errorMap map[string]string) error {
	if err := util.ReadJSONDecoderErr(err); err != nil {
		errBadRequestResp(c, err.Error())
		return nil
	}
	if errMap := util.ReadValidationErr(err, errorMap); errMap != nil {
		errValidationResp(c, errMap)
		return nil
	}
	return err
}

func jsonHandleError(c *gin.Context, err error) {
	var apiErr apiError.API
	if errors.As(err, &apiErr) {
		status, msg := apiErr.APIError()
		c.AbortWithStatusJSON(status, gin.H{
			"status":  "fail",
			"message": msg,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "the server encountered a problem and could not process your request",
		})
	}
}
func uploadFile(c *gin.Context, uploader attachmentManager, validMimeTypes []string, prefix string) (key string, mediaType string) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrNotMultipart) || errors.Is(err, http.ErrMissingBoundary) {
			errBadRequestResp(c, "content-Type header is not valid")
			return "", ""
		}
		if errors.Is(err, http.ErrMissingFile) {
			errBadRequestResp(c, "request did not contain a file")
			return "", ""
		}
		if errors.Is(err, multipart.ErrMessageTooLarge) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"status":  "fail",
				"message": "max byte to upload is 8mB",
			})
			return "", ""
		}
		errServerResp(c, err)
		return "", ""
	}
	file, err := fileHeader.Open()
	if err != nil {
		errServerResp(c, err)
		return "", ""
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		errServerResp(c, err)
		return "", ""
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		errServerResp(c, err)
		return "", ""
	}

	var isValidType bool
	contentType := http.DetectContentType(buff)
	for _, validTypes := range validMimeTypes {
		if contentType == validTypes {
			isValidType = true
			mediaType = contentType
			break
		}
	}
	if !isValidType {
		errUnprocessableEntityResp(c, "not valid mime-type")
		return "", ""
	}
	key, err = uploader.UploadBlob(file, attachmentEntity.Uploader{
		Length:      fileHeader.Size,
		ContentType: contentType,
		Prefix:      prefix,
		Ext:         filepath.Ext(fileHeader.Filename),
	})
	if err != nil {
		errServerResp(c, err)
		return "", ""
	}
	return key, mediaType
}
