package main

import (
//""
)

// this is a unit of work where we store the host-service id
type job struct {
	HostServiceID int
}

// we don't use a pointer receiver for the receiver
// which is unusual but appropriate here
func (j job) run() {

}

func startMonitoring() {
	if preferenceMap["monitoring_live"] == "1" {
		data := make(map[string]string)
		data["message"] = "starting"

		// TODO trigger a message to broadcast to all clients
		// that app is starting to monitor

		// get all of the services that we want to monitor

		// range throught the services
		// -> get the schedule unit and number
		// -> create a job
		// -> save the id of the job so we can start/stop it
		// -> broadcast over websockets the fact that the service is scheduled
		// end of range
	}
}
