package util

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinary() *Cloudinary {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatalf("failed to initialize cloudinary, %v", err)
	}
	log.Println("Cloudinary connected")

	return &Cloudinary{
		cld: cld,
	}
}

func (c *Cloudinary) UploadFile(ctx context.Context, formHeader *multipart.FileHeader, file multipart.File, filePath string, fileType string) (*uploader.UploadResult, error) {
	var uploadParams uploader.UploadParams

	switch {
	case strings.HasPrefix(fileType, "image/"):
		if formHeader.Size > 2*1024*1024 {
			return nil, fmt.Errorf("Image file is larger than 2MB")
		}

		uploadParams = uploader.UploadParams{
			ResourceType: "image",
			Format:       "webp", // Convert images to WebP
		}
	case strings.HasPrefix(fileType, "video/"):
		if formHeader.Size > 10*1024*1024 {
			return nil, fmt.Errorf("Video file is larger than 10MB")
		}

		uploadParams = uploader.UploadParams{
			ResourceType:   "video",
			Format:         "mp4",           // Convert video to MP4
			Transformation: "f_auto,q_auto", // Optimize quality and size
		}
	default:
		if formHeader.Size > 2*1024*1024 {
			return nil, fmt.Errorf("Uploaded file is larger than 2MB")
		}

		uploadParams = uploader.UploadParams{
			ResourceType: "auto", // Keep original format for other files
		}
	}

	result, err := c.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		log.Println("error uploading file", err)
		return nil, err
	}

	return result, nil
}

func (c *Cloudinary) DestroyByPrefix(ctx context.Context, prefix string, filePath string) (bool, error) {
	newBool := func(b bool) *bool {
		return &b
	}

	result, err := c.cld.Admin.DeleteAssetsByPrefix(ctx, admin.DeleteAssetsByPrefixParams{
		Prefix: api.CldAPIArray{
			prefix,
		},
		Invalidate: newBool(true),
	})

	if err != nil {
		return false, err
	}

	return len(result.DeletedCounts) > 0, nil
}

func (c *Cloudinary) Destroy(ctx context.Context, publicID string) (string, error) {
	newBool := func(b bool) *bool {
		return &b
	}

	result, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:   publicID,
		Invalidate: newBool(true),
	})

	if err != nil {
		return "", err
	}

	return result.Result, nil
}
