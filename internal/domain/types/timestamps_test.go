package types

import (
	"testing"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
)

func TestCreatedUpdated_SetCreatedBy(t *testing.T) {
	type fields struct {
		createdInfo createdInfo
		updatedInfo updatedInfo
	}
	type args struct {
		userID ref.UUID
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
			args:    args{userID: "seme user UUID"},
			wantErr: false,
		},
		{
			name: "createdBy was already set => returns error",
			fields: fields{createdInfo: createdInfo{
				createdBy: "5a4b6317-f99c-4c21-aa82-9ca5671d7f18",
			}},
			args:    args{userID: "seme user UUID"},
			wantErr: true,
		},
		{
			name: "createdAt was already set",
			fields: fields{createdInfo: createdInfo{
				createdAt: "some time string",
			}},
			args:    args{userID: "seme user UUID"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreatedUpdated{
				createdInfo: tt.fields.createdInfo,
				updatedInfo: tt.fields.updatedInfo,
			}
			if err := o.SetCreatedBy(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("SetCreatedBy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatedUpdated_SetCreated(t *testing.T) {
	type fields struct {
		createdInfo createdInfo
		updatedInfo updatedInfo
	}
	type args struct {
		userID   ref.UUID
		dateTime DateTime
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
			args:    args{userID: "864e20a3-51ae-49fa-ae74-27768f8d2d48", dateTime: "actual time"},
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
				createdBy: "5a4b6317-f99c-4c21-aa82-9ca5671d7f18",
			}},
			args:    args{userID: "864e20a3-51ae-49fa-ae74-27768f8d2d48", dateTime: "actual time"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreatedUpdated{
				createdInfo: tt.fields.createdInfo,
				updatedInfo: tt.fields.updatedInfo,
			}
			if err := o.SetCreated(tt.args.userID, tt.args.dateTime); (err != nil) != tt.wantErr {
				t.Errorf("SetCreatedAt error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatedUpdated_SetUpdatedBy(t *testing.T) {
	type fields struct {
		createdInfo createdInfo
		updatedInfo updatedInfo
	}
	type args struct {
		userID ref.UUID
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
			args:    args{userID: "seme user UUID"},
			wantErr: false,
		},
		{
			name: "updatedBy was already set",
			fields: fields{updatedInfo: updatedInfo{
				updatedBy: "5a4b6317-f99c-4c21-aa82-9ca5671d7f18",
			}},
			args:    args{userID: "seme user UUID"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreatedUpdated{
				createdInfo: tt.fields.createdInfo,
				updatedInfo: tt.fields.updatedInfo,
			}
			if err := o.SetUpdatedBy(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("SetUpdatedBy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatedUpdated_SetUpdated(t *testing.T) {
	type fields struct {
		createdInfo createdInfo
		updatedInfo updatedInfo
	}
	type args struct {
		userID   ref.UUID
		dateTime DateTime
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
			args:    args{userID: "864e20a3-51ae-49fa-ae74-27768f8d2d48", dateTime: "actual time"},
			wantErr: false,
		},
		{
			name: "updatedBy was already set",
			fields: fields{updatedInfo: updatedInfo{
				updatedBy: "5a4b6317-f99c-4c21-aa82-9ca5671d7f18",
			}},
			args:    args{userID: "864e20a3-51ae-49fa-ae74-27768f8d2d48", dateTime: "actual time"},
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
				createdBy: "5a4b6317-f99c-4c21-aa82-9ca5671d7f18",
			}},
			args:    args{userID: "864e20a3-51ae-49fa-ae74-27768f8d2d48", dateTime: "actual time"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreatedUpdated{
				createdInfo: tt.fields.createdInfo,
				updatedInfo: tt.fields.updatedInfo,
			}
			if err := o.SetUpdated(tt.args.userID, tt.args.dateTime); (err != nil) != tt.wantErr {
				t.Errorf("SetUpdatedBy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
