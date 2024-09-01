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

var (
	testContactUserName = "test_contacts_username"
	testContactUserUUID = "test_contacts_uuid_1234"
	testUserUUID        = "test_user_uuid_1234"
)

func TestContactsUseCase_GetContacts(t *testing.T) {
	type args struct {
		ctx      context.Context
		userUuid string
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo)
		want       []entity.Contacts
		wantErr    bool
	}

	successfulGetContacts := entity.Contacts{
		UserProfile: entity.UserProfile{
			UserUUID:  testContactUserUUID,
			FirstName: "test_firstname",
			LastName:  "test_lastname",
			Avatar:    "test_avatar",
		},
		ConversationUUID: "conversation_1234",
		Blocked:          false,
	}

	tests := []testCase{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				userUuid: testUserUUID,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockRepo.EXPECT().
					GetContactsByUserUUID(gomock.Any(), testUserUUID).
					Return([]entity.Contacts{successfulGetContacts}, nil)
			},
			want:    []entity.Contacts{successfulGetContacts},
			wantErr: false,
		},
		{
			name: "error fetching contacts",
			args: args{
				ctx:      context.Background(),
				userUuid: testUserUUID,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockRepo.EXPECT().
					GetContactsByUserUUID(gomock.Any(), testUserUUID).
					Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockContactsRepo(ctrl)
			mockUserRepo := mocks.NewMockUserRepo(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo, mockUserRepo)
			}

			uc := &ContactsUseCase{
				repo:         mockRepo,
				userInfoRepo: mockUserRepo,
			}

			got, err := uc.GetContacts(tt.args.ctx, tt.args.userUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsUseCase.GetContacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContactsUseCase.GetContacts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactsUseCase_AddContact(t *testing.T) {
	type args struct {
		ctx             context.Context
		contactUserName string
		userUuid        string
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo)
		wantErr    bool
	}

	tests := []testCase{
		{
			name: "success",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil)

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(false, nil)

				mockRepo.EXPECT().
					StoreContacts(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "contacts username not found",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "contact already exists",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil)

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(true, nil)

				mockRepo.EXPECT().
					UpdateRemovedStatus(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockContactsRepo(ctrl)
			mockUserRepo := mocks.NewMockUserRepo(ctrl)

			// Set up the mocks as per the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo, mockUserRepo)
			}

			uc := &ContactsUseCase{
				repo:         mockRepo,
				userInfoRepo: mockUserRepo,
			}

			err := uc.AddContact(tt.args.ctx, tt.args.contactUserName, tt.args.userUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsUseCase.AddContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContactsUseCase_RemoveContact(t *testing.T) {
	type args struct {
		ctx             context.Context
		contactUserName string
		userUuid        string
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo)
		wantErr    bool
	}

	tests := []testCase{
		{
			name: "success",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil)

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(true, nil)

				mockRepo.EXPECT().
					UpdateRemovedStatus(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "contact does not exist",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil)

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(false, nil)
			},
			wantErr: true,
		},
		{
			name: "username not found",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockContactsRepo(ctrl)
			mockUserRepo := mocks.NewMockUserRepo(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo, mockUserRepo)
			}

			uc := &ContactsUseCase{
				repo:         mockRepo,
				userInfoRepo: mockUserRepo,
			}

			err := uc.RemoveContact(tt.args.ctx, tt.args.contactUserName, tt.args.userUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsUseCase.RemoveContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContactsUseCase_UpdateBlockContact(t *testing.T) {
	type args struct {
		ctx             context.Context
		contactUserName string
		userUuid        string
		block           bool
	}
	type testCase struct {
		name       string
		args       args
		setupMocks func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo)
		wantErr    bool
	}

	tests := []testCase{
		{
			name: "success - block contact",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
				block:           true,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil)

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(true, nil)

				mockRepo.EXPECT().
					UpdateBlockedStatus(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - unblock contact",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
				block:           false,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil)

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(true, nil)

				mockRepo.EXPECT().
					UpdateBlockedStatus(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "username not found",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
				block:           true,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "contact does not exist",
			args: args{
				ctx:             context.Background(),
				contactUserName: testContactUserName,
				userUuid:        testUserUUID,
				block:           true,
			},
			setupMocks: func(mockRepo *mocks.MockContactsRepo, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.EXPECT().
					GetUserUUIDByUsername(gomock.Any(), testContactUserName).
					Return(&testContactUserUUID, nil)

				mockRepo.EXPECT().
					CheckContactExist(gomock.Any(), testUserUUID, testContactUserUUID).
					Return(false, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockContactsRepo(ctrl)
			mockUserRepo := mocks.NewMockUserRepo(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo, mockUserRepo)
			}

			uc := &ContactsUseCase{
				repo:         mockRepo,
				userInfoRepo: mockUserRepo,
			}

			err := uc.UpdateBlockContact(tt.args.ctx, tt.args.contactUserName, tt.args.userUuid, tt.args.block)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsUseCase.UpdateBlockContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
