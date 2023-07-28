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

func (rp *Repository) CreateOrderItem() http.HandlerFunc {
	type request struct {
		ProductID       int64   `json:"product_id"`
		ProductVariants []int32 `json:"product_variants"`
		Quantity        int64   `json:"quantity"`
		TotalPrice      int64   `json:"total_price"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.CreateOrderItem(r.Context(), models.OrderProduct{
			Product: models.Product{
				ID: req.ProductID,
			},
			ProductVariants: req.ProductVariants,
			Quantity:        req.Quantity,
			TotalPrice:      req.TotalPrice,
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

func (rp *Repository) UpdateOrderItem() http.HandlerFunc {
	type request struct {
		ProductID       int64   `json:"product_id"`
		ProductVariants []int32 `json:"product_variants"`
		Quantity        int64   `json:"quantity"`
		TotalPrice      int64   `json:"total_price"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.UpdateOrderItem(r.Context(), models.OrderProduct{
			ID: int64(id),
			Product: models.Product{
				ID: req.ProductID,
			},
			ProductVariants: req.ProductVariants,
			Quantity:        req.Quantity,
			TotalPrice:      req.TotalPrice,
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

func (rp *Repository) DeleteOrderItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		err := rp.store.DeleteOrderItem(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewDeleteResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) ListOrderItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := rp.store.ListOrderItems(r.Context())
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, items)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) GetOrderItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		orderProducts, err := rp.store.GetOrderProduct(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, orderProducts)
		render.Respond(w, r, data)
	}
}

// Orders
func (rp *Repository) CreateOrder() http.HandlerFunc {
	type request struct {
		GrandTotal   int64  `json:"grand_total"`
		SerialNumber string `json:"serial_number"`
		ShippingID   int64  `json:"shipping_address"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		userID := rp.GetUserPayload(w, r)
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.CreateOrder(r.Context(), models.Order{
			User: models.User{
				ID: userID,
			},
			SerialNumber: req.SerialNumber,
			GrandTotal:   req.GrandTotal,
			Shipping: models.Shipping{
				ID: req.ShippingID,
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

func (rp *Repository) UpdateOrder() http.HandlerFunc {
	type request struct {
		GrandTotal int64 `json:"grand_total"`
		Confirmed  bool  `json:"confirmed"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.UpdateOrder(r.Context(), models.Order{
			ID:         int64(id),
			GrandTotal: req.GrandTotal,
			Confirmed:  req.Confirmed,
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
func (rp *Repository) DeleteOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		err := rp.store.DeleteOrder(r.Context(), int64(id))
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewDeleteResponse(SuccessMessage, nil)
		render.Respond(w, r, data)
	}
}

func (rp *Repository) ListOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		order, err := rp.store.ListOrders(r.Context())
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, order)
		render.Respond(w, r, data)
	}
}
func (rp *Repository) ListUserOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := rp.GetUserPayload(w, r)
		order, err := rp.store.ListUserOrders(r.Context(), userID)
		if err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorProcessingRequest)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		data := NewStatusOkResponse(SuccessMessage, order)
		render.Respond(w, r, data)
	}
}
