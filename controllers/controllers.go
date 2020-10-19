package controllers

import (
	"context"
	"fmt"
	"go_rest_api/helpers"
	"go_rest_api/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateMeeting Function to create a meeting
func CreateMeeting(newEntry models.Meeting) helpers.Result {
	var newMeeting models.Meeting
	var newParticipant models.Participant
	var newParticipantMeetingRealtion models.ParticipantMeeting
	var insertParticipants []interface{}

	meetingsCollection := helpers.DBClient.Database("appointy").Collection("meetings")
	participantCollection := helpers.DBClient.Database("appointy").Collection("participants")
	participantMeetingCollection := helpers.DBClient.Database("appointy").Collection("participantMeeting")

	for _, participant := range newEntry.Participants {

		//Checking if the email of the given pariticipant is already in use
		var result models.Participant
		filter := bson.M{"_id": participant.Email}
		err := participantCollection.FindOne(context.TODO(), filter).Decode(&result)

		if err != nil {

			//If not in use then add the participant to Participant Collection
			newParticipant = models.Participant{
				Name:  participant.Name,
				Email: participant.Email,
				RSVP:  participant.RSVP,
			}
			insertParticipants = append(insertParticipants, newParticipant)

		} else if result.Name != participant.Name {

			// If already in use then check the names of the user,
			// if the name is different then its a different user,
			// so respond with an error message that the email is already in use
			var message string = "User with email-id " + participant.Email + " already exists"
			res := helpers.Result{
				Message:  message,
				Success:  false,
				Instance: models.Meeting{},
			}
			return res

		} else if result.Name == participant.Name {

			//Fetching all the meetings of the participant
			filter2 := bson.M{"participantEmail": participant.Email}
			scheduledMeetings, err := participantMeetingCollection.Find(context.TODO(), filter2)

			if err != nil {
				fmt.Println(err.Error())
				continue
			} else {

				//Parsing the string to type time.Time
				newStartTime, _ := time.Parse(time.RFC3339Nano, newEntry.StartTime)
				newEndTime, _ := time.Parse(time.RFC3339Nano, newEntry.EndTime)

				for scheduledMeetings.Next(context.TODO()) {
					var relation models.ParticipantMeeting
					err := scheduledMeetings.Decode(&relation)
					if err != nil {
						log.Fatal(err)
					}

					var meet models.Meeting
					filter3 := bson.M{"_id": relation.MeetingID}
					err = meetingsCollection.FindOne(context.TODO(), filter3).Decode(&meet)
					if err != nil {
						log.Fatal(err)
					}

					//Parsing the string to type time.Time
					meetStartTime, _ := time.Parse(time.RFC3339Nano, meet.StartTime)
					meetEndTime, _ := time.Parse(time.RFC3339Nano, meet.EndTime)

					// Checking if there is a time overlap between the meeting the participant,
					// is already a part of and the new meeting he is being added to, and if
					// he is RSVPed Yes in both the meetings
					if (((meetStartTime.Before(newStartTime) || meetStartTime.Equal(newStartTime)) &&
						(meetEndTime.After(newStartTime) || meetEndTime.Equal(newStartTime))) ||
						((meetStartTime.Before(newEndTime) || meetStartTime.Equal(newEndTime)) &&
							(meetEndTime.After(newEndTime) || meetEndTime.Equal(newEndTime)))) &&
						relation.RSVP == "Yes" && participant.RSVP == "Yes" {

						//If there is an overlap then return an error message.
						var message string = participant.Name + " is already RSVPed in someother meeting"
						res := helpers.Result{
							Message:  message,
							Success:  false,
							Instance: models.Meeting{},
						}
						return res
					}
				}
			}
		}
	}

	newMeeting = models.Meeting{
		ID:                newEntry.ID,
		Title:             newEntry.Title,
		Participants:      newEntry.Participants,
		StartTime:         newEntry.StartTime,
		EndTime:           newEntry.EndTime,
		CreationTimestamp: newEntry.CreationTimestamp,
	}

	//Creating a new meeting record
	_, err := meetingsCollection.InsertOne(context.TODO(), newMeeting)
	if err != nil {
		res := helpers.Result{
			Message:  err.Error(),
			Success:  false,
			Instance: models.Meeting{},
		}
		return res
	}

	if len(insertParticipants) > 0 {
		_, err = participantCollection.InsertMany(context.TODO(), insertParticipants)
		if err != nil {
			res := helpers.Result{
				Message:  err.Error(),
				Success:  false,
				Instance: models.Meeting{},
			}
			return res
		}
	}

	for _, participant := range newEntry.Participants {
		newParticipantMeetingRealtion = models.ParticipantMeeting{
			MeetingID:        newEntry.ID,
			ParticipantEmail: participant.Email,
			RSVP:             participant.RSVP,
		}
		_, err = participantMeetingCollection.InsertOne(context.TODO(), newParticipantMeetingRealtion)
		if err != nil {
			res := helpers.Result{
				Message:  err.Error(),
				Success:  false,
				Instance: models.Meeting{},
			}
			return res
		}
	}

	// Returning the newly created meeting
	var instance models.Meeting
	filter := bson.D{primitive.E{Key: "_id", Value: newEntry.ID}}
	meetingsCollection.FindOne(context.TODO(), filter).Decode(&instance)
	res := helpers.Result{
		Message:  "New Meeting Scheduled",
		Success:  true,
		Instance: instance,
	}

	return res
}

//GetMeeting function to fetch a meeting by its ID
func GetMeeting(id int64) helpers.Result {
	meetingsCollection := helpers.DBClient.Database("appointy").Collection("meetings")

	var instance models.Meeting
	var res helpers.Result

	//Fetching a meeting by its ID
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := meetingsCollection.FindOne(context.TODO(), filter).Decode(&instance)

	if err != nil {
		res = helpers.Result{
			Message:  err.Error(),
			Success:  false,
			Instance: models.Meeting{},
		}
	} else {
		res = helpers.Result{
			Message:  "Found",
			Success:  true,
			Instance: instance,
		}
	}

	return res
}

//GetMeetingsByTimeRange function to get list of meetings within a time range
func GetMeetingsByTimeRange(st time.Time, ed time.Time) helpers.MeetingList {
	meetingsCollection := helpers.DBClient.Database("appointy").Collection("meetings")
	filter := bson.D{}
	cur, err := meetingsCollection.Find(context.TODO(), filter)

	if err != nil {
		res := helpers.MeetingList{
			Message: err.Error(),
			Success: false,
			List:    []models.Meeting{},
		}
		return res
	}

	var meets []models.Meeting
	for cur.Next(context.TODO()) {
		var tmp models.Meeting
		err := cur.Decode(&tmp)

		if err != nil {
			continue
		}
		//Checking if the start and end time falls within the given time range
		tmpStartTime, _ := time.Parse(time.RFC3339Nano, tmp.StartTime)
		tmpEndTime, _ := time.Parse(time.RFC3339Nano, tmp.EndTime)

		if (tmpStartTime.After(st) || tmpStartTime.Equal(st)) &&
			(tmpEndTime.Before(ed) || tmpEndTime.Equal(ed)) {
			meets = append(meets, tmp)
		}
	}

	res := helpers.MeetingList{
		Message: "List of Meetings",
		Success: true,
		List:    meets,
	}

	return res
}

//GetParticipantMeetings function to return a list of meetings of a participant
func GetParticipantMeetings(email string) helpers.MeetingList {
	meetingsCollection := helpers.DBClient.Database("appointy").Collection("meetings")
	participantMeetingCollection := helpers.DBClient.Database("appointy").Collection("participantMeeting")

	// Fetching all the meetings the participant is a part of
	filter := bson.M{"participantEmail": bson.M{"$eq": email}}

	cur, err := participantMeetingCollection.Find(context.TODO(), filter)

	if err != nil {
		res := helpers.MeetingList{
			Message: err.Error(),
			Success: false,
			List:    []models.Meeting{},
		}
		return res
	}

	var meets []models.Meeting
	for cur.Next(context.TODO()) {

		var tmp models.ParticipantMeeting
		err := cur.Decode(&tmp)
		if err != nil {
			continue
		}

		var meet models.Meeting

		//Fetching each meeting details
		err = meetingsCollection.FindOne(context.TODO(), bson.M{"_id": tmp.MeetingID}).Decode(&meet)
		if err != nil {
			continue
		}

		meets = append(meets, meet)
	}

	res := helpers.MeetingList{
		Message: "List of Meetings",
		Success: true,
		List:    meets,
	}
	return res
}
