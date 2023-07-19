package services

import (
	"context"
	"encoding/json"
	"schooli-api/internal/models"
	"schooli-api/internal/repository/postgresql/db"
	"schooli-api/pkg/resterrors"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProductsService interface {
	CreateProduct(ctx context.Context, pm models.Product) error
	UpdateProduct(ctx context.Context, pm models.Product) error
	DeleteProduct(ctx context.Context, productID int64) error
	AllProducts(ctx context.Context) ([]*models.Product, error)
	ByCategory(ctx context.Context, categoryID int64, lastID int64) ([]*models.Product, error)
	GetProduct(ctx context.Context, productID int64) (*models.Product, error)
	//
	//
	// Wishlist
	AddWishlist(ctx context.Context, wm models.Wishlist) error
	ListWishlists(ctx context.Context, userID int64) ([]*models.Wishlist, error)
	//
	//
	// Categories
	CreateCategory(ctx context.Context, cm models.Category) error
	UpdateCategory(ctx context.Context, cm models.Category) error
	DeleteCategory(ctx context.Context, categoryID int64) error
	ListCategories(ctx context.Context) ([]*models.Category, error)
	GetCategory(ctx context.Context, categoryID int64) (*models.Category, error)
	//
	// Product variants
	CreateVariant(ctx context.Context, vm models.ProductVariant) error
	UpdateVariant(ctx context.Context, vm models.ProductVariant) error
	DeleteVariant(ctx context.Context, variantID int64) error
	ProductVariaties(ctx context.Context, productID int64) ([]*models.ProductVariant, error)
	GetVariant(ctx context.Context, varietyID int64) (*models.ProductVariant, error)
}

func (sq *SQLStore) CreateProduct(ctx context.Context, pm models.Product) error {
	if err := sq.store.CreateProduct(ctx, db.CreateProductParams{
		Name:          pm.Name,
		Price:         int32(pm.Price),
		DiscountPrice: pgtype.Int4{Int32: int32(pm.DiscountPrice), Valid: true},
		Sku:           pm.Sku,
		Description:   pm.Description,
		StockCount:    int32(pm.StockCount),
		CategoryID:    pm.Category.ID,
		DefaultImage:  pm.DefaultImage,
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.Create")

	}
	return nil
}
func (sq *SQLStore) UpdateProduct(ctx context.Context, pm models.Product) error {
	if err := sq.store.UpdateProduct(ctx, db.UpdateProductParams{
		Name:          pm.Name,
		Price:         int32(pm.Price),
		DiscountPrice: pgtype.Int4{Int32: int32(pm.DiscountPrice), Valid: true},
		Sku:           pm.Sku,
		Description:   pm.Description,
		StockCount:    int32(pm.StockCount),
		CategoryID:    pm.Category.ID,
		DefaultImage:  pm.DefaultImage,
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.Update")
	}
	return nil
}

func (sq *SQLStore) DeleteProduct(ctx context.Context, productID int64) error {
	if err := sq.store.DeleteProduct(ctx, productID); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.Delete")
	}
	return nil
}

func (sq *SQLStore) AllProducts(ctx context.Context) ([]*models.Product, error) {
	products, err := sq.store.ListProducts(ctx)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.List")
	}
	prod := make([]*models.Product, len(products))
	for i, v := range products {
		prod[i] = &models.Product{
			ID:            v.ID,
			CreatedAt:     v.CreatedAt.Time,
			UpdatedAt:     v.UpdatedAt.Time,
			Name:          v.Name,
			Price:         int64(v.Price),
			DiscountPrice: int64(v.DiscountPrice.Int32),
			Sku:           v.Sku,
			Description:   v.Description,
			StockCount:    int64(v.StockCount),
			DefaultImage:  sq.fileStore.BuildFilePath(v.DefaultImage),
			AverageRating: int(v.AverageRating.Int32),
		}
	}
	return prod, nil

}

func (sq *SQLStore) ByCategory(ctx context.Context, categoryID int64, lastID int64) ([]*models.Product, error) {
	products, err := sq.store.ProductsByCategory(ctx, db.ProductsByCategoryParams{
		CategoryID: categoryID,
		ID:         lastID,
	})
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.ListByCategory")
	}
	prod := make([]*models.Product, len(products))
	for i, v := range products {
		prod[i] = &models.Product{
			ID:            v.ID,
			CreatedAt:     v.CreatedAt.Time,
			UpdatedAt:     v.UpdatedAt.Time,
			Name:          v.Name,
			Price:         int64(v.Price),
			DiscountPrice: int64(v.DiscountPrice.Int32),
			Sku:           v.Sku,
			Description:   v.Description,
			StockCount:    int64(v.StockCount),
			DefaultImage:  sq.fileStore.BuildFilePath(v.DefaultImage),
			AverageRating: int(v.AverageRating.Int32),
			// AvarageRating: int(v.AverageRating.Int32),
		}
	}
	return prod, nil
}
func (sq *SQLStore) GetProduct(ctx context.Context, productID int64) (*models.Product, error) {
	v, err := sq.store.OneProduct(ctx, productID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.OneProduct")
	}

	// Handle empty product specifications
	var productSpecs []models.ProductSpecification
	if v.ProductSpecifications != nil {
		productSpecsBytes, ok := v.ProductSpecifications.([]byte)
		if !ok {
			productSpecsBytes, err = json.Marshal(v.ProductSpecifications)
			if err != nil {
				return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.OneProduct.MarshalSpecs")
			}
		}

		err = json.Unmarshal(productSpecsBytes, &productSpecs)
		if err != nil {
			return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.OneProduct.UnmarshalSpecs")
		}
	}

	// Handle empty product variants
	var productVariants []models.ProductVariant
	if v.ProductVariants != nil {
		productVariantsBytes, ok := v.ProductVariants.([]byte)
		if !ok {
			productVariantsBytes, err = json.Marshal(v.ProductVariants)
			if err != nil {
				return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.OneProduct.MarshalVariants")
			}
		}

		err = json.Unmarshal(productVariantsBytes, &productVariants)
		if err != nil {
			return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.OneProduct.UnmarshalVariants")
		}
	}

	return &models.Product{
		ID:                   v.ID,
		CreatedAt:            v.CreatedAt.Time,
		UpdatedAt:            v.UpdatedAt.Time,
		Name:                 v.Name,
		Price:                int64(v.Price),
		DiscountPrice:        int64(v.DiscountPrice.Int32),
		Sku:                  v.Sku,
		Description:          v.Description,
		StockCount:           int64(v.StockCount),
		DefaultImage:         sq.fileStore.BuildFilePath(v.DefaultImage),
		AverageRating:        int(v.AverageRating.Int32),
		TotalRatings:         int(v.TotalRatings),
		TotalViews:           int(v.TotalView.Int32),
		MinStockCount:        int64(v.MinStockCount),
		ProductVariant:       productVariants,
		ProductSpecification: productSpecs,
	}, nil
}

// Wishlist
func (sq *SQLStore) AddWishlist(ctx context.Context, wm models.Wishlist) error {
	wishlist, _ := sq.store.GetWishlist(ctx, db.GetWishlistParams{
		UserID:    wm.ID,
		ProductID: wm.Product.ID,
	})
	if wishlist.ProductID == 0 {
		err := sq.store.CreateWishlist(ctx, db.CreateWishlistParams{
			ProductID: wm.Product.ID,
			UserID:    wm.User.ID,
		})
		if err != nil {
			return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.CreateWishlist")
		}
		return nil
	}
	err := sq.store.DeleteWishlist(ctx, wishlist.ID)
	if err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.DeleteWishlist")
	}
	return nil
}

func (sq *SQLStore) ListWishlists(ctx context.Context, userID int64) ([]*models.Wishlist, error) {
	wish, err := sq.store.ListWishlist(ctx, userID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.ListWishlist")
	}
	wishlist := make([]*models.Wishlist, len(wish))

	for i, v := range wish {
		wishlist[i] = &models.Wishlist{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Time,
			Product: models.Product{
				ID:            v.ID_2.Int64,
				Name:          v.Name.String,
				Price:         int64(v.Price.Int32),
				DiscountPrice: int64(v.DiscountPrice.Int32),
				DefaultImage:  sq.fileStore.BuildFilePath(v.DefaultImage.String),
			},
		}
	}
	return wishlist, nil
}

// Category
func (sq *SQLStore) CreateCategory(ctx context.Context, cm models.Category) error {
	if err := sq.store.CreateCategory(ctx, db.CreateCategoryParams{
		Name:  cm.Name,
		Image: cm.Image,
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.CreateCategory")
	}
	return nil
}
func (sq *SQLStore) UpdateCategory(ctx context.Context, cm models.Category) error {
	if err := sq.store.UpdateCategory(ctx, db.UpdateCategoryParams{
		ID:    cm.ID,
		Name:  cm.Name,
		Image: cm.Image,
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.UpdateCategory")
	}
	return nil
}

func (sq *SQLStore) DeleteCategory(ctx context.Context, categoryID int64) error {
	if err := sq.store.DeleteCategory(ctx, categoryID); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.DeleteCategory")
	}
	return nil
}

func (sq *SQLStore) ListCategories(ctx context.Context) ([]*models.Category, error) {
	categories, err := sq.store.ListCategories(ctx)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.UpdateCategory")
	}
	cat := make([]*models.Category, len(categories))
	for i, v := range categories {
		cat[i] = &models.Category{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Time,
			UpdatedAt: v.UpdatedAt.Time,
			Name:      v.Name,
			Image:     sq.fileStore.BuildFilePath(v.Image),
		}
	}
	return cat, nil
}
func (sq *SQLStore) GetCategory(ctx context.Context, categoryID int64) (*models.Category, error) {
	category, err := sq.store.GetCategory(ctx, categoryID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.UpdateCategory")
	}

	return &models.Category{
		ID:        category.ID,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
		Name:      category.Name,
		Image:     sq.fileStore.BuildFilePath(category.Image),
	}, nil
}

//
//
// Variants

func (sq *SQLStore) CreateVariant(ctx context.Context, vm models.ProductVariant) error {
	if err := sq.store.CreateVariant(ctx, db.CreateVariantParams{
		Name:      vm.Name,
		ProductID: vm.ProductID,
		Type:      int32(vm.Type),
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.CreateVariant")
	}
	return nil
}

func (sq *SQLStore) UpdateVariant(ctx context.Context, vm models.ProductVariant) error {
	if err := sq.store.UpdateVariant(ctx, db.UpdateVariantParams{
		ID:        vm.ID,
		Name:      vm.Name,
		ProductID: vm.ProductID,
		Type:      int32(vm.Type),
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.UpdateVariant")
	}
	return nil
}

func (sq *SQLStore) DeleteVariant(ctx context.Context, variantID int64) error {
	if err := sq.store.DeleteVariant(ctx, variantID); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.DeleteVariant")
	}
	return nil
}

func (sq *SQLStore) ProductVariaties(ctx context.Context, productID int64) ([]*models.ProductVariant, error) {
	variants, err := sq.store.ListVariants(ctx, productID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.ProductVariaties")
	}
	vars := make([]*models.ProductVariant, len(variants))
	for i, v := range variants {
		vars[i] = &models.ProductVariant{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Time,
			Name:      v.Name,
			ProductID: v.ProductID,
		}
	}
	return vars, nil
}

func (sq *SQLStore) GetVariant(ctx context.Context, varietyID int64) (*models.ProductVariant, error) {
	variant, err := sq.store.GetVariant(ctx, varietyID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "ProductService.GettVariaty")
	}
	return &models.ProductVariant{
		ID:        variant.ID,
		CreatedAt: variant.CreatedAt.Time,
		Name:      variant.Name,
		ProductID: variant.ProductID,
	}, nil
}
