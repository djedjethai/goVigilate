package main

import (
	"fmt"
	"log"
	"strconv"
)

//""

// this is a unit of work where we store the host-service id
type job struct {
	HostServiceID int
}

// we don't use a pointer receiver for the receiver
// which is unusual but appropriate here
func (j job) run() {
	repo.ScheduledCheck(j.HostServiceID)
}

func startMonitoring() {
	if preferenceMap["monitoring_live"] == "1" {
		// trigger a message to broadcast to all clients
		data := make(map[string]string)
		data["message"] = "Monitoring is starting..."

		// that app is starting to monitor
		err := app.WsClient.Trigger("public-channel", "app-starting", data)
		if err != nil {
			log.Println(err)
		}

		// get all of the services that we want to monitor
		servicesToMonitor, err := repo.DB.GetServicesToMonitor()
		if err != nil {
			log.Println(err)
		}

		// range throught the services
		for _, x := range servicesToMonitor {
			// -> get the schedule unit and number
			var sch string
			// if in db the schedule unit has been set to "d"(day)
			if x.ScheduleUnit == "d" {
				// see the doc of the cron package
				sch := fmt.Sprintf("@every %d%s", x.ScheduleNumber*24, "h")

			} else {
				sch := fmt.Sprintf("@every %d%s", x.ScheduleNumber, x.ScheduleUnit)
			}

			// create a job
			var j job
			j.HostServiceID = x.ID
			scheduleID, err := app.Scheduler.AddJob(sch, j)
			if err != nil {
				log.Println(err)
			}

			// save the id of the job so we can start/stop it
			app.MonitorMap[x.ID] = scheduleID

			// broadcast over websockets the fact that the service is scheduled
			payload := make(map[string]string)
			payload["message"] = "scheduling"
			data["host_service_id"] = strconv.Itoa(x.ID)
			yearone := time.Date(".....to finish here ......")
			err := app.WsClient.Trigger("public-channel", "test-event", "monitoring")

			// end of range
		}
	}
}
