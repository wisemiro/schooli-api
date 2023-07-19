package handlers

import (
	"errors"
	"net/http"
	"schooli-api/internal/models"
	"schooli-api/pkg/resterrors"
	"schooli-api/pkg/web"

	"github.com/go-chi/render"
)

func (rp *Repository) RegisterUser() http.HandlerFunc {
	type req struct {
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var reQ req
		if err := render.DecodeJSON(r.Body, &reQ); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		user, err := models.NewUser(reQ.Email, reQ.Password, reQ.PhoneNumber)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data, saveErr := rp.store.CreateUser(r.Context(), *user)
		if saveErr != nil {
			if errors.Is(saveErr, errors.New(resterrors.ErrorPhoneExists)) {
				e := resterrors.NewBadRequestError(resterrors.ErrorPhoneExists)
				web.Respond(r.Context(), w, r, e, e.Status)
				return
			}
			e := resterrors.NewBadRequestError(resterrors.ErrorEmailExists)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		deta := NewStatusOkResponse(SuccessMessage, data)
		render.Respond(w, r, deta)
	}
}

func (rp *Repository) Login() http.HandlerFunc {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		User         any    `json:"user"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var reQ req
		if err := render.DecodeJSON(r.Body, &reQ); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		user, err := rp.store.GetByEmail(r.Context(), reQ.Email)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorAccountNotFound)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}

		ok := user.CheckPasswordHash(reQ.Password)
		if !ok {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidEmailPassword)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}

		accessToken, err := rp.tokenMaker.Create(user.Email, user.ID)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}

		refreshToken, err := rp.tokenMaker.CreateRefresh(user.Email, user.ID)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := resp{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		}
		deta := NewStatusOkResponse(SuccessMessage, data)
		render.Respond(w, r, deta)
	}
}

func (rp *Repository) RefreshToken() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		User         any    `json:"user"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}

		payload, err := rp.tokenMaker.Verify(req.RefreshToken)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidToken)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		// Check if a user exists
		user, err := rp.store.GetByEmail(r.Context(), payload.Email)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorAccountNotFound)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		accessToken, err := rp.tokenMaker.Create(user.Email, user.ID)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		refreshToken, err := rp.tokenMaker.CreateRefresh(user.Email, user.ID)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := response{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		}

		e := NewStatusOkResponse(SuccessMessage, data)
		web.Respond(r.Context(), w, r, e, e.Status)
	}
}
