package models

//Meeting model for the database
type Meeting struct {
	ID                int           `json:"id" bson:"_id,omitempty"`
	Title             string        `json:"title" bson:"title,omitempty"`
	Participants      []Participant `json:"participants" bson:"participants, omitempty"`
	StartTime         string        `json:"startTime" bson:"startTime,omitempty"`
	EndTime           string        `json:"endTime" bson:"endTime,omitempty"`
	CreationTimestamp string        `json:"creationTimestamp" bson:"creationTimestamp,omitempty"`
}

//Participant model for the database
type Participant struct {
	Name  string `json:"name" bson:"name,omitempty"`
	Email string `json:"email" bson:"_id,omitempty"`
	RSVP  string `json:"rsvp" bson:"rsvp,omitempty"`
}

//ParticipantMeeting Participant and Meeting Relation
type ParticipantMeeting struct {
	MeetingID        int    `json:"meetingId" bson:"meetingId,omitempty"`
	ParticipantEmail string `json:"participantEmail" bson:"participantEmail,omitempty"`
	RSVP             string `json:"rsvp" bson:"rsvp,omitempty"`
}
