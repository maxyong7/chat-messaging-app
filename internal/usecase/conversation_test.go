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

func TestConversationUseCase_GetConversationList(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx      context.Context
		reqParam entity.RequestParams
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                                     // Name of the test case
		args       args                                       // Input arguments for the test case
		setupMocks func(mockRepo *mocks.MockConversationRepo) // Function to set up mock behavior
		want       []entity.ConversationList                  // Expected output
		wantErr    bool                                       // Whether an error is expected
	}

	// List of test cases to run
	tests := []testCase{
		{
			// Test case for successful retrieval of conversations
			name: "success",
			args: args{
				ctx: context.Background(),
				reqParam: entity.RequestParams{
					UserID: testUserUUID,
				},
			},
			setupMocks: func(mockRepo *mocks.MockConversationRepo) {
				reqParamDTO := entity.RequestParamsDTO(entity.RequestParams{UserID: testUserUUID})
				mockRepo.EXPECT().
					GetConversationList(gomock.Any(), reqParamDTO).
					Return([]entity.ConversationList{{ConversationUUID: &testConversationUUID}}, nil)
			},
			want:    []entity.ConversationList{{ConversationUUID: &testConversationUUID}},
			wantErr: false,
		},
		{
			// Test case where the conversation list is empty
			name: "empty conversation list",
			args: args{
				ctx:      context.Background(),
				reqParam: entity.RequestParams{UserID: testUserUUID},
			},
			setupMocks: func(mockRepo *mocks.MockConversationRepo) {
				reqParamDTO := entity.RequestParamsDTO(entity.RequestParams{UserID: testUserUUID})
				mockRepo.EXPECT().
					GetConversationList(gomock.Any(), reqParamDTO).
					Return([]entity.ConversationList{}, nil)
			},
			want:    nil,
			wantErr: false,
		},
		{
			// Test case where an error occurs while fetching conversations
			name: "error fetching conversations",
			args: args{
				ctx:      context.Background(),
				reqParam: entity.RequestParams{UserID: testUserUUID},
			},
			setupMocks: func(mockRepo *mocks.MockConversationRepo) {
				reqParamDTO := entity.RequestParamsDTO(entity.RequestParams{UserID: testUserUUID})
				mockRepo.EXPECT().
					GetConversationList(gomock.Any(), reqParamDTO).
					Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller, which helps in managing the lifecycle of mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the controller checks the expectations and cleans up after the test

			// Create mock instances of the ConversationRepo interface
			mockRepo := mocks.NewMockConversationRepo(ctrl)

			// Set up the mock expectations using the provided setupMocks function
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of the ConversationUseCase using the mock repository
			uc := &ConversationUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			got, err := uc.GetConversationList(tt.args.ctx, tt.args.reqParam)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationUseCase.GetConversationList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the actual output with the expected output using reflect.DeepEqual
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConversationUseCase.GetConversationList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConversationUseCase_StoreConversationAndMessage(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx  context.Context
		conv entity.Conversation
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                                     // Name of the test case
		args       args                                       // Input arguments for the test case
		setupMocks func(mockRepo *mocks.MockConversationRepo) // Function to set up mock behavior
		wantErr    bool                                       // Whether an error is expected
	}

	// Define a sample time for testing
	testTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

	// List of test cases to run
	tests := []testCase{
		{
			// Test case for successfully storing a conversation and message
			name: "success",
			args: args{
				ctx: context.Background(),
				conv: entity.Conversation{
					SenderUUID:       "sender_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					MessageUUID:      "msg_uuid_1234",
					Content:          "Hello!",
					CreatedAt:        testTime,
				},
			},
			setupMocks: func(mockRepo *mocks.MockConversationRepo) {
				// Define the expected behavior of the mock
				convDTO := entity.ConversationDTO{
					SenderUUID:       "sender_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					MessageUUID:      "msg_uuid_1234",
					Content:          "Hello!",
					CreatedAt:        testTime,
				}
				mockRepo.EXPECT().
					InsertConversationAndMessage(gomock.Any(), convDTO).
					Return(nil) // Simulate success by returning nil (no error)
			},
			wantErr: false,
		},
		{
			// Test case where an error occurs while storing the conversation and message
			name: "error storing conversation and message",
			args: args{
				ctx: context.Background(),
				conv: entity.Conversation{
					SenderUUID:       "sender_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					MessageUUID:      "msg_uuid_1234",
					Content:          "Hello!",
					CreatedAt:        testTime,
				},
			},
			setupMocks: func(mockRepo *mocks.MockConversationRepo) {
				// Define the expected behavior of the mock
				convDTO := entity.ConversationDTO{
					SenderUUID:       "sender_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					MessageUUID:      "msg_uuid_1234",
					Content:          "Hello!",
					CreatedAt:        testTime,
				}
				mockRepo.EXPECT().
					InsertConversationAndMessage(gomock.Any(), convDTO).
					Return(fmt.Errorf("some error")) // Simulate an error condition
			},
			wantErr: true,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller, which helps in managing the lifecycle of mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the controller checks the expectations and cleans up after the test

			// Create a mock instance of the ConversationRepo interface
			mockRepo := mocks.NewMockConversationRepo(ctrl)

			// Set up the mock expectations using the provided setupMocks function
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of the ConversationUseCase using the mock repository
			uc := &ConversationUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.StoreConversationAndMessage(tt.args.ctx, tt.args.conv)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationUseCase.StoreConversationAndMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testConversationUUID = "conv_uuid_1234"
