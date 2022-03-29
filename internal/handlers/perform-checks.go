package handlers

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

func (m *DBRepo) TestCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	oldStatus := chi.URLParam(r, "oldStatus")

	log.Println(hostServiceID, oldStatus)

	//
}
