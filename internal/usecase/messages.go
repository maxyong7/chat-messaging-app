package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type MessageUseCase struct {
	msgRepo      MessageRepo
	reactionRepo ReactionRepo
}

func NewMessage(m MessageRepo, r ReactionRepo) *MessageUseCase {
	return &MessageUseCase{
		msgRepo:      m,
		reactionRepo: r,
	}
}

func (uc *MessageUseCase) GetMessagesFromConversation(ctx context.Context, reqParam entity.RequestParams, conversationUUID string) ([]entity.GetMessageDTO, error) {
	// Convert request parameter entity object into reqParamDTO
	reqParamDTO := entity.RequestParamsDTO(reqParam)

	// Get messages from message data repository by querying 'messages' table
	// Then join 'user_info' table on 'user_uuid' to get user's firstname, lastname and avatar
	messages, err := uc.msgRepo.GetMessages(ctx, reqParamDTO, conversationUUID)
	if err != nil {
		return nil, fmt.Errorf("MessageUseCase - GetMessages - uc.msgRepo.GetMessages: %w", err)
	}
	if len(messages) == 0 {
		return nil, nil
	}

	for i, msg := range messages {
		// For each message, get reactions by querying 'reaction' table on 'message_uuid' from reaction data repository
		reactions, err := uc.reactionRepo.GetReactions(ctx, msg.MessageUUID)
		if err != nil {
			return nil, fmt.Errorf("MessageUseCase - GetMessages - uc.reactionRepo.GetReactions: %w", err)
		}
		messages[i].Reaction = reactions
	}

	return messages, nil
}

func (uc *MessageUseCase) UpdateSeenStatus(ctx context.Context, seenStatus entity.SeenStatus) error {
	// Convert seen status entity object into seenStatusDTO
	seenStatusDTO := entity.SeenStatusDTO{
		UserUUID:         seenStatus.UserUUID,
		ConversationUUID: seenStatus.ConversationUUID,
	}

	// Update 'seen_status' table in message data repository
	err := uc.msgRepo.UpdateSeenStatus(ctx, seenStatusDTO)
	if err != nil {
		return fmt.Errorf("MessageUseCase - UpdateSeenStatus - uc.msgRepo.UpdateSeenStatus: %w", err)
	}
	return nil
}

func (uc *MessageUseCase) GetSeenStatus(ctx context.Context, messageUUID string) ([]entity.GetSeenStatusDTO, error) {
	// Get seen status by querting 'seen_status' table from message data repository
	seenStatus, err := uc.msgRepo.GetSeenStatus(ctx, messageUUID)
	if err != nil {
		return nil, fmt.Errorf("MessageUseCase - GetSeenStatus - uc.msgRepo.GetSeenStatus: %w", err)
	}
	return seenStatus, nil
}

func (uc *MessageUseCase) SearchMessage(ctx context.Context, keyword string, conversationUUID string) ([]entity.SearchMessageDTO, error) {
	// Search message from message data repository by querying 'messages' table
	// Then join 'user_info' table on 'user_uuid' to get user's firstname, lastname and avatar
	messages, err := uc.msgRepo.SearchMessage(ctx, keyword, conversationUUID)
	if err != nil {
		return nil, fmt.Errorf("MessageUseCase - GetSeenStatus - uc.msgRepo.GetSeenStatus: %w", err)
	}

	for i, msg := range messages {
		// For each message found, set the cursor time to include the keyword
		inclusiveCursorTime := msg.CreatedAt.Add(-1 * time.Millisecond)
		messages[i].Cursor = encodeCursor(&inclusiveCursorTime)
	}
	return messages, nil
}

func (uc *MessageUseCase) DeleteMessage(ctx context.Context, msg entity.Message) (bool, error) {
	// Convert message entity object into msgDTO
	msgDTO := entity.MessageDTO{
		UserUUID:    msg.SenderUUID,
		MessageUUID: msg.MessageUUID,
	}

	// Validate message is sent by user in message data repository by querying 'messages' table
	valid, err := uc.msgRepo.ValidateMessageSentByUser(ctx, msgDTO)
	if err != nil {
		return valid, fmt.Errorf("MessageUseCase - DeleteMessage - uc.msgRepo.DeleteMessage: %w", err)
	}

	// If message is not sent by user, return fail validation. Will be handled by controller.
	if !valid {
		return valid, nil
	}

	// If message sent by user, delete message in 'messages' table via message data repository
	err = uc.msgRepo.DeleteMessage(ctx, msgDTO)
	if err != nil {
		return valid, fmt.Errorf("MessageUseCase - DeleteMessage - uc.msgRepo.DeleteMessage: %w", err)
	}
	return valid, nil
}

func encodeCursor(cursor *time.Time) string {
	if cursor == nil {
		return ""
	}
	if cursor.IsZero() {
		return ""
	}
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}
	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}
