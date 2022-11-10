package api

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/util"
)

func jsonBindingErrResp(err error, c *gin.Context, errorMap map[string]string) error {
	err1 := util.ReadJSONDecoderErr(err)
	if err1 != nil {
		errBadRequestResp(c, err1.Error())
		return nil
	}
	errMap := util.ReadValidationErr(err, errorMap)
	if errMap != nil {
		errValidationResp(c, errMap)
		return nil
	}
	return err
}

func uploadFile(c *gin.Context, uploader attachmentManager, validImageTypes []string) string {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrNotMultipart) || errors.Is(err, http.ErrMissingBoundary) {
			errBadRequestResp(c, "content-Type header is not valid")
			return ""
		}
		if errors.Is(err, http.ErrMissingFile) {
			errBadRequestResp(c, "request did not contain a file")
			return ""
		}
		if errors.Is(err, multipart.ErrMessageTooLarge) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"status":  "fail",
				"message": "max byte to upload is 32mB",
			})
			return ""
		}
		errServerResp(c, err)
		return ""
	}
	file, err := fileHeader.Open()
	if err != nil {
		errServerResp(c, err)
		return ""
	}
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		errServerResp(c, err)
		return ""
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		errServerResp(c, err)
		return ""
	}

	var isValidType bool
	contentType := http.DetectContentType(buff)
	for _, validTypes := range validImageTypes {
		if contentType == validTypes {
			isValidType = true
			break
		}
	}
	if !isValidType {
		errUnprocessableEntityResp(c, "not valid mime-type")
		return ""
	}
	key, err := uploader.UploadBlob(file, fileHeader.Size, contentType)
	if err != nil {
		errServerResp(c, err)
		return ""
	}
	return key
}
