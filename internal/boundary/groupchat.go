package boundary

import "github.com/maxyong7/chat-messaging-app/internal/entity"

type GroupChatCreationForm struct {
	Title            string             `json:"title"`
	ConversationUUID string             `json:"conversation_uuid"`
	Participants     []ParticipantsForm `json:"participants"`
}

type ParticipantsForm struct {
	Username        string `json:"username"`
	ParticipantUUID string `json:"participant_uuid"`
}

func (r GroupChatCreationForm) ToGroupChat(userUUID string) entity.GroupChat {
	participantsEntity := []entity.Participant{}
	for _, p := range r.Participants {
		participantEntity := entity.Participant{
			Username:        p.Username,
			ParticipantUUID: p.ParticipantUUID,
		}
		participantsEntity = append(participantsEntity, participantEntity)

	}
	return entity.GroupChat{
		UserUUID:         userUUID,
		Title:            r.Title,
		ConversationUUID: r.ConversationUUID,
		Participants:     participantsEntity,
	}
}
