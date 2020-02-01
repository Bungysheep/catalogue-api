package productcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bungysheep/catalogue-api/pkg/commons/contextkey"
	"github.com/bungysheep/catalogue-api/pkg/controllers/v1/basecontroller"
	productmodel "github.com/bungysheep/catalogue-api/pkg/models/v1/product"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/signinclaimresource"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/productrepository"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/unitofmeasurerepository"
	"github.com/gorilla/mux"
)

// ProductController type
type ProductController struct {
	basecontroller.BaseResource
}

// NewProductController - Creates product controller
func NewProductController() *ProductController {
	return &ProductController{}
}

// GetByCatalogue - Return produts by catalogue
func (prodCtl *ProductController) GetByCatalogue(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	clgCode := params["clg_code"]

	log.Printf("Retrieving Products by Catalogue '%v'.\n", clgCode)

	prodRepo := productrepository.NewProductRepository()
	result, err := prodRepo.GetByCatalogue(r.Context(), clgCode)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	prodCtl.WriteResponse(w, http.StatusOK, true, result, "")
}

// GetByID - Return a product
func (prodCtl *ProductController) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	log.Printf("Retrieving Product '%v'.\n", id)

	prodRepo := productrepository.NewProductRepository()
	result, err := prodRepo.GetByID(r.Context(), id)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	prodCtl.WriteResponse(w, http.StatusOK, true, result, "")
}

// Create - Create new product
func (prodCtl *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating Product.\n")

	authClaims := r.Context().Value(contextkey.ClaimToken).(signinclaimresource.SignInClaimResource)

	newProd := productmodel.NewProduct()
	err := json.NewDecoder(r.Body).Decode(newProd)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Invalid create product request.")
		return
	}

	valid, message := newProd.DoValidate()
	if !valid {
		prodCtl.WriteResponse(w, http.StatusBadRequest, false, nil, message)
		return
	}

	newProd.Status = "A"
	newProd.CreatedBy = authClaims.GetUsername()
	newProd.ModifiedBy = authClaims.GetUsername()
	newProd.Vers = 1

	prodRepo := productrepository.NewProductRepository()
	lastID, err := prodRepo.Create(r.Context(), newProd)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if lastID == 0 {
		prodCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Product was not created.")
		return
	}

	result, err := prodRepo.GetByID(r.Context(), lastID)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if result != nil {
		uomRepo := unitofmeasurerepository.NewUnitOfMeasureRepository()
		for _, newUom := range newProd.GetAllUoms() {
			newUom.ProdID = lastID
			newUom.Vers = 1

			lastUomID, err := uomRepo.Create(r.Context(), newUom)
			if err != nil {
				prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
				return
			}

			if lastUomID == 0 {
				prodCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Unit of Measure was not created.")
				return
			}
		}

		uoms, err := uomRepo.GetByProduct(r.Context(), lastID)
		if err != nil {
			prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
			return
		}
		result.UnitOfMeasures = uoms
	}

	prodCtl.WriteResponse(w, http.StatusAccepted, true, result, "Product has been created.")
}

// Update - Update product
func (prodCtl *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	log.Printf("Updating Product '%v'.\n", id)

	authClaims := r.Context().Value(contextkey.ClaimToken).(signinclaimresource.SignInClaimResource)

	prodRepo := productrepository.NewProductRepository()
	oldProd, err := prodRepo.GetByID(r.Context(), id)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	updProd := productmodel.NewProduct()
	err = json.NewDecoder(r.Body).Decode(updProd)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Invalid update product request.")
		return
	}

	if oldProd.GetVers() != updProd.GetVers() {
		prodCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Invalid product version.")
		return
	}

	oldProd.Description = updProd.GetDescription()
	oldProd.Details = updProd.GetDetails()
	oldProd.Status = updProd.GetStatus()
	oldProd.ModifiedBy = authClaims.GetUsername()

	nbrRows, err := prodRepo.Update(r.Context(), oldProd)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if nbrRows == 0 {
		prodCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Product was not updated.")
		return
	}

	result, err := prodRepo.GetByID(r.Context(), oldProd.GetID())
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if result != nil {
		uomRepo := unitofmeasurerepository.NewUnitOfMeasureRepository()
		for _, updUom := range updProd.GetAllUoms() {
			oldUom := oldProd.GetUom(updUom.GetID())

			if oldUom != nil {
				oldUom.Code = updUom.GetCode()
				oldUom.Description = updUom.GetDescription()
				oldUom.Ratio = updUom.GetRatio()

				nbrRow, err := uomRepo.Update(r.Context(), oldUom)
				if err != nil {
					prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
					return
				}

				if nbrRow == 0 {
					prodCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Unit of Measure was not created.")
					return
				}
			}
		}

		uoms, err := uomRepo.GetByProduct(r.Context(), oldProd.GetID())
		if err != nil {
			prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
			return
		}
		result.UnitOfMeasures = uoms
	}

	prodCtl.WriteResponse(w, http.StatusAccepted, true, result, "Product has been updated.")
}

// Delete - Delete product
func (prodCtl *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)

	log.Printf("Deleting Product '%v'.\n", id)

	prodRepo := productrepository.NewProductRepository()
	nbrRows, err := prodRepo.Delete(r.Context(), id)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if nbrRows == 0 {
		prodCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Product does not exist.")
		return
	}

	// Also delete all related unit of measures
	uomRepo := unitofmeasurerepository.NewUnitOfMeasureRepository()
	err = uomRepo.DeleteByProduct(r.Context(), id)
	if err != nil {
		prodCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	prodCtl.WriteResponse(w, http.StatusOK, true, nil, "Product has been deleted.")
}
