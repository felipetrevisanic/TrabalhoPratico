package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"rest/internal/product"
)

type Handler struct {
	repository product.Repository
}

func NewHandler(repository product.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/Product", h.handleProduct)
	mux.HandleFunc("/Product/all", h.handleProductAll)
	mux.HandleFunc("/Product/", h.handleProductByID)
	return mux
}

func (h *Handler) handleProduct(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		h.getProductByID(response, request)
	case http.MethodPost:
		h.createProduct(response, request)
	default:
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleProductAll(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		h.getAllProducts(response)
	case http.MethodDelete:
		h.deleteAllProducts(response)
	default:
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleProductByID(response http.ResponseWriter, request *http.Request) {
	id, ok := parseIDFromPath(request.URL.Path)
	if !ok {
		http.NotFound(response, request)
		return
	}

	switch request.Method {
	case http.MethodPut:
		h.updateProduct(response, request, id)
	case http.MethodDelete:
		h.deleteProduct(response, id)
	default:
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getProductByID(response http.ResponseWriter, request *http.Request) {
	idValue := request.URL.Query().Get("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		http.Error(response, "invalid id", http.StatusBadRequest)
		return
	}

	item, err := h.repository.GetByID(id)
	if err != nil {
		log.Printf("get product by id: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}

	// Preserve the C# behavior: return a placeholder when the record does not exist.
	if item == nil {
		item = &product.Product{
			ID:            id,
			Name:          "Product " + idValue,
			Description:   "Product not found in sample list",
			Price:         0,
			StockQuantity: 0,
			CreatedAt:     time.Now().UTC(),
		}
	}

	writeJSON(response, http.StatusOK, item)
}

func (h *Handler) getAllProducts(response http.ResponseWriter) {
	items, err := h.repository.GetAll()
	if err != nil {
		log.Printf("get all products: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(response, http.StatusOK, items)
}

func (h *Handler) createProduct(response http.ResponseWriter, request *http.Request) {
	var payload product.CreateRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		http.Error(response, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.repository.Create(payload)
	if err != nil {
		log.Printf("create product: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(response, http.StatusOK, item)
}

func (h *Handler) updateProduct(response http.ResponseWriter, request *http.Request, id int) {
	var payload product.UpdateRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		http.Error(response, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.repository.Update(id, payload)
	if err != nil {
		log.Printf("update product: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(response, http.StatusOK, item)
}

func (h *Handler) deleteProduct(response http.ResponseWriter, id int) {
	deleted, err := h.repository.Delete(id)
	if err != nil {
		log.Printf("delete product: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}

	if !deleted {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteAllProducts(response http.ResponseWriter) {
	deleted, err := h.repository.DeleteAll()
	if err != nil {
		log.Printf("delete all products: %v", err)
		http.Error(response, "internal server error", http.StatusInternalServerError)
		return
	}

	if !deleted {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func parseIDFromPath(path string) (int, bool) {
	idText := strings.TrimPrefix(path, "/Product/")
	if idText == "" || strings.Contains(idText, "/") {
		return 0, false
	}

	id, err := strconv.Atoi(idText)
	if err != nil {
		return 0, false
	}

	return id, true
}

func writeJSON(response http.ResponseWriter, statusCode int, payload any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	if err := json.NewEncoder(response).Encode(payload); err != nil {
		log.Printf("write json response: %v", err)
	}
}
