package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"schooli-api/internal/models"
	"schooli-api/pkg/resterrors"
	"schooli-api/pkg/web"

	"github.com/go-chi/render"
)

func (rp *Repository) UploadCarouselImages() http.HandlerFunc {
	type request struct {
		CaouselName string `json:"carousel_name"`
		Files       []struct {
			Filename string `json:"filename"`
		} `json:"files"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			e := resterrors.NewBadRequestError("Failed to parse json")
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		filenames := make([]*multipart.FileHeader, len(req.Files))
		for i, file := range req.Files {
			fileHeader := &multipart.FileHeader{
				Filename: file.Filename,
				Header:   textproto.MIMEHeader(http.Header{}),
			}
			filenames[i] = fileHeader
		}
		_, uploadErrs := rp.storageService.MultipleFileUpload(r.Context(), &models.MultipleFileUploadModel{
			FileNames:  filenames,
			FolderName: req.CaouselName,
		})
		if len(uploadErrs) > 0 {
			e := resterrors.NewBadRequestError("Failed to upload some files")
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}

		resp := NewStatusCreatedResponse(SuccessMessage, nil)
		render.Respond(w, r, resp)
	}
}

func (rp *Repository) ListCarouselImages() http.HandlerFunc {
	type listResponse struct {
		Total  int `json:"total"`
		Images any `json:"images"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		folderName := r.URL.Query().Get("folder_name")

		imgs, err := rp.storageService.ListFiles(r.Context(), folderName)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := listResponse{
			Total:  len(imgs),
			Images: imgs,
		}
		e := NewStatusOkResponse(SuccessMessage, data)
		web.Respond(r.Context(), w, r, e, e.Status)
	}
}

func (rp *Repository) DeleteCarousel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		folderName := r.URL.Query().Get("folder_name")
		err := rp.storageService.DeleteFolder(r.Context(), folderName)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		e := NewDeleteResponse(SuccessMessage, nil)
		web.Respond(r.Context(), w, r, e, e.Status)
	}
}

func (rp *Repository) DeleteCarouselImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Query().Get("file_path")
		err := rp.storageService.DeleteFolderFile(r.Context(), filePath)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		e := NewDeleteResponse(SuccessMessage, nil)
		web.Respond(r.Context(), w, r, e, e.Status)
	}
}
