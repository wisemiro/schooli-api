package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"schooli-api/pkg/resterrors"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func (rp *Repository) SetupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "*", "Strict-Transport-Security"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			/* Authentication */
			r.Route("/auth", func(r chi.Router) {
				r.Post("/signup", rp.RegisterUser())
				r.Post("/signin", rp.Login())
				r.Post("/refresh", rp.RefreshToken())
			})
			/* Admin Routes
			------ /api/v1/admin/ ------
			*/
			r.Route("/admin", func(r chi.Router) {
				// Require tokens on all requests
				//
				r.Use(rp.AuthMiddleware)
				// User routes
				r.Route("/user", func(r chi.Router) {
					r.Get("/", rp.ListUsers())
					r.Get("/{id}", rp.GetUser())
					r.Put("/{id}", rp.UpdateUser())
					r.Delete("/{id}", rp.DeleteUser())
				})
				// Category routes
				r.Route("/category", func(r chi.Router) {
					r.Get("/", rp.ListCategories())
					r.Get("/{id}", rp.GetCategory())
					r.Post("/", rp.CreateCategory())
					r.Put("/{id}", rp.UpdateCategory())
					r.Delete("/{id}", rp.DeleteCategory())
				})
				// Products
				r.Route("/products", func(r chi.Router) {
					r.Post("/", rp.CreateProduct())
					r.Put("/{id}", rp.UpdateProduct())
					r.Delete("/{id}", rp.DeleteProduct())
					r.Get("/", rp.ListProducts())
					r.Get("/{id}", rp.GetProduct())
					r.Get("/by-category/{id}", rp.ListProductsByCategory())
					r.Get("/ratings/{id}", rp.ProductRating())
					r.Post("/upload-images", rp.UploadProductImages())
					r.Get("/images", rp.ListProductImages())
					r.Delete("/images", rp.DeleteProductImages())
					r.Delete("/image", rp.DeleteProductImage())
				})

				// Product Variants
				r.Route("/product-variants", func(r chi.Router) {
					r.Post("/", rp.CreateVariant())
					r.Put("/{id}", rp.UpdateVariant())
					r.Delete("/{id}", rp.DeleteVariant())
					r.Get("/product/{id}", rp.ProductVariants())
					r.Get("/{id}", rp.GetVariant())
				})
				// Order items
				r.Route("/orders-items", func(r chi.Router) {
					r.Post("/", rp.CreateOrderItem())
					r.Put("/{id}", rp.UpdateOrderItem())
					r.Delete("/{id}", rp.DeleteOrderItem())
					r.Get("/", rp.ListOrderItem())
					r.Get("/{id}", rp.GetOrderItem())
				})
				// Orders
				r.Route("/orders", func(r chi.Router) {
					r.Post("/", rp.CreateOrder())
					r.Put("/{id}", rp.UpdateOrder())
					r.Delete("/{id}", rp.DeleteOrder())
					r.Get("/", rp.ListOrders())
					// r.Get("/{id}", rp.GetOrder())
				})
				// Ratings
				r.Route("/ratings", func(r chi.Router) {
					r.Delete("/{id}", rp.DeleteRating())
					r.Get("/{id}", rp.GetRating())
				})
				// carousel
				r.Route("/carousel", func(r chi.Router) {
					r.Post("/", rp.UploadCarouselImages())
					r.Get("/", rp.ListCarouselImages())
					r.Delete("/", rp.DeleteCarousel())
					r.Delete("/image/", rp.DeleteCarouselImage())
				})
			})
			/* Client routes
			------ /api/v1/client/-------
			*/
			r.Route("/client", func(r chi.Router) {
				// User routes
				r.Route("/user", func(r chi.Router) {
					r.Use(rp.AuthMiddleware)
					r.Put("/{id}", rp.UpdateUser())
					r.Get("/{id}", rp.GetUser())
				})
				// Categories
				r.Route("/category", func(r chi.Router) {
					r.Get("/", rp.ListCategories())
				})
				// Products
				r.Route("/products", func(r chi.Router) {
					r.Get("/", rp.ListProducts())
					r.Get("/by-category/{id}/{last_id}", rp.ListProductsByCategory())
					r.Get("/{id}", rp.GetProduct())
					r.Post("/", rp.AddWishlist())
					r.Get("/ratings/{id}", rp.ProductRating())
					r.Get("/images", rp.ListProductImages())
				})
				// Wishlist
				r.Route("/wishlist", func(r chi.Router) {
					r.Use(rp.AuthMiddleware)
					r.Post("/", rp.AddWishlist())
					r.Get("/", rp.ListWishlist())

				})
				// Product Variants
				r.Route("/product-variants", func(r chi.Router) {
					r.Get("/{id}", rp.ProductVariants())
				})
				// Order items
				r.Route("/order-items", func(r chi.Router) {
					r.Post("/", rp.CreateOrderItem())
					r.Put("/{id}", rp.UpdateOrderItem())
					r.Delete("/{id}", rp.DeleteOrderItem())
					r.Get("/", rp.ListOrderItem())
					r.Get("/{id}", rp.GetOrderItem())
				})
				// Orders
				r.Route("/orders", func(r chi.Router) {
					r.Post("/", rp.CreateOrder())
					r.Put("/{id}", rp.UpdateOrder())
					r.Delete("/{id}", rp.DeleteOrder())
					r.Get("/", rp.ListOrders())
					r.Get("/{id}", rp.ListUserOrders())
				})
				// Ratings
				r.Route("/ratings", func(r chi.Router) {
					r.Use(rp.AuthMiddleware)
					r.Post("/", rp.CreateRating())
					r.Put("/{id}", rp.UpdateRating())
					r.Delete("/{id}", rp.DeleteRating())
				})
				// carousel
				r.Route("/carousel", func(r chi.Router) {
					r.Get("/", rp.ListCarouselImages())
				})
			})
		})
	})
	return router

}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func (rp *Repository) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			render.Respond(w, r, resterrors.NewError("authorization header is not provided"))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			render.Respond(w, r, resterrors.NewError("invalid authorization header format"))
			return
		}

		authorizationType := strings.ToLower(fields[0])

		if authorizationType != authorizationTypeBearer {
			err := fmt.Sprintf("unsupported authorization type %s", authorizationType)
			render.Respond(w, r, err)
			return
		}
		accessToken := fields[1]
		payload, err := rp.tokenMaker.Verify(accessToken)
		if err != nil {
			switch err {
			case ErrExpiredToken:
				w.WriteHeader(http.StatusBadRequest)
				render.Respond(w, r, resterrors.NewError("Token has expired"))
				return
			case ErrInvalidToken:
				w.WriteHeader(http.StatusBadRequest)
				render.Respond(w, r, resterrors.NewError("Token is invalid"))
				return
			default:
				w.WriteHeader(http.StatusBadRequest)
				render.Respond(w, r, err.Error())
			}
			render.Respond(w, r, resterrors.NewError("invalid authorization header format"))
			return
		}
		ctx := context.WithValue(r.Context(), "user", payload.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
