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

func (rp *Repository) CreateShippingAddress() http.HandlerFunc {
	type request struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitde"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		userID := rp.GetUserPayload(w, r)
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.CreateShipping(r.Context(), models.Shipping{
			User: models.User{
				ID: userID,
			},
			Geo: models.Geo{
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
			},
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

func (rp *Repository) UpdateShippingAddress() http.HandlerFunc {
	type request struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitde"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.UpdateShippingAddress(r.Context(), models.Shipping{
			ID: int64(id),
			Geo: models.Geo{
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
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

func (rp *Repository) ListShippingAddresses() http.HandlerFunc {
	type response struct {
		Total    int `json:"total"`
		Shipping any `json:"shipping"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		shipping, err := rp.store.ListShipping(r.Context())
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		sh := response{
			Total:    len(shipping),
			Shipping: shipping,
		}
		data := NewStatusOkResponse(SuccessMessage, sh)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) UserShippingAddresses() http.HandlerFunc {
	type response struct {
		Total    int `json:"total"`
		Shipping any `json:"shipping"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userID := rp.GetUserPayload(w, r)

		shipping, err := rp.store.ListUserShipping(r.Context(), userID)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		sh := response{
			Total:    len(shipping),
			Shipping: shipping,
		}
		data := NewStatusOkResponse(SuccessMessage, sh)
		render.Respond(w, r, data)
	}
}
