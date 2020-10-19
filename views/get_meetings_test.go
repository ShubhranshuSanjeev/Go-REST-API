package views

import (
	"go_rest_api/helpers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetMeetingsByTimeRange(t *testing.T) {
	helpers.ConnectDB()

	req, err := http.NewRequest("GET", "/meetings", nil)
	if err != nil {
		t.Fatal(err)
	}

	stTime := time.Now().Add(time.Hour*(-2) + time.Minute*30 + time.Second*0)
	edTime := time.Now().Add(time.Hour*2 + time.Minute*30 + time.Second*0)

	params := req.URL.Query()
	params.Add("start", stTime.Format(time.RFC3339Nano))
	params.Add("end", edTime.Format(time.RFC3339Nano))
	req.URL.RawQuery = params.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MeetingsCreateListAPIView)
	handler.ServeHTTP(rr, req)

	expectedResponse := `{"message":"List of Meetings","success":true,"list":[{"id":1,"title":"Golang","participants":[{"name":"Shubhranshu","email":"shubh@mail.com","rsvp":"Yes"},{"name":"Shivesh","email":"shivesh@mail.com","rsvp":"No"},{"name":"Umang","email":"umag@mail.com","rsvp":"No"},{"name":"Apurv","email":"apurv@mail.com","rsvp":"No"}],"startTime":"2020-10-19T12:37:04.4319963+05:30","endTime":"2020-10-19T15:07:04.4319963+05:30","creationTimestamp":"2020-10-19T12:37:04.4319963+05:30"}]}`
	response := rr.Body.String()
	if strings.Compare(strings.Trim(response, "\n"), expectedResponse) != 0 {
		t.Errorf("handler returned unexpected body: got %v want \n%v",
			response, expectedResponse)
	}
}

func TestGetMeetingsByID(t *testing.T) {
	helpers.ConnectDB()

	req, err := http.NewRequest("GET", "/meeting/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMeetingByIDAPIView)
	handler.ServeHTTP(rr, req)

	expectedResponse := `{"message":"Found","success":true,"instance":{"id":1,"title":"Golang","participants":[{"name":"Shubhranshu","email":"shubh@mail.com","rsvp":"Yes"},{"name":"Shivesh","email":"shivesh@mail.com","rsvp":"No"},{"name":"Umang","email":"umag@mail.com","rsvp":"No"},{"name":"Apurv","email":"apurv@mail.com","rsvp":"No"}],"startTime":"2020-10-19T12:37:04.4319963+05:30","endTime":"2020-10-19T15:07:04.4319963+05:30","creationTimestamp":"2020-10-19T12:37:04.4319963+05:30"}}`
	response := rr.Body.String()
	if strings.Compare(strings.Trim(response, "\n"), expectedResponse) != 0 {
		t.Errorf("handler returned unexpected body: got %v want \n%v",
			response, expectedResponse)
	}
}

func TestGetMeetingsByParticipantEmail(t *testing.T) {
	helpers.ConnectDB()

	req, err := http.NewRequest("GET", "/meetings", nil)
	if err != nil {
		t.Fatal(err)
	}

	participant := "shubh@mail.com"

	params := req.URL.Query()
	params.Add("participant", participant)
	req.URL.RawQuery = params.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MeetingsCreateListAPIView)
	handler.ServeHTTP(rr, req)

	expectedResponse := `{"message":"List of Meetings","success":true,"list":[{"id":1,"title":"Golang","participants":[{"name":"Shubhranshu","email":"shubh@mail.com","rsvp":"Yes"},{"name":"Shivesh","email":"shivesh@mail.com","rsvp":"No"},{"name":"Umang","email":"umag@mail.com","rsvp":"No"},{"name":"Apurv","email":"apurv@mail.com","rsvp":"No"}],"startTime":"2020-10-19T12:37:04.4319963+05:30","endTime":"2020-10-19T15:07:04.4319963+05:30","creationTimestamp":"2020-10-19T12:37:04.4319963+05:30"},{"id":2,"title":"Golang","participants":[{"name":"Shubhranshu","email":"shubh@mail.com","rsvp":"No"},{"name":"Shivesh","email":"shivesh@mail.com","rsvp":"No"},{"name":"Umang","email":"umag@mail.com","rsvp":"No"},{"name":"Apurv","email":"apurv@mail.com","rsvp":"No"}],"startTime":"2020-10-18T20:14:05.0630561+05:30","endTime":"2020-10-18T21:14:05.0630561+05:30","creationTimestamp":"2020-10-18T19:14:05.0630561+05:30"}]}`
	response := rr.Body.String()
	if strings.Compare(strings.Trim(response, "\n"), expectedResponse) != 0 {
		t.Errorf("handler returned unexpected body: got %v want \n%v",
			response, expectedResponse)
	}
}
