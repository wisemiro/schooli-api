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

func (rp *Repository) GetUserPayload(_ http.ResponseWriter, r *http.Request) int64 {
	ctx := r.Context()
	d := ctx.Value("user")
	return d.(int64)
}

func (rp *Repository) GetInt(w http.ResponseWriter, r *http.Request, key string) int64 {
	id := chi.URLParam(r, key)
	value, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		render.Respond(w, r, resterrors.NewBadRequestError("unable to parse id"))
		return 0
	}
	if id == "" {
		render.Respond(w, r, resterrors.NewBadRequestError("id is not provided"))
		return 0
	}
	if value == 0 {
		render.Respond(w, r, resterrors.NewBadRequestError("id provided is not valid"))
		return 0
	}

	return value
}

func (rp *Repository) UpdateUser() http.HandlerFunc {
	type request struct {
		Email       string `json:"email" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := rp.GetUserPayload(w, r)
		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.UpdateUser(r.Context(), models.User{
			ID:          id,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
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

func (rp *Repository) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		err := rp.store.DeleteUser(r.Context(), models.User{
			ID: int64(id),
		})
		if err != nil {
			e := resterrors.NewNotFoundError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		res := NewStatusOkResponse(SuccessMessage, nil)
		render.Respond(w, r, res)
	}
}

func (rp *Repository) ListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := rp.store.ListUsers(r.Context())
		if err != nil {
			e := resterrors.NewNotFoundError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		res := NewStatusOkResponse(SuccessMessage, users)
		render.Respond(w, r, res)
	}
}

func (rp *Repository) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		user, err := rp.store.OneUser(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewNotFoundError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		res := NewStatusOkResponse(SuccessMessage, user)
		render.Respond(w, r, res)
	}
}
