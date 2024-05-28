package services

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"
	"todo/api/entities"
	"todo/api/enum"
	"todo/api/models/request"
	"todo/pkg/base"
	"todo/pkg/base/mock"
	"todo/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func Test_taskService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		repository         *mock.MockBaseRepository[any]
		repositoryBehavior func(*mock.MockBaseRepository[any])
	}
	type args struct {
		req request.CreatedTaskRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: mock.NewMockBaseRepository[any](ctrl),
				repositoryBehavior: func(mbr *mock.MockBaseRepository[any]) {
					mbr.EXPECT().Create(gomock.Any()).Return(mbr)
					mbr.EXPECT().Error()
				},
			},
			args: args{
				req: request.CreatedTaskRequest{
					Title:       "foo",
					Description: "foo",
					Image:       "foo",
					Status:      "COMPLETED",
				},
			},
			wantErr: false,
		},
		{
			name: "create failed",
			fields: fields{
				repository: mock.NewMockBaseRepository[any](ctrl),
				repositoryBehavior: func(mbr *mock.MockBaseRepository[any]) {
					mbr.EXPECT().Create(gomock.Any()).Return(mbr)
					mbr.EXPECT().Error().Return(errors.New("foo"))
				},
			},
			args: args{
				req: request.CreatedTaskRequest{
					Title:       "foo",
					Description: "foo",
					Image:       "foo",
					Status:      "COMPLETED",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.repositoryBehavior(tt.fields.repository)
			s := taskService{
				repository: tt.fields.repository,
				log:        logger.WithPrefix("test"),
			}

			if err := s.CreateTask(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("taskService.CreateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskService_GetTasks(t *testing.T) {
	var (
		tn    = time.Now()
		tasks = []entities.Task{{
			ID:          1,
			Title:       "foo",
			Description: "foo",
			Image:       "foo",
			Status:      enum.TaskStatusCompleted,
			CreatedAt:   tn,
			UpdatedAt:   tn,
		}}
	)

	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	type fields struct {
		repository         base.BaseRepository[any]
		repositoryBehavior func()
	}
	type args struct {
		query request.TaskListQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entities.Task
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: base.NewBaseRepository[any](db),
				repositoryBehavior: func() {
					users := sqlmock.NewRows([]string{"id", "title", "description", "image", "status", "created_at", "updated_at"}).
						AddRow(1, "foo", "foo", "foo", "COMPLETED", tn, tn)
					expectedSQL := `SELECT (.+) FROM "tasks" WHERE title LIKE (.+) AND description LIKE (.+)`
					mock.ExpectQuery(expectedSQL).WillReturnRows(users)
				},
			},
			args: args{
				query: request.TaskListQuery{
					Title:       "foo",
					Description: "foo",
					SortBy:      "title",
					SortOrder:   "asc",
				},
			},
			want:    tasks,
			wantErr: false,
		},
		{
			name: "find tasks failed",
			fields: fields{
				repository: base.NewBaseRepository[any](db),
				repositoryBehavior: func() {
					users := sqlmock.NewRows([]string{"asd"}).
						AddRow(1)
					expectedSQL := `SELECT (.+) FROM "tasks" WHERE title LIKE (.+) AND description LIKE (.+)`
					mock.ExpectQuery(expectedSQL).WillReturnRows(users)
				},
			},
			args: args{
				query: request.TaskListQuery{
					Title:       "foo",
					Description: "foo",
					SortBy:      "title",
					SortOrder:   "asc",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.repositoryBehavior()
			s := taskService{
				repository: tt.fields.repository,
				log:        logger.WithPrefix("test"),
			}

			got, err := s.GetTasks(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskService.GetTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskService.GetTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskService_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		repository         *mock.MockBaseRepository[any]
		repositoryBehavior func(*mock.MockBaseRepository[any])
	}
	type args struct {
		id  int
		req request.UpdatedTaskRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repository: mock.NewMockBaseRepository[any](ctrl),
				repositoryBehavior: func(mbr *mock.MockBaseRepository[any]) {
					mbr.EXPECT().Model(gomock.Any()).Return(mbr)
					mbr.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mbr)
					mbr.EXPECT().Updates(gomock.Any()).Return(mbr)
					mbr.EXPECT().Error().Return(nil)
				},
			},
			args: args{
				id: 1,
				req: request.UpdatedTaskRequest{
					Title:       "foo",
					Description: "foo",
					Image:       "foo",
					Status:      enum.TaskStatusCompleted,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				repository: mock.NewMockBaseRepository[any](ctrl),
				repositoryBehavior: func(mbr *mock.MockBaseRepository[any]) {
					mbr.EXPECT().Model(gomock.Any()).Return(mbr)
					mbr.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mbr)
					mbr.EXPECT().Updates(gomock.Any()).Return(mbr)
					mbr.EXPECT().Error().Return(errors.New("foo"))
				},
			},
			args: args{
				id: 1,
				req: request.UpdatedTaskRequest{
					Title:       "foo",
					Description: "foo",
					Image:       "foo",
					Status:      enum.TaskStatusCompleted,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.repositoryBehavior(tt.fields.repository)
			s := taskService{
				repository: tt.fields.repository,
				log:        logger.WithPrefix("test"),
			}
			if err := s.UpdateTask(tt.args.id, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("taskService.UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
