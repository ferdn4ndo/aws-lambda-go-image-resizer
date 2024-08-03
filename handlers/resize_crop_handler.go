package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"

	"github.com/ferdn4ndo/aws-lambda-go-image-resizer/model"
)

type ResizeCropHandler struct {
	initialized bool
	s3Handler   S3Bucket
}

var (
	bucket, originalFolder, resizedFolder, regional string
)

func (handler *ResizeCropHandler) init() error {
	if !handler.initialized {
		if err := handler.getConfig(); err != nil {
			return err
		}
		handler.s3Handler = new(S3Handler)
		handler.initialized = true
	}
	return nil
}

func (handler *ResizeCropHandler) getConfig() error {
	bucket = os.Getenv("bucket")
	originalFolder = os.Getenv("original_folder")
	resizedFolder = os.Getenv("resized_folder")
	regional = os.Getenv("regional")

	if bucket == "" || originalFolder == "" || resizedFolder == "" {
		fmt.Printf("Config: %v | %v | %v | %v\n", bucket, originalFolder, resizedFolder, regional)

		return errors.New("Couldn't read config from environment!")
	}

	return nil
}

func (handler *ResizeCropHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := handler.init(); err != nil {
		fmt.Printf("Error here: %v\n", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	vars := mux.Vars(r)

	// Bind params to model
	img := &model.Image{
		Optional: vars["optional"],
	}

	// Validate model state
	if !img.IsMatchFormat() {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	fmt.Printf("Img: %v | %v | %v | %v | %v | %v\n", img.FileName, img.Crop, img.Extension, img.Dimension, img.Height, img.Width)
	ctx := request.Context()
	currentSession := session.Must(session.NewSession())
	currentSession.Config.Region = aws.String(regional)

	originalKey := img.GetS3Key(originalFolder, img.FileName)
	exist, data, err := handler.s3Handler.DownloadImage(ctx, currentSession, bucket, originalKey)
	if !exist {
		fmt.Printf("Not found image with key: %v\n", originalKey)
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	if err != nil {
		fmt.Printf("Download image error : %v\n", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	// Decode from downloaded data
	originalImage, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Printf("Decode image error: %v\n", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	// Resize image
	resizedImage := img.ResizeOrCrop(originalImage)
	if resizedImage == nil {
		fmt.Printf("Resized image error: %v\n", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	// Encode resize image and upload to s3
	var bufferEncode = new(bytes.Buffer)
	errEncode := imaging.Encode(bufferEncode, resizedImage, model.ParseExtension(model.ParseContentType(img.Extension)))
	if errEncode != nil {
		fmt.Printf("Encode image error: %v\n", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	// Upload to S3
	targetKey := img.GetS3Key(resizedFolder, img.GetOutputFileName())
	output, err := handler.s3Handler.UploadImage(ctx, currentSession, bucket, targetKey, bufferEncode.Bytes())

	if err != nil {
		fmt.Printf("Upload image error: %v\n", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	// Serve file just uploaded
	fmt.Printf("New image has uploaded at: %v\n", output.Location)
	http.Redirect(writer, request, output.Location, http.StatusTemporaryRedirect)
}
