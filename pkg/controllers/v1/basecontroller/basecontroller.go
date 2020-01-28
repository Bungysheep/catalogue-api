package basecontroller

import (
	"encoding/json"
	"net/http"
)

// BaseResource type
type BaseResource struct {
}

// NewBaseController - Creates new base controller
func NewBaseController() *BaseResource {
	return &BaseResource{}
}

// WriteResponse - Writes http response
func (resource *BaseResource) WriteResponse(w http.ResponseWriter, statusCode int, success bool, data interface{}, message string) {
	resp := map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
