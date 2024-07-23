package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// ConversationUseCase -.
type ConversationUseCase struct {
	repo         ConversationRepo
	userRepo     UserRepo
	reactionRepo ReactionRepo
	messageRepo  MessageRepo
	// webAPI  TranslationWebAPI
}

// // MessageRequest struct to hold message data
// type MessageRequest struct {
// 	MessageType string `json:"message_type"`
// 	Data        Data   `json:"data"`
// }

//	type Data struct {
//		MessageUUID      string `json:"message_uuid,omitempty"`
//		SenderUUID       string `json:"sender_uuid"`
//		ConversationUUID string `json:"conversation_uuid"`
//		ReactionData
//		MessageData
//	}
type ReactionData struct {
	ReactionType string `json:"reaction_type"`
}

// type MessageData struct {
// 	Content   string    `json:"content"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// type MessageResponse struct {
// 	MessageType  string       `json:"message_type"`
// 	ResponseData ResponseData `json:"data"`
// }

// type ResponseData struct {
// 	MessageUUID      string `json:"message_uuid,omitempty"`
// 	ConversationUUID string `json:"conversation_uuid"`
// 	ResponseReaction
// 	ResponseMessage
// 	ResponseError
// }

// type ResponseReaction struct {
// 	Reaction string `json:"reaction,omitempty"`
// }

// type ResponseMessage struct {
// 	SenderFirstName string    `json:"sender_first_name"`
// 	SenderLastName  string    `json:"sender_last_name"`
// 	Content         string    `json:"content"`
// 	CreatedAt       time.Time `json:"created_at"`
// }

// type ResponseError struct {
// 	ErrorMessage string `json:"error_msg"`
// 	SenderUUID   string `json:"sender_uuid"`
// }

// New -.
func NewConversation(r ConversationRepo, userRepo UserRepo, reactionRepo ReactionRepo, msgRepo MessageRepo) *ConversationUseCase {
	return &ConversationUseCase{
		repo:         r,
		userRepo:     userRepo,
		reactionRepo: reactionRepo,
		messageRepo:  msgRepo,
	}
}

// var hub = Hub{
// 	Clients:    make(map[string]*Client),
// 	Register:   make(chan *Client),
// 	Unregister: make(chan *Client),
// 	Broadcast:  make(chan []byte),
// }

// // function to remove client from room
// func (h *Hub) RemoveClient(client *Client) {
// 	if _, ok := h.Clients[client.ID]; ok {
// 		delete(h.Clients[client.ID], client)
// 		close(client.send)
// 		fmt.Println("Removed client")
// 	}
// }

// func (c *Client) handleConversation(convReq boundary.ConversationRequestModel, senderUUID string) {
// 	conv := entity.Conversation{
// 		SenderUUID:       senderUUID,
// 		ConversationUUID: c.ID,
// 	}
// 	switch convReq.MessageType {
// 	case sendMessageType:
// 		msg := entity.Message{
// 			MessageUUID: uuid.New().String(),
// 			Content:     convReq.Data.SendMessageRequest.Content,
// 			CreatedAt:   time.Now(),
// 		}
// 		err := c.repo.StoreConversation(conv, msg)
// 		if err != nil {
// 			fmt.Println("Conversation - handleConversation - StoreConversation err: ", err)
// 			errorMsg := c.buildErrorMessage(conv, errProcessingMessage)
// 			c.hub.Broadcast <- errorMsg
// 			break
// 		}
// 		msgResponse := conversationMessageToResponse(msg, sendMessageType)
// 		c.hub.Broadcast <- msgResponse
// 	case deleteMessageType:
// 		msg.Data.ConversationUUID = c.ID
// 		valid, err := c.messageRepo.ValidateMessageSentByUser(msg)
// 		if err != nil {
// 			fmt.Println("Conversation - readPump - ValidateMessageSentByUser err: ", err)
// 			errorMsg := c.buildErrorMessage(msg, errProcessingMessage)
// 			c.hub.Broadcast <- errorMsg
// 			break
// 		}
// 		if !valid {
// 			errorMsg := c.buildErrorMessage(msg, errOnlyAuthorCanDeleteMsg)
// 			c.hub.Broadcast <- errorMsg
// 		}
// 		err = c.messageRepo.DeleteMessage(msg)
// 		if err != nil {
// 			fmt.Println("Conversation - readPump - DeleteMessage err: ", err)
// 			errorMsg := c.buildErrorMessage(msg, errProcessingMessage)
// 			c.hub.Broadcast <- errorMsg
// 			break
// 		}
// 		msgResponse := c.buildMessageResponse(msg)
// 		c.hub.Broadcast <- msgResponse
// 	case addReactionMessageType:
// 		err = c.reactionRepo.StoreReaction(msg)
// 		if err != nil {
// 			fmt.Println("Conversation - readPump - StoreReaction err: ", err)
// 			errorMsg := c.buildErrorMessage(msg, errProcessingReaction)
// 			c.hub.Broadcast <- errorMsg
// 			break
// 		}
// 		msgResponse := c.buildMessageResponse(msg)
// 		c.hub.Broadcast <- msgResponse
// 	case removeReactionMessageType:
// 		err = c.reactionRepo.RemoveReaction(msg)
// 		if err != nil {
// 			fmt.Println("Conversation - readPump - RemoveReaction err: ", err)
// 			errorMsg := c.buildErrorMessage(msg, errProcessingReaction)
// 			c.hub.Broadcast <- errorMsg
// 			break
// 		}
// 		msgResponse := c.buildMessageResponse(msg)
// 		c.hub.Broadcast <- msgResponse
// 	}
// }

// func conversationMessageToResponse(convMsg entity.ConversationMessage, messageType string) boundary.ConversationResponseModel {
// 	return boundary.ConversationResponseModel{
// 		MessageType: messageType,
// 		Data: boundary.ConversationResponseData{
// 			Conversation: entity.Conversation{
// 				SenderUUID:          "",
// 				ConversationUUID:    "",
// 				ConversationMessage: entity.ConversationMessage{},
// 			},
// 		},
// 	}
// }

func (uc *ConversationUseCase) StoreConversation(ctx context.Context, conv entity.Conversation, msg entity.Message) error {
	err := uc.repo.StoreConversation(ctx, conv, msg)
	if err != nil {
		return fmt.Errorf("ConversationUseCase - StoreConversation - uc.repo.StoreConversation: %w", err)
	}
	return nil
}
