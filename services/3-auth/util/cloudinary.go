package util

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

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

func (c *Cloudinary) UploadImg(ctx context.Context, file multipart.File, filePath string) (*uploader.UploadResult, error) {
	uploadParams := uploader.UploadParams{
		PublicID:     fmt.Sprintf("jobber/auth/%s", filePath),
		Format:       "webp",
		ResourceType: "image",
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
