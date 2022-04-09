package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/djedjethai/vigilate/internal/models"
	"github.com/go-chi/chi/v5"
)

const (
	STARTIOTA = iota
	HTTP
	HTTPS
	SSLCertificate
)

type jsonResp struct {
	OK            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	LastCheck     time.Time `json:"last_check"`
}

// ScheduledCheck perform a schedule check on a host service by id
func (repo *DBRepo) ScheduledCheck(hostServiceID int) {
	log.Println("********* Running check for: ", hostServiceID)
	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		return
	}

	// testServiceForHost() is a few lines below
	msg, newStatus := repo.testServiceForHost(h, hs)

	hostServiceStatusChanged := false
	if newStatus != hs.Status {
		hostServiceStatusChanged = true
	}

	// if the host service has changed, broadcast to all clients
	if hostServiceStatusChanged {
		data := make(map[string]string)
		data["message"] = fmt.Sprintf("host service %s on %s has changed to %s", hs.Service.ServiceName, h.HostName, newStatus)
		repo.broadcastMessage("public-channel", "host-service-status-changed", data)

		// if appropriate, send email or sms message
	}

	// update host service record in db with status(if changed)
	// and update the last check
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		return
	}

	// update the service count on the front-end
	// we got the datas from the db as we updated it already
	if hostServiceStatusChanged {
		pending, healthy, warning, problem, err := repo.DB.GetAllServiceStatusCounts()
		if err != nil {
			log.Println(err)
			return
		}

		data := make(map[string]string)
		data["pending_count"] = strconv.Itoa(pending)
		data["healthy_count"] = strconv.Itoa(healthy)
		data["warning_count"] = strconv.Itoa(warning)
		data["problem_count"] = strconv.Itoa(problem)

		repo.broadcastMessage("public-channel", "host-service-count-changed", data)
	}

	log.Println("New status is: ", newStatus, " and message is message is: ", msg)
}

func (repo *DBRepo) broadcastMessage(channel, messageType string, data map[string]string) {

	err := app.WsClient.Trigger(channel, messageType, data)
	if err != nil {
		log.Println(err)
	}
}

func (repo *DBRepo) TestCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	oldStatus := chi.URLParam(r, "oldStatus")
	okay := true

	// get host service
	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		okay = false
		log.Println(err)
	}

	// get host?
	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		okay = false
		log.Println(err)
	}

	// test the service
	msg, newStatus := repo.testServiceForHost(h, hs)

	// update the host service in the database with status(if changed) and last check
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	hs.UpdatedAt = time.Now()
	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		okay = false
		log.Println(err)
	}

	// broadcast service status changed event

	var resp jsonResp

	if okay {
		resp = jsonResp{
			OK:            okay,
			Message:       msg,
			ServiceID:     hs.ServiceID,
			HostServiceID: hs.ID,
			HostID:        hs.HostID,
			OldStatus:     oldStatus,
			NewStatus:     newStatus,
			LastCheck:     time.Now(),
		}
	} else {
		resp.OK = false
		resp.Message = "Something went wrong"
	}

	out, _ := json.MarshalIndent(resp, "", "   ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *DBRepo) testServiceForHost(h models.Host, hs models.HostService) (string, string) {

	var msg, newStatus string

	switch hs.ServiceID {
	case HTTP:
		fmt.Println("this is case http")
		msg, newStatus = testHTTPForHost(h.URL)
		break
	}

	return msg, newStatus
}

func testHTTPForHost(url string) (string, string) {
	if strings.HasSuffix(url, "/") {
		strings.TrimSuffix(url, "/")
	}

	// replace https to http
	url = strings.Replace(url, "https://", "http://", -1)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("in the errrr")
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}
