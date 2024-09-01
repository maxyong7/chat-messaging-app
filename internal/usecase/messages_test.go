package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
)

func TestMessageUseCase_GetMessagesFromConversation(t *testing.T) {
	type args struct {
		ctx              context.Context
		reqParam         entity.RequestParams
		conversationUUID string
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockMsgRepo *mocks.MockMessageRepo, mockReactionRepo *mocks.MockReactionRepo)
		want       []entity.GetMessageDTO
		wantErr    bool
	}

	tests := []testCase{
		{
			name: "success",
			args: args{
				ctx:              context.Background(),
				reqParam:         entity.RequestParams{UserID: "user_uuid_1234"},
				conversationUUID: "conv_uuid_1234",
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo, mockReactionRepo *mocks.MockReactionRepo) {
				reqParamDTO := entity.RequestParamsDTO(entity.RequestParams{UserID: "user_uuid_1234"})
				mockMsgRepo.EXPECT().
					GetMessages(gomock.Any(), reqParamDTO, "conv_uuid_1234").
					Return([]entity.GetMessageDTO{{MessageUUID: "msg_uuid_1234"}}, nil)

				mockReactionRepo.EXPECT().
					GetReactions(gomock.Any(), "msg_uuid_1234").
					Return([]entity.GetReactionDTO{{ReactionType: "like"}}, nil)
			},
			want: []entity.GetMessageDTO{
				{MessageUUID: "msg_uuid_1234", Reaction: []entity.GetReactionDTO{{ReactionType: "like"}}},
			},
			wantErr: false,
		},
		{
			name: "error getting messages",
			args: args{
				ctx:              context.Background(),
				reqParam:         entity.RequestParams{UserID: "user_uuid_1234"},
				conversationUUID: "conv_uuid_1234",
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo, mockReactionRepo *mocks.MockReactionRepo) {
				reqParamDTO := entity.RequestParamsDTO(entity.RequestParams{UserID: "user_uuid_1234"})
				mockMsgRepo.EXPECT().
					GetMessages(gomock.Any(), reqParamDTO, "conv_uuid_1234").
					Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error getting reactions",
			args: args{
				ctx:              context.Background(),
				reqParam:         entity.RequestParams{UserID: "user_uuid_1234"},
				conversationUUID: "conv_uuid_1234",
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo, mockReactionRepo *mocks.MockReactionRepo) {
				reqParamDTO := entity.RequestParamsDTO(entity.RequestParams{UserID: "user_uuid_1234"})
				mockMsgRepo.EXPECT().
					GetMessages(gomock.Any(), reqParamDTO, "conv_uuid_1234").
					Return([]entity.GetMessageDTO{{MessageUUID: "msg_uuid_1234"}}, nil)

				mockReactionRepo.EXPECT().
					GetReactions(gomock.Any(), "msg_uuid_1234").
					Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create mock instances of the MessageRepo and ReactionRepo interfaces
			mockMsgRepo := mocks.NewMockMessageRepo(ctrl)
			mockReactionRepo := mocks.NewMockReactionRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockMsgRepo, mockReactionRepo)
			}

			// Create an instance of MessageUseCase using the mock repositories
			uc := &MessageUseCase{
				msgRepo:      mockMsgRepo,
				reactionRepo: mockReactionRepo,
			}

			// Call the method under test with the provided arguments
			got, err := uc.GetMessagesFromConversation(tt.args.ctx, tt.args.reqParam, tt.args.conversationUUID)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageUseCase.GetMessagesFromConversation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the actual output with the expected output
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageUseCase.GetMessagesFromConversation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageUseCase_UpdateSeenStatus(t *testing.T) {
	type args struct {
		ctx        context.Context
		seenStatus entity.SeenStatus
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockMsgRepo *mocks.MockMessageRepo)
		wantErr    bool
	}

	tests := []testCase{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				seenStatus: entity.SeenStatus{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
				},
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				seenStatusDTO := entity.SeenStatusDTO{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
				}
				mockMsgRepo.EXPECT().
					UpdateSeenStatus(gomock.Any(), seenStatusDTO).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error updating seen status",
			args: args{
				ctx: context.Background(),
				seenStatus: entity.SeenStatus{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
				},
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				seenStatusDTO := entity.SeenStatusDTO{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
				}
				mockMsgRepo.EXPECT().
					UpdateSeenStatus(gomock.Any(), seenStatusDTO).
					Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create a mock instance of the MessageRepo interface
			mockMsgRepo := mocks.NewMockMessageRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockMsgRepo)
			}

			// Create an instance of MessageUseCase using the mock repository
			uc := &MessageUseCase{
				msgRepo: mockMsgRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.UpdateSeenStatus(tt.args.ctx, tt.args.seenStatus)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageUseCase.UpdateSeenStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageUseCase_GetSeenStatus(t *testing.T) {
	type args struct {
		ctx         context.Context
		messageUUID string
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockMsgRepo *mocks.MockMessageRepo)
		want       []entity.GetSeenStatusDTO
		wantErr    bool
	}

	tests := []testCase{
		{
			name: "success",
			args: args{
				ctx:         context.Background(),
				messageUUID: "msg_uuid_1234",
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				mockMsgRepo.EXPECT().
					GetSeenStatus(gomock.Any(), "msg_uuid_1234").
					Return([]entity.GetSeenStatusDTO{{SeenTimestamp: "2024-01-01T00:00:00Z"}}, nil)
			},
			want:    []entity.GetSeenStatusDTO{{SeenTimestamp: "2024-01-01T00:00:00Z"}},
			wantErr: false,
		},
		{
			name: "error getting seen status",
			args: args{
				ctx:         context.Background(),
				messageUUID: "msg_uuid_1234",
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				mockMsgRepo.EXPECT().
					GetSeenStatus(gomock.Any(), "msg_uuid_1234").
					Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create a mock instance of the MessageRepo interface
			mockMsgRepo := mocks.NewMockMessageRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockMsgRepo)
			}

			// Create an instance of MessageUseCase using the mock repository
			uc := &MessageUseCase{
				msgRepo: mockMsgRepo,
			}

			// Call the method under test with the provided arguments
			got, err := uc.GetSeenStatus(tt.args.ctx, tt.args.messageUUID)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageUseCase.GetSeenStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the actual output with the expected output
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageUseCase.GetSeenStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageUseCase_SearchMessage(t *testing.T) {
	type args struct {
		ctx              context.Context
		keyword          string
		conversationUUID string
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockMsgRepo *mocks.MockMessageRepo)
		want       []entity.SearchMessageDTO
		wantErr    bool
	}

	tests := []testCase{
		{
			name: "success",
			args: args{
				ctx:              context.Background(),
				keyword:          "hello",
				conversationUUID: "conv_uuid_1234",
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				mockMsgRepo.EXPECT().
					SearchMessage(gomock.Any(), "hello", "conv_uuid_1234").
					Return([]entity.SearchMessageDTO{{MessageUUID: "msg_uuid_1234", CreatedAt: time.Now()}}, nil)
			},
			want:    []entity.SearchMessageDTO{{MessageUUID: "msg_uuid_1234"}},
			wantErr: false,
		},
		{
			name: "error searching message",
			args: args{
				ctx:              context.Background(),
				keyword:          "hello",
				conversationUUID: "conv_uuid_1234",
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				mockMsgRepo.EXPECT().
					SearchMessage(gomock.Any(), "hello", "conv_uuid_1234").
					Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create a mock instance of the MessageRepo interface
			mockMsgRepo := mocks.NewMockMessageRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockMsgRepo)
			}

			// Create an instance of MessageUseCase using the mock repository
			uc := &MessageUseCase{
				msgRepo: mockMsgRepo,
			}

			// Call the method under test with the provided arguments
			got, err := uc.SearchMessage(tt.args.ctx, tt.args.keyword, tt.args.conversationUUID)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageUseCase.SearchMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the actual output with the expected output
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageUseCase.SearchMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageUseCase_DeleteMessage(t *testing.T) {
	type args struct {
		ctx context.Context
		msg entity.Message
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockMsgRepo *mocks.MockMessageRepo)
		want       bool
		wantErr    bool
	}

	tests := []testCase{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				msg: entity.Message{
					SenderUUID:  "user_uuid_1234",
					MessageUUID: "msg_uuid_1234",
				},
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				msgDTO := entity.MessageDTO{
					UserUUID:    "user_uuid_1234",
					MessageUUID: "msg_uuid_1234",
				}
				mockMsgRepo.EXPECT().
					ValidateMessageSentByUser(gomock.Any(), msgDTO).
					Return(true, nil)
				mockMsgRepo.EXPECT().
					DeleteMessage(gomock.Any(), msgDTO).
					Return(nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "message not sent by user",
			args: args{
				ctx: context.Background(),
				msg: entity.Message{
					SenderUUID:  "user_uuid_1234",
					MessageUUID: "msg_uuid_1234",
				},
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				msgDTO := entity.MessageDTO{
					UserUUID:    "user_uuid_1234",
					MessageUUID: "msg_uuid_1234",
				}
				mockMsgRepo.EXPECT().
					ValidateMessageSentByUser(gomock.Any(), msgDTO).
					Return(false, nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "error deleting message",
			args: args{
				ctx: context.Background(),
				msg: entity.Message{
					SenderUUID:  "user_uuid_1234",
					MessageUUID: "msg_uuid_1234",
				},
			},
			setupMocks: func(mockMsgRepo *mocks.MockMessageRepo) {
				msgDTO := entity.MessageDTO{
					UserUUID:    "user_uuid_1234",
					MessageUUID: "msg_uuid_1234",
				}
				mockMsgRepo.EXPECT().
					ValidateMessageSentByUser(gomock.Any(), msgDTO).
					Return(true, nil)
				mockMsgRepo.EXPECT().
					DeleteMessage(gomock.Any(), msgDTO).
					Return(fmt.Errorf("some error"))
			},
			want:    false,
			wantErr: true,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create a mock instance of the MessageRepo interface
			mockMsgRepo := mocks.NewMockMessageRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockMsgRepo)
			}

			// Create an instance of MessageUseCase using the mock repository
			uc := &MessageUseCase{
				msgRepo: mockMsgRepo,
			}

			// Call the method under test with the provided arguments
			got, err := uc.DeleteMessage(tt.args.ctx, tt.args.msg)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageUseCase.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the actual output with the expected output
			if got != tt.want {
				t.Errorf("MessageUseCase.DeleteMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
