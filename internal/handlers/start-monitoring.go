package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

//""

// this is a unit of work where we store the host-service id
type job struct {
	HostServiceID int
}

// we don't use a pointer receiver for the receiver
// which is unusual but appropriate here
func (j job) Run() {
	Repo.ScheduledCheck(j.HostServiceID)
}

func (repo *DBRepo) StartMonitoring() {
	if app.PreferenceMap["monitoring_live"] == "1" {
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
				sch = fmt.Sprintf("@every %d%s", x.ScheduleNumber*24, "h")

			} else {
				sch = fmt.Sprintf("@every %d%s", x.ScheduleNumber, x.ScheduleUnit)
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
			payload["host_service_id"] = strconv.Itoa(x.ID)
			// year 0001 then some numbers for the days, time
			// the thing is in Go we don't want to deal with empty value
			// so instead of we put the year 0001
			yearOne := time.Date(0001, 11, 17, 20, 34, 58, 65138737, time.UTC)
			if app.Scheduler.Entry(app.MonitorMap[x.ID]).Next.After(yearOne) {
				payload["next_run"] = app.Scheduler.Entry(app.MonitorMap[x.ID]).Next.Format("2006-01-03 3:04:05 PM")
			} else {
				payload["next_run"] = "Pending...."
			}
			payload["host"] = x.HostName
			payload["service"] = x.Service.ServiceName
			if x.LastCheck.After(yearOne) {
				payload["last_run"] = x.LastCheck.Format("2006-01-03 3:04:05 PM")
			} else {
				payload["next_run"] = "Pending...."
			}
			payload["schedule"] = fmt.Sprintf("@every %d%s", x.ScheduleNumber, x.ScheduleUnit)

			err = app.WsClient.Trigger("public-channel", "next-run-event", payload)
			if err != nil {
				log.Println(err)
			}

			err = app.WsClient.Trigger("public-channel", "schedule-changed-event", payload)
			if err != nil {
				log.Println(err)
			}

			// end of range
		}
	}
}
