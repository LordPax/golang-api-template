package services

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

func UploadFile(file any) (string, error) {
	ctx := context.Background()
	name := uuid.New().String()

	cld, _ := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)

	param := uploader.UploadParams{PublicID: name}
	resp, err := cld.Upload.Upload(ctx, file, param)
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
