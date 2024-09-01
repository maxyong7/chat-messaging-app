package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUseCase_VerifyCredentials(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx             context.Context
		userCredentials entity.UserCredentials
	}

	// Define the structure of each test case
	type testCase struct {
		name       string                             // Name of the test case, used to identify the test in the output
		args       args                               // The input arguments for the test case
		setupMocks func(mockRepo *mocks.MockUserRepo) // Function to set up mock behavior for the test
		wantUUID   string                             // The expected UUID returned on success
		wantMatch  bool                               // Whether the credentials are expected to match
		wantErr    bool                               // Whether the test expects an error to occur
	}

	// Example hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Define the test cases
	tests := []testCase{
		{
			name: "success - credentials match", // Test case name
			args: args{
				ctx: context.Background(),
				userCredentials: entity.UserCredentials{
					Username: "testuser",
					Password: "password123",
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				userCredentialsDTO := entity.UserCredentialsDTO{
					Username: "testuser",
					Password: "password123",
				}
				mockRepo.EXPECT().
					GetUserCredentials(gomock.Any(), userCredentialsDTO).
					Return(&entity.UserCredentialsDTO{UserUuid: "user_uuid_1234", Password: string(hashedPassword)}, nil)
			},
			wantUUID:  "user_uuid_1234",
			wantMatch: true,
			wantErr:   false,
		},
		{
			name: "error - user not found", // Test case for when the user is not found
			args: args{
				ctx: context.Background(),
				userCredentials: entity.UserCredentials{
					Username: "unknownuser",
					Password: "password123",
				},
			},
			// This function sets up the mock to simulate that the user is not found
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				userCredentialsDTO := entity.UserCredentialsDTO{
					Username: "unknownuser",
					Password: "password123",
				}
				mockRepo.EXPECT().
					GetUserCredentials(gomock.Any(), userCredentialsDTO).
					Return(nil, nil)
			},
			wantUUID:  "",
			wantMatch: false,
			wantErr:   true,
		},
		{
			name: "error - incorrect password", // Test case for when the password does not match
			args: args{
				ctx: context.Background(),
				userCredentials: entity.UserCredentials{
					Username: "testuser",
					Password: "wrongpassword",
				},
			},
			// This function sets up the mock to simulate that the password is incorrect
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				userCredentialsDTO := entity.UserCredentialsDTO{
					Username: "testuser",
					Password: "wrongpassword",
				}
				mockRepo.EXPECT().
					GetUserCredentials(gomock.Any(), userCredentialsDTO).
					Return(&entity.UserCredentialsDTO{UserUuid: "user_uuid_1234", Password: string(hashedPassword)}, nil)
			},
			wantUUID:  "",
			wantMatch: false,
			wantErr:   true,
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

			// Create an instance of LoginUseCase using the mock repository
			uc := &LoginUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			gotUUID, gotMatch, err := uc.VerifyCredentials(tt.args.ctx, tt.args.userCredentials)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUseCase.VerifyCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check if the UUID returned matches the expected UUID
			if gotUUID != tt.wantUUID {
				t.Errorf("LoginUseCase.VerifyCredentials() gotUUID = %v, want %v", gotUUID, tt.wantUUID)
			}

			// Check if the password match status matches the expected value
			if gotMatch != tt.wantMatch {
				t.Errorf("LoginUseCase.VerifyCredentials() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}

func TestLoginUseCase_RegisterUser(t *testing.T) {
	// Define the input arguments for the method being tested
	type args struct {
		ctx              context.Context
		userRegistration entity.UserRegistration
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
			name: "success - user registration", // Test case name
			args: args{
				ctx: context.Background(),
				userRegistration: entity.UserRegistration{
					UserCredentials: entity.UserCredentials{
						Username: "newuser",
						Password: "newpassword",
						Email:    "newuser@example.com",
					},
					FirstName: "New",
					LastName:  "User",
					Avatar:    "",
				},
			},
			// This function sets up the expected behavior of the mock repository for this test case
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				// Simulate that the user does not already exist and the user info is stored successfully
				mockRepo.EXPECT().
					CheckUserExist(gomock.Any(), gomock.Any()).
					Return(false, nil)
				mockRepo.EXPECT().
					StoreUserInfo(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false, // The test does not expect an error to occur
		},
		{
			name: "error - user already exists", // Test case for when the user already exists
			args: args{
				ctx: context.Background(),
				userRegistration: entity.UserRegistration{
					UserCredentials: entity.UserCredentials{
						Username: "newuser",
						Password: "newpassword",
						Email:    "newuser@example.com",
					},
					FirstName: "Existing",
					LastName:  "User",
				},
			},
			// This function sets up the mock to simulate that the user already exists
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				mockRepo.EXPECT().
					CheckUserExist(gomock.Any(), gomock.Any()).
					Return(true, nil)
			},
			wantErr: true, // The test expects an error to occur
		},
		{
			name: "error - storing user info", // Test case for when there is an error storing the user info
			args: args{
				ctx: context.Background(),
				userRegistration: entity.UserRegistration{
					UserCredentials: entity.UserCredentials{
						Username: "newuser",
						Password: "newpassword",
						Email:    "newuser@example.com",
					},
					FirstName: "New",
					LastName:  "User",
				},
			},
			// This function sets up the mock to simulate an error when storing user info
			setupMocks: func(mockRepo *mocks.MockUserRepo) {
				mockRepo.EXPECT().
					CheckUserExist(gomock.Any(), gomock.Any()).
					Return(false, nil)
				mockRepo.EXPECT().
					StoreUserInfo(gomock.Any(), gomock.Any()).
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

			// Create a mock instance of the UserRepo interface
			mockRepo := mocks.NewMockUserRepo(ctrl)

			// Set up the mock expectations using the setupMocks function provided in the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create an instance of LoginUseCase using the mock repository
			uc := &LoginUseCase{
				repo: mockRepo,
			}

			// Call the method under test with the provided arguments
			err := uc.RegisterUser(tt.args.ctx, tt.args.userRegistration)

			// Check if the error status matches the expected value
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUseCase.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
