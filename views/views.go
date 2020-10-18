package views

import (
	"encoding/json"
	"go_rest_api/controllers"
	"go_rest_api/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//MeetingsCreateListAPIView View for invoking creation and fetching of records
func MeetingsCreateListAPIView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		queries := r.URL.Query()

		if queries.Get("participant") != "" {

			//get meetings by participant email id
			participant := queries["participant"][0]
			res := controllers.GetParticipantMeetings(participant)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusAccepted)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)
			}

		} else if queries.Get("start") != "" && queries.Get("end") != "" {

			//Parsing string to time.Time
			start, _ := time.Parse(time.RFC3339Nano, queries["start"][0])
			end, _ := time.Parse(time.RFC3339Nano, queries["end"][0])

			//get meeting within time range.
			res := controllers.GetMeetingsByTimeRange(start, end)

			//JSON response
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusAccepted)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)
			}

		}

	} else if r.Method == "POST" {

		//Reading the contents of body of the HTTP packet
		var meeting models.Meeting
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		//De-Searializing JSON data read from the HTTP packet to meeting of type "strcut Meeting"
		if err := json.Unmarshal(body, &meeting); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		//Calling function to schedule a new meeting
		response := controllers.CreateMeeting(meeting)

		//JSON response
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}
}

//GetMeetingByIDAPIView View to get any meeting by its ID
func GetMeetingByIDAPIView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		urlString := strings.Split(r.URL.Path, "/")

		if len(urlString) != 3 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(404)
			if err := json.NewEncoder(w).Encode(nil); err != nil {
				panic(err)
			}
		} else {
			meetingID, _ := strconv.ParseInt(urlString[2], 10, 32)

			//Fetching meeting with id = meetingID
			response := controllers.GetMeeting(meetingID)

			//JSON response
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				panic(err)
			}
		}
	}
}
