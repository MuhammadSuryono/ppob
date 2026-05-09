package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/product-service/internal/dto"
	"github.com/yontech/ppob/product-service/internal/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	resp, err := h.productService.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetProductByCode(c *gin.Context) {
	code := c.Param("code")

	resp, err := h.productService.GetProductByCode(c.Request.Context(), code)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.productService.ListProducts(c.Request.Context(), uint(categoryID), status, limit, offset)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
	keyword := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.productService.SearchProducts(c.Request.Context(), keyword, limit, offset)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.productService.UpdateProduct(c.Request.Context(), uint(id), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), uint(id)); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *ProductHandler) ListCategories(c *gin.Context) {
	resp, err := h.productService.ListCategories(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": resp})
}

func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.productService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *ProductHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	resp, err := h.productService.GetCategory(c.Request.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.productService.UpdateCategory(c.Request.Context(), uint(id), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	if err := h.productService.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}