package handlers

import (
	"github.com/pusher/pusher-http-go"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// this handler is protected as only authenticated user have access to it
// the connection which arrive here(from the browser),
// request to ugrade the connection to ws
func (repo *DBRepo) PusherAuth(w http.ResponseWriter, r *http.Request) {
	// get userID from the session
	userID := repo.App.Session.GetInt(r.Context(), "userID")
	u, _ := repo.DB.GetUserById(userID)

	// from the browser
	params, _ := ioutil.ReadAll(r.Body)

	// identicate the user to the pusher
	presenceData := pusher.MemberData{
		UserID: strconv.Itoa(userID),
		UserInfo: map[string]string{
			"name": u.FirstName,
			"id":   strconv.Itoa(userID),
		},
	}

	response, err := app.WsClient.AuthenticatePresenceChannel(params, presenceData)
	if err != nil {
		log.Println(err)
		return
	}

	// we respond to the browser, which activate the socket connection
	// the client will be connected(directly) to the pusher
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(response)
}

func (repo *DBRepo) TestPusher(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["message"] = "Hello world"

	// trigger the public-channel(the only one we setted up)
	// means that any one who reach this handler will trigger
	// the event(test-event) link to it "on the browser page" listening to this channel
	err := repo.App.WsClient.Trigger("public-channel", "test-event", data)
	if err != nil {
		log.Println(err)
	}
}
