package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
)

// Variables for test data used in the test cases
var (
	testContactUserName = "test_contacts_username"  // Test contact username
	testContactUserUUID = "test_contacts_uuid_1234" // Test contact UUID
	testUserUUID        = "test_user_uuid_1234"     // Test user UUID
)

// Test function to validate the GetContacts method of ContactsUseCase.
func TestContactsUseCase_GetContacts(t *testing.T) {
	// Structure to hold the arguments passed to GetContacts
	type args struct {
		ctx      context.Context // Context passed to the method
		userUuid string          // User UUID for which contacts are fetched
	}
	// Structure to hold test case data
	type testCase struct {
		name       string                                                                   // Test case name
		args       args                                                                     // Arguments passed to the method
		setupMocks func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) // Function to set up mocks
		want       []entity.Contacts                                                        // Expected result
		wantErr    bool                                                                     // Indicates if error is expected
	}

	// Defining an expected successful contact to be returned
	successfulGetContacts := entity.Contacts{
		UserProfile: entity.UserProfile{
			UserUUID:  testContactUserUUID, // Test contact UUID
			FirstName: "test_firstname",    // First name of the contact
			LastName:  "test_lastname",     // Last name of the contact
			Avatar:    "test_avatar",       // Avatar of the contact
		},
		ConversationUUID: "conversation_1234", // Example conversation UUID
		Blocked:          false,               // Block status of the contact
	}

	// Defining test cases for GetContacts method
	tests := []testCase{
		{
			name: "success", // Successful retrieval of contacts
			args: args{
				ctx:      context.Background(), // Context to be passed
				userUuid: testUserUUID,         // Test user UUID
			},
			// Mocking the GetContactsByUserUUID method for a successful case
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockRepo.EXPECT().
					GetContactsByUserUUID(gomock.Any(), testUserUUID).
					Return([]entity.Contacts{successfulGetContacts}, nil) // Returning success response
			},
			want:    []entity.Contacts{successfulGetContacts}, // Expected contact to be returned
			wantErr: false,                                    // No error expected
		},
		{
			name: "error fetching contacts", // Error case when fetching contacts
			args: args{
				ctx:      context.Background(),
				userUuid: testUserUUID,
			},
			// Mocking the GetContactsByUserUUID method to return an error
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockRepo.EXPECT().
					GetContactsByUserUUID(gomock.Any(), testUserUUID).
					Return(nil, fmt.Errorf("some error")) // Returning error response
			},
			want:    nil,  // No contacts expected due to error
			wantErr: true, // Error expected
		},
	}

	// Loop through each test case and run the test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // Creating a mock controller
			defer ctrl.Finish()             // Ensuring mock controller is cleaned up

			// Create mock instances for ContactsRepo and UserRepo
			mockRepo := mocks.NewMockContactsRepo(ctrl)
			mockUserRepo := mocks.NewMockUserRepo(ctrl)

			// Setup mocks for the specific test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo, mockUserRepo)
			}

			// Create instance of ContactsUseCase with mocked dependencies
			uc := &ContactsUseCase{
				repo:         mockRepo,
				userInfoRepo: mockUserRepo,
			}

			// Call the method being tested and capture the result
			got, err := uc.GetContacts(tt.args.ctx, tt.args.userUuid)
			// Verify if the returned error matches the expected error
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsUseCase.GetContacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Verify if the returned contacts match the expected contacts
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContactsUseCase.GetContacts() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test function to validate the AddContact method of ContactsUseCase
func TestContactsUseCase_AddContact(t *testing.T) {
	// Structure to hold the arguments passed to AddContact
	type args struct {
		ctx             context.Context // Context passed to the method
		contactUserName string          // Contact username to be added
		userUuid        string          // UUID of the user adding the contact
	}
	// Structure to hold test case data
	type testCase struct {
		name       string                                                                   // Test case name
		args       args                                                                     // Arguments passed to the method
		setupMocks func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) // Function to set up mocks
		wantErr    bool                                                                     // Indicates if error is expected
	}

	// Defining test cases for AddContact method
	tests := []testCase{
		{
			name: "success", // Successful addition of contact
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName, // Test contact username
				userUuid:        testUserUUID,        // Test user UUID
			},
			// Mocking the GetUserUUIDByUsername and other methods for success case
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil) // Returning success response for fetching UUID

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(false, nil) // Contact does not exist

				mockRepo.EXPECT().
					StoreContacts(gomock.Any(), gomock.Any()).
					Return(nil) // Contact successfully stored
			},
			wantErr: false, // No error expected
		},
		{
			name: "contacts username not found", // Case when the username does not exist
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			// Mocking the GetUserUUIDByUsername method to return no UUID
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(nil, nil) // Username not found
			},
			wantErr: true, // Error expected
		},
		{
			name: "contact already exists", // Case when the contact already exists
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			// Mocking the GetUserUUIDByUsername and CheckContactExist methods
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil) // Returning contact UUID

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(true, nil) // Contact already exists

				mockRepo.EXPECT().
					UpdateRemovedStatus(gomock.Any(), gomock.Any()).
					Return(nil) // Successfully updated the removed status
			},
			wantErr: false, // No error expected
		},
	}

	// Loop through each test case and run the test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // Creating a mock controller
			defer ctrl.Finish()             // Ensuring mock controller is cleaned up

			// Create mock instances for ContactsRepo and UserRepo
			mockRepo := mocks.NewMockContactsRepo(ctrl)
			mockUserRepo := mocks.NewMockUserRepo(ctrl)

			// Setup mocks for the specific test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo, mockUserRepo)
			}

			// Create instance of ContactsUseCase with mocked dependencies
			uc := &ContactsUseCase{
				repo:         mockRepo,
				userInfoRepo: mockUserRepo,
			}

			// Call the method being tested and capture the result
			err := uc.AddContact(tt.args.ctx, tt.args.contactUserName, tt.args.userUuid)
			// Verify if the returned error matches the expected error
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsUseCase.AddContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test function to validate the RemoveContact method of ContactsUseCase
func TestContactsUseCase_RemoveContact(t *testing.T) {
	// Structure to hold the arguments passed to RemoveContact
	type args struct {
		ctx             context.Context // Context passed to the method
		contactUserName string          // Contact username to be removed
		userUuid        string          // UUID of the user removing the contact
	}
	// Structure to hold test case data
	type testCase struct {
		name       string                                                                   // Test case name
		args       args                                                                     // Arguments passed to the method
		setupMocks func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) // Function to set up mocks
		wantErr    bool                                                                     // Indicates if error is expected
	}

	// Defining test cases for RemoveContact method
	tests := []testCase{
		{
			name: "success", // Successful removal of contact
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName, // Test contact username
				userUuid:        testUserUUID,        // Test user UUID
			},
			// Mocking the GetUserUUIDByUsername and other methods for success case
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil) // Returning contact UUID

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(true, nil) // Contact exists

				mockRepo.EXPECT().
					UpdateRemovedStatus(gomock.Any(), gomock.Any()).
					Return(nil) // Contact successfully removed
			},
			wantErr: false, // No error expected
		},
		{
			name: "contact does not exist", // Case when the contact does not exist
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			// Mocking the GetUserUUIDByUsername and CheckContactExist methods
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil) // Returning contact UUID

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(false, nil) // Contact does not exist
			},
			wantErr: true, // Error expected
		},
		{
			name: "username not found", // Case when the username is not found
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			// Mocking the GetUserUUIDByUsername method to return no UUID
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(nil, nil) // Username not found
			},
			wantErr: true, // Error expected
		},
	}

	// Loop through each test case and run the test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // Creating a mock controller
			defer ctrl.Finish()             // Ensuring mock controller is cleaned up

			// Create mock instances for ContactsRepo and UserRepo
			mockRepo := mocks.NewMockContactsRepo(ctrl)
			mockUserRepo := mocks.NewMockUserRepo(ctrl)

			// Setup mocks for the specific test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo, mockUserRepo)
			}

			// Create instance of ContactsUseCase with mocked dependencies
			uc := &ContactsUseCase{
				repo:         mockRepo,
				userInfoRepo: mockUserRepo,
			}

			// Call the method being tested and capture the result
			err := uc.RemoveContact(tt.args.ctx, tt.args.contactUserName, tt.args.userUuid)
			// Verify if the returned error matches the expected error
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsUseCase.RemoveContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test function to validate the UpdateBlockContact method of ContactsUseCase
func TestContactsUseCase_UpdateBlockContact(t *testing.T) {
	// Structure to hold the arguments passed to UpdateBlockContact
	type args struct {
		ctx             context.Context // Context passed to the method
		contactUserName string          // Contact username to be blocked/unblocked
		userUuid        string          // UUID of the user blocking/unblocking the contact
		block           bool            // Flag to indicate if blocking or unblocking
	}
	// Structure to hold test case data
	type testCase struct {
		name       string                                                                   // Test case name
		args       args                                                                     // Arguments passed to the method
		setupMocks func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) // Function to set up mocks
		wantErr    bool                                                                     // Indicates if error is expected
	}

	// Defining test cases for UpdateBlockContact method
	tests := []testCase{
		{
			name: "success - block contact", // Successful blocking of contact
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName, // Test contact username
				userUuid:        testUserUUID,        // Test user UUID
				block:           true,                // Block the contact
			},
			// Mocking the GetUserUUIDByUsername and other methods for success case
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil) // Returning contact UUID

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(true, nil) // Contact exists

				mockRepo.EXPECT().
					UpdateBlockedStatus(gomock.Any(), gomock.Any()).
					Return(nil) // Contact successfully blocked
			},
			wantErr: false, // No error expected
		},
		{
			name: "success - unblock contact", // Successful unblocking of contact
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName, // Test contact username
				userUuid:        testUserUUID,        // Test user UUID
				block:           false,               // Unblock the contact
			},
			// Mocking the GetUserUUIDByUsername and other methods for success case
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil) // Returning contact UUID

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(true, nil) // Contact exists

				mockRepo.EXPECT().
					UpdateBlockedStatus(gomock.Any(), gomock.Any()).
					Return(nil) // Contact successfully unblocked
			},
			wantErr: false, // No error expected
		},
		{
			name: "username not found", // Case when the username is not found
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
				block:           true, // Attempting to block the contact
			},
			// Mocking the GetUserUUIDByUsername method to return no UUID
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(nil, nil) // Username not found
			},
			wantErr: true, // Error expected
		},
		{
			name: "contact does not exist", // Case when the contact does not exist
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
				block:           true, // Attempting to block the contact
			},
			// Mocking the GetUserUUIDByUsername and CheckContactExist methods
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil) // Returning contact UUID

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(false, nil) // Contact does not exist
			},
			wantErr: true, // Error expected
		},
	}

	// Loop through each test case and run the test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // Creating a mock controller
			defer ctrl.Finish()             // Ensuring mock controller is cleaned up

			// Create mock instances for ContactsRepo and UserRepo
			mockRepo := mocks.NewMockContactsRepo(ctrl)
			mockUserRepo := mocks.NewMockUserRepo(ctrl)

			// Setup mocks for the specific test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo, mockUserRepo)
			}

			// Create instance of ContactsUseCase with mocked dependencies
			uc := &ContactsUseCase{
				repo:         mockRepo,
				userInfoRepo: mockUserRepo,
			}

			// Call the method being tested and capture the result
			err := uc.UpdateBlockContact(tt.args.ctx, tt.args.contactUserName, tt.args.userUuid, tt.args.block)
			// Verify if the returned error matches the expected error
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsUseCase.UpdateBlockContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
