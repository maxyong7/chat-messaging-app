package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
)

func TestReactionUseCase_StoreReaction(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx      context.Context
		reaction entity.Reaction
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                                         // Name of the test case, used to identify the test in the output
		args       args                                           // The input arguments for the test case
		setupMocks func(mockReactionRepo *mocks.MockReactionRepo) // Function to set up mock behavior for the test
		wantErr    bool                                           // Whether the test expects an error to occur
	}

	// Define the test cases
	tests := []testCase{
		{
			name: "success", // Test case name
			args: args{
				ctx: context.Background(),
				reaction: entity.Reaction{
					MessageUUID:  "msg_uuid_1234",  // The UUID of the message being reacted to
					SenderUUID:   "user_uuid_1234", // The UUID of the user sending the reaction
					ReactionType: "like",           // The type of reaction being stored
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockReactionRepo *mocks.MockReactionRepo) {
				storeReactionDTO := entity.StoreReactionDTO{
					MessageUUID:  "msg_uuid_1234",
					SenderUUID:   "user_uuid_1234",
					ReactionType: "like",
				}
				// Expect the StoreReaction method to be called with the specified DTO and return no error
				mockReactionRepo.EXPECT().
					StoreReaction(gomock.Any(), storeReactionDTO).
					Return(nil)
			},
			wantErr: false, // The test does not expect an error to occur
		},
		{
			name: "error storing reaction", // Test case for when an error occurs while storing the reaction
			args: args{
				ctx: context.Background(),
				reaction: entity.Reaction{
					MessageUUID:  "msg_uuid_1234",
					SenderUUID:   "user_uuid_1234",
					ReactionType: "like",
				},
			},
			// This function sets up the mock to simulate an error when storing the reaction
			setupMocks: func(mockReactionRepo *mocks.MockReactionRepo) {
				storeReactionDTO := entity.StoreReactionDTO{
					MessageUUID:  "msg_uuid_1234",
					SenderUUID:   "user_uuid_1234",
					ReactionType: "like",
				}
				mockReactionRepo.EXPECT().
					StoreReaction(gomock.Any(), storeReactionDTO).
					Return(fmt.Errorf("some error")) // Simulate an error occurring
			},
			wantErr: true, // The test expects an error to occur
		},
	}

	// Iterate over each test case and run it
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create a mock instance of the ReactionRepo interface
			mockReactionRepo := mocks.NewMockReactionRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockReactionRepo)
			}

			// Create an instance of ReactionUseCase using the mock repository
			uc := &ReactionUseCase{
				reactionRepo: mockReactionRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.StoreReaction(tt.args.ctx, tt.args.reaction)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("ReactionUseCase.StoreReaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReactionUseCase_RemoveReaction(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx      context.Context
		reaction entity.Reaction
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                                         // Name of the test case, used to identify the test in the output
		args       args                                           // The input arguments for the test case
		setupMocks func(mockReactionRepo *mocks.MockReactionRepo) // Function to set up mock behavior for the test
		wantErr    bool                                           // Whether the test expects an error to occur
	}

	// Define the test cases
	tests := []testCase{
		{
			name: "success", // Test case name
			args: args{
				ctx: context.Background(),
				reaction: entity.Reaction{
					MessageUUID: "msg_uuid_1234",  // The UUID of the message from which the reaction is being removed
					SenderUUID:  "user_uuid_1234", // The UUID of the user removing the reaction
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockReactionRepo *mocks.MockReactionRepo) {
				removeReactionDTO := entity.RemoveReactionDTO{
					MessageUUID: "msg_uuid_1234",
					SenderUUID:  "user_uuid_1234",
				}
				// Expect the RemoveReaction method to be called with the specified DTO and return no error
				mockReactionRepo.EXPECT().
					RemoveReaction(gomock.Any(), removeReactionDTO).
					Return(nil)
			},
			wantErr: false, // The test does not expect an error to occur
		},
		{
			name: "error removing reaction", // Test case for when an error occurs while removing the reaction
			args: args{
				ctx: context.Background(),
				reaction: entity.Reaction{
					MessageUUID: "msg_uuid_1234",
					SenderUUID:  "user_uuid_1234",
				},
			},
			// This function sets up the mock to simulate an error when removing the reaction
			setupMocks: func(mockReactionRepo *mocks.MockReactionRepo) {
				removeReactionDTO := entity.RemoveReactionDTO{
					MessageUUID: "msg_uuid_1234",
					SenderUUID:  "user_uuid_1234",
				}
				mockReactionRepo.EXPECT().
					RemoveReaction(gomock.Any(), removeReactionDTO).
					Return(fmt.Errorf("some error")) // Simulate an error occurring
			},
			wantErr: true, // The test expects an error to occur
		},
	}

	// Iterate over each test case and run it
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create a mock instance of the ReactionRepo interface
			mockReactionRepo := mocks.NewMockReactionRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockReactionRepo)
			}

			// Create an instance of ReactionUseCase using the mock repository
			uc := &ReactionUseCase{
				reactionRepo: mockReactionRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.RemoveReaction(tt.args.ctx, tt.args.reaction)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("ReactionUseCase.RemoveReaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
