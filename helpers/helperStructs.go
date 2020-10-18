package helpers

import "go_rest_api/models"

//Result used to structure the JSON response
type Result struct {
	Message  string
	Success  bool
	Instance models.Meeting
}

//MeetingList used to structure list type JSON response
type MeetingList struct {
	Message string
	Success bool
	List    []models.Meeting
}
