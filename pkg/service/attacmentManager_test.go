package service

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/domain"
)

func TestUploadAndDeleteBlob(t *testing.T) {
	err := godotenv.Load("../../.env")
	require.NoError(t, err)
	s3man := NewS3()
	key := testUpload(s3man, t)
	t.Cleanup(func() {
		time.Sleep(10 * time.Second)
		base := filepath.Base(key)
		cleanKey := strings.Split(base, ".")
		err := s3man.DeleteBlob("profile-picture-resized/" + cleanKey[0] + "_640" + "." + cleanKey[1])
		require.NoError(t, err)
		err = s3man.DeleteBlob("profile-picture-resized/" + cleanKey[0] + "_160" + "." + cleanKey[1])
		require.NoError(t, err)
	})
}

func TestUploadAndPresignedIt(t *testing.T) {
	err := godotenv.Load("../../.env")
	require.NoError(t, err)
	s3man := NewS3()
	key := testUpload(s3man, t)

	time.Sleep(10 * time.Second)

	base := filepath.Base(key)
	cleanKey := strings.Split(base, ".")
	url, err := s3man.GetPresignedUrl("profile-picture-resized/" + cleanKey[0] + "_640" + "." + cleanKey[1])
	require.NoError(t, err)
	assert.NotEmpty(t, url)
	res, err := http.Get(url)
	require.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func testUpload(s3man *attachment, t *testing.T) string {
	file, err := os.Open("../../assets/test.png")
	require.NoError(t, err)
	fileInfo, err := os.Stat(file.Name())
	require.NoError(t, err)
	key, err := s3man.UploadBlob(file, domain.Uploader{
		Length:      fileInfo.Size(),
		ContentType: "image/png",
		Prefix:      "profile-picture",
		Ext:         ".png",
	})
	assert.NoError(t, err)
	return key
}
