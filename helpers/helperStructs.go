package helpers

import "go_rest_api/models"

//Result used to structure the JSON response
type Result struct {
	Message  string         `json:"message"`
	Success  bool           `json:"success"`
	Instance models.Meeting `json:"instance"`
}

//MeetingList used to structure list type JSON response
type MeetingList struct {
	Message string           `json:"message"`
	Success bool             `json:"success"`
	List    []models.Meeting `json:"list"`
}
