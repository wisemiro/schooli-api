package handlers

import (
	"net/http"
	"schooli-api/internal/models"
	"schooli-api/pkg/resterrors"
	"schooli-api/pkg/web"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (rp *Repository) CreateRating() http.HandlerFunc {
	type request struct {
		ProductID int64  `json:"product_id"`
		Stars     int    `json:"stars"`
		Feeedback string `json:"feeedback"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		userID := rp.GetUserPayload(w, r)
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.CreateRatings(r.Context(), models.ProductRatings{
			User: models.User{
				ID: userID,
			},
			Product: models.Product{
				ID: req.ProductID,
			},
			Stars:     req.Stars,
			Feeedback: req.Feeedback,
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

func (rp *Repository) UpdateRating() http.HandlerFunc {
	type request struct {
		Stars     int    `json:"stars"`
		Feeedback string `json:"feeedback"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.UpdateRatings(r.Context(), models.ProductRatings{
			ID:        int64(id),
			Stars:     req.Stars,
			Feeedback: req.Feeedback,
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
func (rp *Repository) DeleteRating() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		err := rp.store.DeleteRating(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewDeleteResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) ProductRating() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		ratings, err := rp.store.ProductRatings(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, ratings)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) GetRating() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		rating, err := rp.store.GetRating(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, rating)
		render.Respond(w, r, data)
	}
}
