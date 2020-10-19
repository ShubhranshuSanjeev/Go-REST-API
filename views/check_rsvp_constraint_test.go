package views

import (
	"bytes"
	"go_rest_api/helpers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateMeetingRSVPConstraint(t *testing.T) {
	helpers.ConnectDB()
	stTime := time.Now().Local()
	edTime := time.Now().Local().Add(time.Hour*2 + time.Minute*30 + time.Second*0)

	var reqBody = []byte(`{"id":2,"title":"MongoDB","participants":[{"name":"Shubhranshu","email":"shubh@mail.com","rsvp":"Yes"},{"name":"Shivesh","email":"shivesh@mail.com","rsvp":"No"},{"name":"Umang","email":"umag@mail.com","rsvp":"No"},{"name":"Apurv","email":"apurv@mail.com","rsvp":"No"}],"startTime":"` + stTime.Format(time.RFC3339Nano) + `","endTime":"` + edTime.Format(time.RFC3339Nano) + `","creationTimestamp":"` + stTime.Format(time.RFC3339Nano) + `"}`)

	req, err := http.NewRequest("POST", "/meetings", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MeetingsCreateListAPIView)
	handler.ServeHTTP(rr, req)

	expectedResponse := `{"message":"Shubhranshu is already RSVPed in someother meeting","success":false,"instance":{"id":0,"title":"","participants":null,"startTime":"","endTime":"","creationTimestamp":""}}`
	response := rr.Body.String()
	if strings.Compare(strings.Trim(response, "\n"), expectedResponse) != 0 {
		t.Errorf("handler returned unexpected body: got %v want \n%v",
			response, expectedResponse)
	}
}
