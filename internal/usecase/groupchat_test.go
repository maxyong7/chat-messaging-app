package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
)

func TestGroupChatUseCase_CreateGroupChat(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx       context.Context  // Context provides request-scoped values and cancellation signals
		groupChat entity.GroupChat // The GroupChat entity representing the chat group to be created
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                                  // Name of the test case, used to identify the test in the output
		args       args                                    // The input arguments for the test case
		setupMocks func(mockRepo *mocks.MockGroupChatRepo) // Function to set up mock behavior for the test
		wantErr    bool                                    // Whether the test expects an error to occur
	}

	// Define the test cases
	tests := []testCase{
		{
			name: "success", // Test case name
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234", // The UUID of the user creating the group
					Title:            "Group Title",    // The title of the group
					ConversationUUID: "conv_uuid_1234", // The UUID of the conversation associated with the group
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				groupChatDTO := toGroupChatDTO(entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					Title:            "Group Title",
					ConversationUUID: "conv_uuid_1234",
				})
				// Expect the CreateGroupChat method to be called with the specified DTO and return no error
				mockRepo.EXPECT().
					CreateGroupChat(gomock.Any(), groupChatDTO).
					Return(nil)
			},
			wantErr: false, // The test does not expect an error to occur
		},
		{
			name: "error creating group chat", // Test case for when an error occurs during group chat creation
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					Title:            "Group Title",
					ConversationUUID: "conv_uuid_1234",
				},
			},
			// This function sets up the mock to simulate an error when creating the group chat
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				groupChatDTO := toGroupChatDTO(entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					Title:            "Group Title",
					ConversationUUID: "conv_uuid_1234",
				})
				mockRepo.EXPECT().
					CreateGroupChat(gomock.Any(), groupChatDTO).
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

			// Create a mock instance of the GroupChatRepo interface
			mockRepo := mocks.NewMockGroupChatRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of GroupChatUseCase using the mock repository
			uc := &GroupChatUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.CreateGroupChat(tt.args.ctx, tt.args.groupChat)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupChatUseCase.CreateGroupChat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupChatUseCase_UpdateGroupTitle(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx       context.Context  // Context provides request-scoped values and cancellation signals
		groupChat entity.GroupChat // The GroupChat entity representing the chat group whose title is being updated
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                                  // Name of the test case, used to identify the test in the output
		args       args                                    // The input arguments for the test case
		setupMocks func(mockRepo *mocks.MockGroupChatRepo) // Function to set up mock behavior for the test
		wantErr    bool                                    // Whether the test expects an error to occur
	}

	// Define the test cases
	tests := []testCase{
		{
			name: "success", // Test case name
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234",  // The UUID of the user updating the group title
					ConversationUUID: "conv_uuid_1234",  // The UUID of the conversation associated with the group
					Title:            "New Group Title", // The new title of the group
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				groupChatDTO := toGroupChatDTO(entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Title:            "New Group Title",
				})
				// Simulate that the user is in the group chat and the title is successfully updated
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(true, nil)
				mockRepo.EXPECT().
					UpdateGroupTitle(gomock.Any(), groupChatDTO).
					Return(nil) // Simulate successful update of the group title
			},
			wantErr: false, // The test does not expect an error to occur
		},
		{
			name: "user not in group chat", // Test case for when the user is not part of the group chat
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Title:            "New Group Title",
				},
			},
			// This function sets up the mock to simulate that the user is not in the group chat
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(false, nil) // Simulate that the user is not in the group chat
			},
			wantErr: true, // The test expects an error to occur
		},
		{
			name: "error updating group title", // Test case for when an error occurs during title update
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Title:            "New Group Title",
				},
			},
			// This function sets up the mock to simulate an error occurring when updating the group title
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				groupChatDTO := toGroupChatDTO(entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Title:            "New Group Title",
				})
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(true, nil) // Simulate that the user is in the group chat
				mockRepo.EXPECT().
					UpdateGroupTitle(gomock.Any(), groupChatDTO).
					Return(fmt.Errorf("some error")) // Simulate an error occurring during update
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

			// Create a mock instance of the GroupChatRepo interface
			mockRepo := mocks.NewMockGroupChatRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of GroupChatUseCase using the mock repository
			uc := &GroupChatUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.UpdateGroupTitle(tt.args.ctx, tt.args.groupChat)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupChatUseCase.UpdateGroupTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupChatUseCase_AddParticipant(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx       context.Context  // Context provides request-scoped values and cancellation signals
		groupChat entity.GroupChat // The GroupChat entity representing the chat group to which participants are being added
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                                  // Name of the test case, used to identify the test in the output
		args       args                                    // The input arguments for the test case
		setupMocks func(mockRepo *mocks.MockGroupChatRepo) // Function to set up mock behavior for the test
		wantErr    bool                                    // Whether the test expects an error to occur
	}

	// Define the test cases
	tests := []testCase{
		{
			name: "success", // Test case name
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234", // The UUID of the user adding participants
					ConversationUUID: "conv_uuid_1234", // The UUID of the conversation associated with the group
					Participants: []entity.Participant{ // Participants to be added to the group
						{ParticipantUUID: "participant_uuid_1234"},
					},
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				groupChatDTO := toGroupChatDTO(entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Participants: []entity.Participant{
						{ParticipantUUID: "participant_uuid_1234"},
					},
				})
				// Simulate that the user is in the group chat and the participant is not yet in the group chat
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(true, nil)
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "participant_uuid_1234").
					Return(false, nil)
				mockRepo.EXPECT().
					AddParticipants(gomock.Any(), groupChatDTO).
					Return(nil) // Simulate successful addition of the participant
			},
			wantErr: false, // The test does not expect an error to occur
		},
		{
			name: "user not in group chat", // Test case for when the user is not part of the group chat
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Participants: []entity.Participant{
						{ParticipantUUID: "participant_uuid_1234"},
					},
				},
			},
			// This function sets up the mock to simulate that the user is not in the group chat
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(false, nil) // Simulate that the user is not in the group chat
			},
			wantErr: true, // The test expects an error to occur
		},
		{
			name: "participant already in group chat", // Test case for when the participant is already in the group chat
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Participants: []entity.Participant{
						{ParticipantUUID: "participant_uuid_1234"},
					},
				},
			},
			// This function sets up the mock to simulate that the participant is already in the group chat
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(true, nil) // Simulate that the user is in the group chat
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "participant_uuid_1234").
					Return(true, nil) // Simulate that the participant is already in the group chat
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

			// Create a mock instance of the GroupChatRepo interface
			mockRepo := mocks.NewMockGroupChatRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of GroupChatUseCase using the mock repository
			uc := &GroupChatUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.AddParticipant(tt.args.ctx, tt.args.groupChat)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupChatUseCase.AddParticipant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupChatUseCase_RemoveParticipant(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx       context.Context  // Context provides request-scoped values and cancellation signals
		groupChat entity.GroupChat // The GroupChat entity representing the chat group from which participants are being removed
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                                  // Name of the test case, used to identify the test in the output
		args       args                                    // The input arguments for the test case
		setupMocks func(mockRepo *mocks.MockGroupChatRepo) // Function to set up mock behavior for the test
		wantErr    bool                                    // Whether the test expects an error to occur
	}

	// Define the test cases
	tests := []testCase{
		{
			name: "success", // Test case name
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234", // The UUID of the user removing participants
					ConversationUUID: "conv_uuid_1234", // The UUID of the conversation associated with the group
					Participants: []entity.Participant{ // Participants to be removed from the group
						{ParticipantUUID: "participant_uuid_1234"},
					},
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				groupChatDTO := toGroupChatDTO(entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Participants: []entity.Participant{
						{ParticipantUUID: "participant_uuid_1234"},
					},
				})
				// Simulate that the user is in the group chat and the participant is also in the group chat
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(true, nil)
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "participant_uuid_1234").
					Return(true, nil)
				mockRepo.EXPECT().
					RemoveParticipants(gomock.Any(), groupChatDTO).
					Return(nil) // Simulate successful removal of the participant
			},
			wantErr: false, // The test does not expect an error to occur
		},
		{
			name: "user not in group chat", // Test case for when the user is not part of the group chat
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Participants: []entity.Participant{
						{ParticipantUUID: "participant_uuid_1234"},
					},
				},
			},
			// This function sets up the mock to simulate that the user is not in the group chat
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(false, nil) // Simulate that the user is not in the group chat
			},
			wantErr: true, // The test expects an error to occur
		},
		{
			name: "participant not in group chat", // Test case for when the participant is not part of the group chat
			args: args{
				ctx: context.Background(),
				groupChat: entity.GroupChat{
					UserUUID:         "user_uuid_1234",
					ConversationUUID: "conv_uuid_1234",
					Participants: []entity.Participant{
						{ParticipantUUID: "participant_uuid_1234"},
					},
				},
			},
			// This function sets up the mock to simulate that the participant is not in the group chat
			setupMocks: func(mockRepo *mocks.MockGroupChatRepo) {
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "user_uuid_1234").
					Return(true, nil) // Simulate that the user is in the group chat
				mockRepo.EXPECT().
					ValidateUserInGroupChat(gomock.Any(), "conv_uuid_1234", "participant_uuid_1234").
					Return(false, nil) // Simulate that the participant is not in the group chat
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

			// Create a mock instance of the GroupChatRepo interface
			mockRepo := mocks.NewMockGroupChatRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of GroupChatUseCase using the mock repository
			uc := &GroupChatUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.RemoveParticipant(tt.args.ctx, tt.args.groupChat)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupChatUseCase.RemoveParticipant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
