package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo/api/entities"
	"todo/api/enum"
	"todo/api/models/request"
	"todo/api/models/response"
	"todo/api/services/mock"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_taskHandler_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		taskService         *mock.MockTaskService
		taskServiceBehavior func(*mock.MockTaskService)
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{
		{
			name: "success",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
					mts.EXPECT().CreateTask(gomock.Any()).Return(nil)
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.CreatedTaskRequest{
						Title:       "foo",
						Description: "foo",
						Image:       "foo",
						Status:      enum.TaskStatusCompleted,
					})

					return httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusOK,
		},
		{
			name: "body parser failed",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
				},
			},
			args: args{
				req: func() *http.Request {
					return httptest.NewRequest("POST", "/api/tasks", nil)
				}(),
			},
			code: fiber.StatusBadRequest,
		},
		{
			name: "title long length failed",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.CreatedTaskRequest{
						Title:       "fooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo",
						Description: "foo",
						Image:       "foo",
						Status:      enum.TaskStatusCompleted,
					})

					return httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusBadRequest,
		},
		{
			name: "status is invalid",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.CreatedTaskRequest{
						Title:       "foo",
						Description: "foo",
						Image:       "foo",
						Status:      "foo",
					})

					return httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusBadRequest,
		},
		{
			name: "create task failed",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
					mts.EXPECT().CreateTask(gomock.Any()).Return(errors.New("foo"))
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.CreatedTaskRequest{
						Title:       "foo",
						Description: "foo",
						Image:       "foo",
						Status:      enum.TaskStatusCompleted,
					})

					return httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.taskServiceBehavior(tt.fields.taskService)

			app := fiber.New()
			h := taskHandler{
				taskService: tt.fields.taskService,
			}
			app.Post("/api/tasks", h.CreateTask)

			tt.args.req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(tt.args.req)
			if err != nil {
				t.Fatalf("Error while performing the request: %v", err)
			}

			var result response.Response
			json.NewDecoder(resp.Body).Decode(&result)
			assert.Equal(t, tt.code, resp.StatusCode)
		})
	}
}

func Test_taskHandler_GetTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		taskService         *mock.MockTaskService
		taskServiceBehavior func(*mock.MockTaskService)
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{
		{
			name: "success",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
					mts.EXPECT().GetTasks(gomock.Any()).Return([]entities.Task{}, nil)
				},
			},
			args: args{
				req: func() *http.Request {
					return httptest.NewRequest("GET", "/api/tasks?title=foo&sort_by=title&sort_order=asc", nil)
				}(),
			},
			code: fiber.StatusOK,
		},
		{
			name: "validate failed",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
				},
			},
			args: args{
				req: func() *http.Request {
					return httptest.NewRequest("GET", "/api/tasks?title=foo&sort_by=foo&sort_order=asc", nil)
				}(),
			},
			code: fiber.StatusBadRequest,
		},
		{
			name: "get tasks failed",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
					mts.EXPECT().GetTasks(gomock.Any()).Return([]entities.Task{}, errors.New("foo"))
				},
			},
			args: args{
				req: func() *http.Request {
					return httptest.NewRequest("GET", "/api/tasks?title=foo&sort_by=title&sort_order=asc", nil)
				}(),
			},
			code: fiber.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.taskServiceBehavior(tt.fields.taskService)

			app := fiber.New()
			h := taskHandler{
				taskService: tt.fields.taskService,
			}
			app.Get("/api/tasks", h.GetTasks)

			tt.args.req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(tt.args.req)
			if err != nil {
				t.Fatalf("Error while performing the request: %v", err)
			}

			var result response.Response
			json.NewDecoder(resp.Body).Decode(&result)
			assert.Equal(t, tt.code, resp.StatusCode)
		})
	}
}

func Test_taskHandler_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		taskService         *mock.MockTaskService
		taskServiceBehavior func(*mock.MockTaskService)
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{
		{
			name: "success",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
					mts.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.UpdatedTaskRequest{
						Title:       "foo",
						Description: "foo",
						Image:       "foo",
						Status:      enum.TaskStatusCompleted,
					})

					return httptest.NewRequest("PUT", "/api/tasks/1", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusOK,
		},
		{
			name: "body parser failed",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
				},
			},
			args: args{
				req: func() *http.Request {
					return httptest.NewRequest("PUT", "/api/tasks/1", nil)
				}(),
			},
			code: fiber.StatusBadRequest,
		},
		{
			name: "title long length failed",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.UpdatedTaskRequest{
						Title:       "fooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo",
						Description: "foo",
						Image:       "foo",
						Status:      enum.TaskStatusCompleted,
					})

					return httptest.NewRequest("PUT", "/api/tasks/1", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusBadRequest,
		},
		{
			name: "status is invalid",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.UpdatedTaskRequest{
						Title:       "foo",
						Description: "foo",
						Image:       "foo",
						Status:      "foo",
					})

					return httptest.NewRequest("PUT", "/api/tasks/1", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusBadRequest,
		},
		{
			name: "id is not int",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.UpdatedTaskRequest{
						Title:       "foo",
						Description: "foo",
						Image:       "foo",
						Status:      enum.TaskStatusCompleted,
					})

					return httptest.NewRequest("PUT", "/api/tasks/foo", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusBadRequest,
		},
		{
			name: "update task failed",
			fields: fields{
				taskService: mock.NewMockTaskService(ctrl),
				taskServiceBehavior: func(mts *mock.MockTaskService) {
					mts.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(errors.New("foo"))
				},
			},
			args: args{
				req: func() *http.Request {
					b, _ := json.Marshal(request.UpdatedTaskRequest{
						Title:       "foo",
						Description: "foo",
						Image:       "foo",
						Status:      enum.TaskStatusCompleted,
					})

					return httptest.NewRequest("PUT", "/api/tasks/1", bytes.NewReader(b))
				}(),
			},
			code: fiber.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.taskServiceBehavior(tt.fields.taskService)

			app := fiber.New()
			h := taskHandler{
				taskService: tt.fields.taskService,
			}
			app.Put("/api/tasks/:id", h.UpdateTask)

			tt.args.req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(tt.args.req)
			if err != nil {
				t.Fatalf("Error while performing the request: %v", err)
			}

			var result response.Response
			json.NewDecoder(resp.Body).Decode(&result)
			assert.Equal(t, tt.code, resp.StatusCode)
		})
	}
}
