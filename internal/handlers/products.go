package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"schooli-api/internal/models"
	"schooli-api/pkg/resterrors"
	"schooli-api/pkg/web"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (rp *Repository) CreateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, file, err := r.FormFile("image")
		if err != nil {
			e := resterrors.NewBadRequestError("Attach image")
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		fileContentType := file.Header["Content-Type"][0]
		uploadFile, err := rp.storageService.UploadFile(r.Context(), models.FileUploadModel{
			ObjName:     file.Filename,
			FileBuf:     d,
			FileSize:    file.Size,
			ContentType: fileContentType,
			FolderName:  r.FormValue("name"),
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err = rp.store.CreateCategory(r.Context(), models.Category{
			Name:  r.FormValue("name"),
			Image: uploadFile.Key,
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		info := NewStatusCreatedResponse(SuccessMessage, nil)
		render.Respond(w, r, info)

	}
}

// TODO: Fix this, check nullable
func (rp *Repository) UpdateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		d, file, err := r.FormFile("image")
		if err != nil && err != http.ErrMissingFile {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		} else if d != nil {
			fileContentType := file.Header["Content-Type"][0]
			uploadFile, err := rp.storageService.UploadFile(r.Context(), models.FileUploadModel{
				ObjName:     file.Filename,
				FileBuf:     d,
				FileSize:    file.Size,
				ContentType: fileContentType,
				FolderName:  r.FormValue("name"),
			})
			if err != nil {
				e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
				web.Respond(r.Context(), w, r, e, e.Status)
				return
			}
			err = rp.store.UpdateCategory(r.Context(), models.Category{
				ID:    int64(id),
				Name:  r.FormValue("name"),
				Image: uploadFile.Key,
			})
			if err != nil {
				e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
				web.Respond(r.Context(), w, r, e, e.Status)
				return
			}
		}

		err = rp.store.UpdateCategory(r.Context(), models.Category{
			ID:   int64(id),
			Name: r.FormValue("name"),
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		info := NewStatusOkResponse(SuccessMessage, nil)
		render.Respond(w, r, info)
	}
}

func (rp *Repository) DeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		err := rp.store.DeleteCategory(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewDeleteResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}
func (rp *Repository) ListCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := rp.store.ListCategories(r.Context())
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, categories)
		render.Respond(w, r, data)
	}
}
func (rp *Repository) GetCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		category, err := rp.store.GetCategory(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, category)
		render.Respond(w, r, data)
	}
}

//
//
// Products

func (rp *Repository) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, file, err := r.FormFile("default_image")
		if err != nil {
			e := resterrors.NewBadRequestError("Attach image")
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		fileContentType := file.Header["Content-Type"][0]
		uploadFile, err := rp.storageService.UploadFile(r.Context(), models.FileUploadModel{
			ObjName:     file.Filename,
			FileBuf:     d,
			FileSize:    file.Size,
			ContentType: fileContentType,
			FolderName:  r.FormValue("name"),
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		price, _ := strconv.Atoi(r.FormValue("price"))
		discount_price, _ := strconv.Atoi(r.FormValue("discount_price"))
		stock_count, _ := strconv.Atoi(r.FormValue("stock_count"))
		category_id, _ := strconv.Atoi(r.FormValue("category_id"))

		err = rp.store.CreateProduct(r.Context(), models.Product{
			Name:          r.FormValue("name"),
			Price:         int64(price),
			DiscountPrice: int64(discount_price),
			Sku:           r.FormValue("sku"),
			Description:   r.FormValue("description"),
			StockCount:    int64(stock_count),
			Category: models.Category{
				ID: int64(category_id),
			},
			DefaultImage: uploadFile.Key,
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}

		data := NewStatusCreatedResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

// TODO:
func (rp *Repository) UpdateProduct() http.HandlerFunc {
	type request struct {
		Name          string `json:"name"`
		Price         int64  `json:"price"`
		DiscountPrice int64  `json:"discount_price"`
		Sku           string `json:"sku"`
		Description   string `json:"description"`
		StockCount    int64  `json:"stock_count"`
		CategoryID    int64  `json:"category_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.UpdateProduct(r.Context(), models.Product{
			ID:            int64(id),
			Name:          req.Name,
			Price:         req.Price,
			DiscountPrice: req.DiscountPrice,
			Sku:           req.Sku,
			Description:   req.Description,
			StockCount:    req.StockCount,
			Category: models.Category{
				ID: req.CategoryID,
			},
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}

		data := NewStatusOkResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		productName := r.URL.Query().Get("product_name")
		err := rp.store.DeleteProduct(r.Context(), int64(id), productName)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewDeleteResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) ListProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := rp.store.AllProducts(r.Context())
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, products)
		render.Respond(w, r, data)
	}
}
func (rp *Repository) ListDiscountedProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := rp.store.DiscountedProducts(r.Context())
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, products)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		product, err := rp.store.GetProduct(r.Context(), int64(id))
		if err != nil {
			log.Println(err)
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, product)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) ListProductsByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		lastID, _ := strconv.Atoi(chi.URLParam(r, "last_id"))

		products, err := rp.store.ByCategory(r.Context(), int64(id), int64(lastID))
		if err != nil {
			fmt.Println(err)
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, products)
		render.Respond(w, r, data)
	}
}

// Variants
func (rp *Repository) CreateVariant() http.HandlerFunc {
	type request struct {
		Name      string `json:"name"`
		ProductID int64  `json:"product_id"`
		Type      int64  `json:"type"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.CreateVariant(r.Context(), models.ProductVariant{
			Name:      req.Name,
			ProductID: req.ProductID,
			Type:      int(req.Type),
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusCreatedResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) UpdateVariant() http.HandlerFunc {
	type request struct {
		Name      string `json:"name"`
		ProductID int64  `json:"product_id"`
		Type      int64  `json:"type"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.UpdateVariant(r.Context(), models.ProductVariant{
			ID:        int64(id),
			Name:      req.Name,
			ProductID: req.ProductID,
			Type:      int(req.Type),
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) DeleteVariant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		err := rp.store.DeleteVariant(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewDeleteResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) ProductVariants() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		variants, err := rp.store.ProductVariaties(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewDeleteResponse(SuccessMessage, variants)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) GetVariant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		variants, err := rp.store.GetVariant(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, variants)
		render.Respond(w, r, data)
	}
}

// Wishlist
func (rp *Repository) AddWishlist() http.HandlerFunc {
	type request struct {
		ProductID int64 `json:"product_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		user_id := rp.GetUserPayload(w, r)
		err := rp.store.AddWishlist(r.Context(), models.Wishlist{
			Product: models.Product{
				ID: req.ProductID,
			},
			User: models.User{
				ID: user_id,
			},
		})
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) ListWishlist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id := rp.GetUserPayload(w, r)
		wishes, err := rp.store.ListWishlists(r.Context(), user_id)

		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, wishes)
		render.Respond(w, r, data)
	}
}

// File upload
func (rp *Repository) UploadProductImages() http.HandlerFunc {
	type request struct {
		ProductName string `json:"prodct_name"`
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
			FolderName: req.ProductName,
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

func (rp *Repository) ListProductImages() http.HandlerFunc {
	type listResponse struct {
		Total  int `json:"total"`
		Images any `json:"images"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		productName := r.URL.Query().Get("product_name")

		imgs, err := rp.storageService.ListFiles(r.Context(), productName)
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

func (rp *Repository) DeleteProductImages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productName := r.URL.Query().Get("product_name")
		err := rp.storageService.DeleteFolder(r.Context(), productName)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		e := NewDeleteResponse(SuccessMessage, nil)
		web.Respond(r.Context(), w, r, e, e.Status)
	}
}

func (rp *Repository) DeleteProductImage() http.HandlerFunc {
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

func (rp *Repository) SearchProducts() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}
	type response struct {
		Total    int `json:"total"`
		Products any `json:"products"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			e := resterrors.NewBadRequestError("Failed to parse json")
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		products, err := rp.store.SearchProducts(r.Context(), req.Name)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}

		data := response{
			Total:    len(products),
			Products: products,
		}
		e := NewStatusOkResponse(SuccessMessage, data)
		web.Respond(r.Context(), w, r, e, e.Status)
	}
}
