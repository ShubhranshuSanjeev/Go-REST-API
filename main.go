package main

import (
	"go_rest_api/helpers"
	"go_rest_api/views"
	"log"
	"net/http"
)

func main() {
	helpers.ConnectDB()

	http.HandleFunc("/meetings", views.MeetingsCreateListAPIView)
	http.HandleFunc("/meeting/", views.GetMeetingByIDAPIView)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
