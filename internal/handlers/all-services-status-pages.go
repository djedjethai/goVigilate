package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/djedjethai/vigilate/internal/helpers"
	"log"
	"net/http"
)

// AllHealthyServices lists all healthy services
func (repo *DBRepo) AllHealthyServices(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "healthy", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllWarningServices lists all warning services
func (repo *DBRepo) AllWarningServices(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "warning", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllProblemServices lists all problem services
func (repo *DBRepo) AllProblemServices(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "problems", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllPendingServices lists all pending services
func (repo *DBRepo) AllPendingServices(w http.ResponseWriter, r *http.Request) {
	// get all host services(with host info) for status pending

	hostServices, err := repo.DB.GetServicesByStatus("pending")
	if err != nil {
		log.Println(err)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("services", hostServices)

	err = helpers.RenderPage(w, r, "pending", vars, nil)

	if err != nil {
		printTemplateError(w, err)
	}
}
