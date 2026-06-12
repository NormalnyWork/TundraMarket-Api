package httptransport

import (
	"net/http"

	appproduct "tundraMarket/internal/application/product"
)

type ProductHandler struct {
	uc *appproduct.UseCase
}

func NewProductHandler(uc *appproduct.UseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) Catalog(w http.ResponseWriter, r *http.Request) {
	products, err := h.uc.GetAll(r.Context())
	if err != nil {
		writeProtoError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeAuto(w, r, http.StatusOK, appproduct.ToProtoCatalog(products))
}
