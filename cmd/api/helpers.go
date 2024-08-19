package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/julienschmidt/httprouter"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// writeJSON writes the data into the response
//
// Parameters
//   - w: http.ResponseWriter
//   - status: http status code
//   - data : data enveloped using envelope struct
//   - headers: http headers to send
//
// Returns error if any problem occurs while writing data to json
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	js, err := json.MarshalIndent(data, "", "\n")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}

// readJSON stores the received json into the dst struct
//
// Parameters
//   - w: http.ResponseWriter
//   - r: http.Request
//   - dst: destination struct where each attribute has corresponding json tag
//
// Returns error if any problem occurs while reading data from response object
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)

	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json:unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		fn()
	}()
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("id must be an integer")
	}
	return id, nil
}

func (app *application) uploadToCloudinary(r *http.Request, file multipart.File, ext string) (string, error) {

	tempFile, err := os.CreateTemp("", "upload-*"+ext)
	if err != nil {
		return "", err
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			fmt.Printf("error closing temp file: %s\n", err)
		}
	}(tempFile)
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", err
	}
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return "", err
	}
	uploadResult, err := cld.Upload.Upload(r.Context(), tempFile.Name(), uploader.UploadParams{
		Format: "jpeg",
	})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func (app *application) getImageURL(r *http.Request, key string) (string, error) {

	file, handler, err := r.FormFile(key)
	if err != nil {
		return "", errors.New("error Retrieving File from the form")
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
	}(file)
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
	}
	if !allowedExtensions[ext] {
		return "", data.ErrUnsupportedFileType
	}
	imgURL, err := app.uploadToCloudinary(r, file, ext)
	if err != nil {
		return "", errors.New("error Uploading Image")
	}
	return imgURL, nil
}
