package types

import (
	"testing"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
)

func TestCreatedUpdated_SetCreatedBy(t *testing.T) {
	var basicUser = user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	_ = basicUser.SetUUID("cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0")

	type fields struct {
		createdInfo createdInfo
		updatedInfo updatedInfo
	}
	type args struct {
		basicUser user.BasicUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "createdBy not set",
			fields:  fields{},
			args:    args{basicUser: basicUser},
			wantErr: false,
		},
		{
			name: "createdBy was already set => returns error",
			fields: fields{createdInfo: createdInfo{
				createdBy: basicUser,
			}},
			args:    args{basicUser: basicUser},
			wantErr: true,
		},
		{
			name: "createdAt was already set",
			fields: fields{createdInfo: createdInfo{
				createdAt: "some time string",
			}},
			args:    args{basicUser: basicUser},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreatedUpdated{
				createdInfo: tt.fields.createdInfo,
				updatedInfo: tt.fields.updatedInfo,
			}
			if err := o.SetCreatedBy(tt.args.basicUser); (err != nil) != tt.wantErr {
				t.Errorf("SetCreatedBy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatedUpdated_SetCreated(t *testing.T) {
	var basicUser = user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	_ = basicUser.SetUUID("cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0")

	type fields struct {
		createdInfo createdInfo
		updatedInfo updatedInfo
	}
	type args struct {
		basicUser user.BasicUser
		dateTime  DateTime
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "both createdBy and createdAt not set",
			fields:  fields{},
			args:    args{basicUser: basicUser, dateTime: "actual time"},
			wantErr: false,
		},
		{
			name: "createdAt was already set => returns error",
			fields: fields{createdInfo: createdInfo{
				createdAt: "some time string",
			}},
			args:    args{dateTime: "actual time"},
			wantErr: true,
		},
		{
			name: "createdBy was already set => returns error",
			fields: fields{createdInfo: createdInfo{
				createdBy: basicUser,
			}},
			args:    args{basicUser: basicUser, dateTime: "actual time"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreatedUpdated{
				createdInfo: tt.fields.createdInfo,
				updatedInfo: tt.fields.updatedInfo,
			}
			if err := o.SetCreated(tt.args.basicUser, tt.args.dateTime); (err != nil) != tt.wantErr {
				t.Errorf("SetCreatedAt error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatedUpdated_SetUpdatedBy(t *testing.T) {
	var basicUser = user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	_ = basicUser.SetUUID("cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0")

	type fields struct {
		createdInfo createdInfo
		updatedInfo updatedInfo
	}
	type args struct {
		basicUser user.BasicUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "updatedBy not set",
			fields:  fields{},
			args:    args{basicUser: basicUser},
			wantErr: false,
		},
		{
			name: "updatedBy was already set",
			fields: fields{updatedInfo: updatedInfo{
				updatedBy: basicUser,
			}},
			args:    args{basicUser: basicUser},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreatedUpdated{
				createdInfo: tt.fields.createdInfo,
				updatedInfo: tt.fields.updatedInfo,
			}
			if err := o.SetUpdatedBy(tt.args.basicUser); (err != nil) != tt.wantErr {
				t.Errorf("SetUpdatedBy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatedUpdated_SetUpdated(t *testing.T) {
	var basicUser = user.BasicUser{
		ExternalUserUUID: "b306a60e-a2a5-463f-a6e1-33e8cb21bc3b",
		Name:             "Alfred",
		Surname:          "Koletschko",
		OrgDisplayName:   "KompiTech",
		OrgName:          "a897a407-e41b-4b14-924a-39f5d5a8038f.kompitech.com",
	}
	_ = basicUser.SetUUID("cb2fe2a7-ab9f-4f6d-9fd6-c7c209403cf0")

	type fields struct {
		createdInfo createdInfo
		updatedInfo updatedInfo
	}
	type args struct {
		basicUser user.BasicUser
		dateTime  DateTime
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "both updatedBy and updatedAt not set",
			fields:  fields{},
			args:    args{basicUser: basicUser, dateTime: "actual time"},
			wantErr: false,
		},
		{
			name: "updatedBy was already set",
			fields: fields{updatedInfo: updatedInfo{
				updatedBy: basicUser,
			}},
			args:    args{basicUser: basicUser, dateTime: "actual time"},
			wantErr: false,
		},
		{
			name: "updatedAt was already set",
			fields: fields{updatedInfo: updatedInfo{
				updatedAt: "some time string",
			}},
			args:    args{dateTime: "actual time"},
			wantErr: false,
		},
		{
			name: "createdBy was already set",
			fields: fields{createdInfo: createdInfo{
				createdBy: basicUser,
			}},
			args:    args{basicUser: basicUser, dateTime: "actual time"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreatedUpdated{
				createdInfo: tt.fields.createdInfo,
				updatedInfo: tt.fields.updatedInfo,
			}
			if err := o.SetUpdated(tt.args.basicUser, tt.args.dateTime); (err != nil) != tt.wantErr {
				t.Errorf("SetUpdatedBy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
