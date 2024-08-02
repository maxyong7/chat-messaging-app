package boundary

import "github.com/maxyong7/chat-messaging-app/internal/entity"

type GroupChatRequestModel struct {
	Title            string                    `json:"title"`
	ConversationUUID string                    `json:"conversation_uuid"`
	Participants     []ParticipantRequestModel `json:"participants"`
}

type ParticipantRequestModel struct {
	Username        string `json:"username"`
	ParticipantUUID string `json:"participant_uuid"`
}

func (r GroupChatRequestModel) ToGroupChat(userUUID string) entity.GroupChat {
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
