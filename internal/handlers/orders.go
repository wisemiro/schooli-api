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
		ProductID  int64  `json:"product_id"`
		VariantID  int64  `json:"variant_id"`
		Quantity   int64  `json:"quantity"`
		TotalPrice int64  `json:"total_price"`
		DeviceID   string `json:"device_id"`
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
			ProductVariant: models.ProductVariant{
				ID: req.VariantID,
			},
			Quantity:   req.Quantity,
			TotalPrice: req.TotalPrice,
			DeviceID:   req.DeviceID,
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
		ProductID  int64  `json:"product_id"`
		VariantID  int64  `json:"variant_id"`
		Quantity   int64  `json:"quantity"`
		TotalPrice int64  `json:"total_price"`
		DeviceID   string `json:"device_id"`
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
			ProductVariant: models.ProductVariant{
				ID: req.VariantID,
			},
			Quantity:   req.Quantity,
			TotalPrice: req.TotalPrice,
			DeviceID:   req.DeviceID,
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
	return func(w http.ResponseWriter, r *http.Request) {}
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
		OrderProductID int64  `json:"order_product_id"`
		GrandTotal     int64  `json:"grand_total"`
		SerialNumber   string `json:"serial_number"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			e := resterrors.NewBadRequestError(resterrors.ErrorInvalidJSONBody)
			web.Respond(r.Context(), w, r, e, e.Status)
			return
		}
		err := rp.store.CreateOrder(r.Context(), models.Order{
			OrderProduct: models.OrderProduct{
				ID: req.OrderProductID,
			},
			SerialNumber: req.SerialNumber,
			GrandTotal:   req.GrandTotal,
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
		OrderProductID int64 `json:"order_product_id"`
		GrandTotal     int64 `json:"grand_total"`
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
			ID:           int64(id),
			OrderProduct: models.OrderProduct{},
			GrandTotal:   req.GrandTotal,
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

// TODO: by user or device?
func (rp *Repository) ListOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
func (rp *Repository) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
