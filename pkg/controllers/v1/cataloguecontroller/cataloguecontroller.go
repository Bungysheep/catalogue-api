package cataloguecontroller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bungysheep/catalogue-api/pkg/commons/contextkey"
	"github.com/bungysheep/catalogue-api/pkg/commons/status"
	"github.com/bungysheep/catalogue-api/pkg/controllers/v1/basecontroller"
	cataloguemodel "github.com/bungysheep/catalogue-api/pkg/models/v1/catalogue"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/signinclaimresource"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/cataloguerepository"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/customfielddefinitionrepository"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/productrepository"
	"github.com/gorilla/mux"
)

// CatalogueController type
type CatalogueController struct {
	basecontroller.BaseResource
}

// NewCatalogueController - Creates catalogue controller
func NewCatalogueController() *CatalogueController {
	return &CatalogueController{}
}

// GetAll - Return all catalogues
func (clgCtl *CatalogueController) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Retrieving all Catalogues.\n")

	clgRepo := cataloguerepository.NewCatalogueRepository()
	result, err := clgRepo.GetAll(r.Context())
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	clgCtl.WriteResponse(w, http.StatusOK, true, result, "")
}

// GetByID - Return a catalogue
func (clgCtl *CatalogueController) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	code := params["id"]

	log.Printf("Retrieving Catalogue '%v'.\n", code)

	clgRepo := cataloguerepository.NewCatalogueRepository()
	result, err := clgRepo.GetByID(r.Context(), code)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	clgCtl.WriteResponse(w, http.StatusOK, true, result, "")
}

// Create - Create new catalogue
func (clgCtl *CatalogueController) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating Catalogue.\n")

	authClaims := r.Context().Value(contextkey.ClaimToken).(signinclaimresource.SignInClaimResource)

	newClg := cataloguemodel.NewCatalogue()
	err := json.NewDecoder(r.Body).Decode(newClg)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Invalid create catalogue request.")
		return
	}

	valid, message := newClg.DoValidate()
	if !valid {
		clgCtl.WriteResponse(w, http.StatusBadRequest, false, nil, message)
		return
	}

	newClg.Status = status.Active.String()
	newClg.CreatedBy = authClaims.GetUsername()
	newClg.ModifiedBy = authClaims.GetUsername()
	newClg.Vers = 1

	clgRepo := cataloguerepository.NewCatalogueRepository()
	nbrRows, err := clgRepo.Create(r.Context(), newClg)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if nbrRows == 0 {
		clgCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Catalogue was not created.")
		return
	}

	result, err := clgRepo.GetByID(r.Context(), newClg.GetCode())
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if result != nil {
		fieldDefRepo := customfielddefinitionrepository.NewCustomFieldDefinitionRepository()
		for _, newFieldDef := range newClg.GetAllCustomFieldDefinitions() {
			newFieldDef.CatalogueCode = newClg.GetCode()
			newFieldDef.CreatedBy = newClg.GetCreatedBy()
			newFieldDef.CreatedAt = newClg.GetCreatedAt()
			newFieldDef.ModifiedBy = newClg.GetModifiedBy()
			newFieldDef.ModifiedAt = newClg.GetModifiedAt()
			newFieldDef.Vers = 1

			nbrRows, err := fieldDefRepo.Create(r.Context(), newFieldDef)
			if err != nil {
				clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
				return
			}

			if nbrRows == 0 {
				clgCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Custom Field Definition was not created.")
				return
			}
		}

		fieldDefs, err := fieldDefRepo.GetByCatalogue(r.Context(), newClg.GetCode())
		if err != nil {
			clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
			return
		}
		result.CustomFieldDefinitions = fieldDefs
	}

	clgCtl.WriteResponse(w, http.StatusAccepted, true, result, "Catalogue has been created.")
}

// Update - Update catalogue
func (clgCtl *CatalogueController) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	code := params["id"]

	log.Printf("Updating Catalogue '%v'.\n", code)

	authClaims := r.Context().Value(contextkey.ClaimToken).(signinclaimresource.SignInClaimResource)

	clgRepo := cataloguerepository.NewCatalogueRepository()
	oldClg, err := clgRepo.GetByID(r.Context(), code)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	updClg := cataloguemodel.NewCatalogue()
	err = json.NewDecoder(r.Body).Decode(updClg)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Invalid update catalogue request.")
		return
	}

	if oldClg.GetVers() != updClg.GetVers() {
		clgCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Invalid catalogue version.")
		return
	}

	oldClg.Description = updClg.GetDescription()
	oldClg.Details = updClg.GetDetails()
	oldClg.Status = updClg.GetStatus()
	oldClg.ModifiedBy = authClaims.GetUsername()

	nbrRows, err := clgRepo.Update(r.Context(), oldClg)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if nbrRows == 0 {
		clgCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Catalogue was not updated.")
		return
	}

	result, err := clgRepo.GetByID(r.Context(), oldClg.GetCode())
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	clgCtl.WriteResponse(w, http.StatusAccepted, true, result, "Catalogue has been updated.")
}

// Delete - Delete catalogue
func (clgCtl *CatalogueController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	code := params["id"]

	log.Printf("Deleting Catalogue '%v'.\n", code)

	clgRepo := cataloguerepository.NewCatalogueRepository()
	nbrRows, err := clgRepo.Delete(r.Context(), code)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if nbrRows == 0 {
		clgCtl.WriteResponse(w, http.StatusNotFound, false, nil, "Catalogue does not exist.")
		return
	}

	// Also delete all related custom field definitions
	fieldDefRepo := customfielddefinitionrepository.NewCustomFieldDefinitionRepository()
	err = fieldDefRepo.DeleteByCatalogue(r.Context(), code)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	// Also delete all related products
	prodRepo := productrepository.NewProductRepository()
	err = prodRepo.DeleteByCatalogue(r.Context(), code)
	if err != nil {
		clgCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	clgCtl.WriteResponse(w, http.StatusOK, true, nil, "Catalogue has been deleted.")
}
