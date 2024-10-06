package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
)

func TestUserProfileUseCase_GetUserProfile(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx      context.Context
		userUUID string
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                             // Name of the test case, used to identify the test in the output
		args       args                               // The input arguments for the test case
		setupMocks func(mockRepo *mocks.MockUserRepo) // Function to set up mock behavior for the test
		want       entity.UserProfile                 // The expected user profile returned on success
		wantErr    bool                               // Whether the test expects an error to occur
	}

	// Define the test cases
	tests := []testCase{
		{
			name: "success - user found", // Test case name
			args: args{
				ctx:      context.Background(),
				userUUID: "user_uuid_1234",
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				userProfileDTO := &entity.UserProfileDTO{
					UserUUID:  "user_uuid_1234",
					FirstName: "Test",
					LastName:  "User",
					Avatar:    "avatar_url",
				}
				mockRepo.EXPECT().
					GetUserProfile(gomock.Any(), "user_uuid_1234").
					Return(userProfileDTO, nil)
			},
			want: entity.UserProfile{
				UserUUID:  "user_uuid_1234",
				FirstName: "Test",
				LastName:  "User",
				Avatar:    "avatar_url",
			},
			wantErr: false,
		},
		{
			name: "error - user not found", // Test case for when the user is not found
			args: args{
				ctx:      context.Background(),
				userUUID: "user_uuid_5678",
			},
			// This function sets up the mock to simulate that the user is not found
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				mockRepo.EXPECT().
					GetUserProfile(gomock.Any(), "user_uuid_5678").
					Return(nil, nil)
			},
			want:    entity.UserProfile{},
			wantErr: true,
		},
		{
			name: "error - repository error", // Test case for when there is an error in the repository
			args: args{
				ctx:      context.Background(),
				userUUID: "user_uuid_1234",
			},
			// This function sets up the mock to simulate an error occurring in the repository
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				mockRepo.EXPECT().
					GetUserProfile(gomock.Any(), "user_uuid_1234").
					Return(nil, fmt.Errorf("some error"))
			},
			want:    entity.UserProfile{},
			wantErr: true,
		},
	}

	// Iterate over each test case and run it
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create a mock instance of the UserRepo interface
			mockRepo := mocks.NewMockUserRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of UserProfileUseCase using the mock repository
			uc := &UserProfileUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			got, err := uc.GetUserProfile(tt.args.ctx, tt.args.userUUID)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("UserProfileUseCase.GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the actual output with the expected output
			if got != tt.want {
				t.Errorf("UserProfileUseCase.GetUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserProfileUseCase_UpdateUserProfile(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx         context.Context
		userProfile entity.UserProfile
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                             // Name of the test case, used to identify the test in the output
		args       args                               // The input arguments for the test case
		setupMocks func(mockRepo *mocks.MockUserRepo) // Function to set up mock behavior for the test
		wantErr    bool                               // Whether the test expects an error to occur
	}

	// Define the test cases
	tests := []testCase{
		{
			name: "success - profile updated", // Test case name
			args: args{
				ctx: context.Background(),
				userProfile: entity.UserProfile{
					UserUUID:  "user_uuid_1234",
					FirstName: "Updated",
					LastName:  "User",
					Avatar:    "updated_avatar_url",
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				userProfileDTO := entity.UserProfileDTO{
					UserUUID:  "user_uuid_1234",
					FirstName: "Updated",
					LastName:  "User",
					Avatar:    "updated_avatar_url",
				}
				mockRepo.EXPECT().
					UpdateUserProfile(gomock.Any(), userProfileDTO).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - repository error", // Test case for when there is an error in the repository
			args: args{
				ctx: context.Background(),
				userProfile: entity.UserProfile{
					UserUUID:  "user_uuid_1234",
					FirstName: "Updated",
					LastName:  "User",
					Avatar:    "updated_avatar_url",
				},
			},
			// This function sets up the mock to simulate an error occurring in the repository
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				userProfileDTO := entity.UserProfileDTO{
					UserUUID:  "user_uuid_1234",
					FirstName: "Updated",
					LastName:  "User",
					Avatar:    "updated_avatar_url",
				}
				mockRepo.EXPECT().
					UpdateUserProfile(gomock.Any(), userProfileDTO).
					Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	// Iterate over each test case and run it
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock controller for managing the lifecycle of the mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure that the mock expectations are checked and cleaned up after the test

			// Create a mock instance of the UserRepo interface
			mockRepo := mocks.NewMockUserRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of UserProfileUseCase using the mock repository
			uc := &UserProfileUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.UpdateUserProfile(tt.args.ctx, tt.args.userProfile)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("UserProfileUseCase.UpdateUserProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
