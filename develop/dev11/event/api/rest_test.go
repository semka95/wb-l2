package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"calendar/event"
	"calendar/event/repository/bolt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var eventTime = time.Date(2022, 7, 5, 15, 4, 1, 0, time.UTC)

var tEvent = event.Event{
	ID:    1,
	Date:  eventTime,
	Title: "birthday",
}

var tEventForDay = []event.Event{
	{
		ID:    1,
		Date:  time.Date(2022, 7, 5, 15, 4, 1, 0, time.UTC),
		Title: "test",
	},
	{
		ID:    1,
		Date:  time.Date(2022, 7, 5, 21, 12, 37, 0, time.UTC),
		Title: "123",
	},
}

var tEventForWeek = []event.Event{
	{
		ID:    1,
		Date:  time.Date(2022, 7, 6, 15, 4, 1, 0, time.UTC),
		Title: "test",
	},
	{
		ID:    1,
		Date:  time.Date(2022, 7, 5, 21, 12, 37, 0, time.UTC),
		Title: "123",
	},
}

var tEventForMonth = []event.Event{
	{
		ID:    1,
		Date:  time.Date(2022, 7, 21, 15, 4, 1, 0, time.UTC),
		Title: "test",
	},
	{
		ID:    1,
		Date:  time.Date(2022, 7, 5, 21, 12, 37, 0, time.UTC),
		Title: "123",
	},
}

type jsonError struct {
	Details string `json:"details,omitempty"`
	Error   string `json:"error,omitempty"`
}

func TestCreate(t *testing.T) {
	api := API{}
	req := new(http.Request)

	testCases := []struct {
		desc           string
		store          *bolt.EventRepositoryMock
		reqBody        string
		checkMockCalls func(tr *bolt.EventRepositoryMock)
		checkResponse  func(rec *httptest.ResponseRecorder)
	}{
		{
			desc: "success",
			store: &bolt.EventRepositoryMock{
				CreateFunc: func(user_id uint64, e event.Event) (event.Event, error) {
					tr := event.Event{
						ID:    1,
						Title: e.Title,
						Date:  e.Date,
					}
					return tr, nil
				},
			},
			reqBody: "user_id=3&date=2022-07-05T15:04:01Z&title=birthday",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.CreateCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				got := event.Event{}
				err := json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				assert.EqualValues(t, tEvent, got)
				assert.Equal(t, http.StatusCreated, rec.Code)
			},
		},
		{
			desc:           "bad user_id",
			store:          &bolt.EventRepositoryMock{},
			reqBody:        "user_id=bad data",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse user_id", jsonErr.Details)
				assert.EqualValues(t, "strconv.Atoi: parsing \"bad data\": invalid syntax", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc:           "bad date",
			store:          &bolt.EventRepositoryMock{},
			reqBody:        "user_id=3&date=bad date",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse date, use RFC3339 format", jsonErr.Details)
				assert.EqualValues(t, "parsing time \"bad date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"bad date\" as \"2006\"", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc:           "empty title",
			store:          &bolt.EventRepositoryMock{},
			reqBody:        "user_id=3&date=2022-07-05T15:04:01Z",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "no title provided", jsonErr.Details)
				assert.EqualValues(t, "empty title", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "store server error",
			store: &bolt.EventRepositoryMock{
				CreateFunc: func(user_id uint64, e event.Event) (event.Event, error) {
					return event.Event{}, fmt.Errorf("can't create record")
				},
			},
			reqBody: "user_id=3&date=2022-07-05T15:04:01Z&title=birthday",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.CreateCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't create event", jsonErr.Details)
				assert.EqualValues(t, "can't create record", jsonErr.Error)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			api.eventStore = tC.store

			req = httptest.NewRequest("POST", "/create_event", strings.NewReader(tC.reqBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rec := httptest.NewRecorder()
			api.Create(rec, req)

			tC.checkMockCalls(tC.store)

			tC.checkResponse(rec)
		})
	}
}

func TestUpdate(t *testing.T) {
	api := API{}
	req := new(http.Request)

	testCases := []struct {
		desc           string
		store          *bolt.EventRepositoryMock
		reqBody        string
		checkMockCalls func(tr *bolt.EventRepositoryMock)
		checkResponse  func(rec *httptest.ResponseRecorder)
	}{
		{
			desc: "success",
			store: &bolt.EventRepositoryMock{
				UpdateFunc: func(user_id uint64, e event.Event) error {
					return nil
				},
			},
			reqBody: "user_id=3&id=1&date=2022-07-05T15:04:01Z&title=birthday",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.UpdateCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			desc:           "bad user_id",
			store:          &bolt.EventRepositoryMock{},
			reqBody:        "user_id=bad data",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse user_id", jsonErr.Details)
				assert.EqualValues(t, "strconv.Atoi: parsing \"bad data\": invalid syntax", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc:           "bad id",
			store:          &bolt.EventRepositoryMock{},
			reqBody:        "user_id=3",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse id", jsonErr.Details)
				assert.EqualValues(t, "strconv.Atoi: parsing \"\": invalid syntax", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc:           "bad date",
			store:          &bolt.EventRepositoryMock{},
			reqBody:        "user_id=3&id=1&date=bad date",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse date, use RFC3339 format", jsonErr.Details)
				assert.EqualValues(t, "parsing time \"bad date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"bad date\" as \"2006\"", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc:           "empty title",
			store:          &bolt.EventRepositoryMock{},
			reqBody:        "user_id=3&id=1&date=2022-07-05T15:04:01Z",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "no title provided", jsonErr.Details)
				assert.EqualValues(t, "empty title", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "store server error",
			store: &bolt.EventRepositoryMock{
				UpdateFunc: func(user_id uint64, e event.Event) error {
					return fmt.Errorf("can't update record")
				},
			},
			reqBody: "user_id=3&id=1&date=2022-07-05T15:04:01Z&title=birthday",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.UpdateCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't update event", jsonErr.Details)
				assert.EqualValues(t, "can't update record", jsonErr.Error)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "event not found",
			store: &bolt.EventRepositoryMock{
				UpdateFunc: func(user_id uint64, e event.Event) error {
					return event.ErrNotFound
				},
			},
			reqBody: "user_id=3&id=1&date=2022-07-05T15:04:01Z&title=birthday",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.UpdateCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't update event", jsonErr.Details)
				assert.EqualValues(t, "your requested item is not found", jsonErr.Error)
				assert.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			api.eventStore = tC.store

			req = httptest.NewRequest("PUT", "/update_event", strings.NewReader(tC.reqBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rec := httptest.NewRecorder()
			api.Update(rec, req)

			tC.checkMockCalls(tC.store)

			tC.checkResponse(rec)
		})
	}
}

func TestDelete(t *testing.T) {
	api := API{}
	req := new(http.Request)

	testCases := []struct {
		desc           string
		store          *bolt.EventRepositoryMock
		user_id        string
		id             string
		checkMockCalls func(tr *bolt.EventRepositoryMock)
		checkResponse  func(rec *httptest.ResponseRecorder)
	}{
		{
			desc: "success",
			store: &bolt.EventRepositoryMock{
				DeleteFunc: func(user_id uint64, event_id uint64) error {
					return nil
				},
			},
			user_id: "3",
			id:      "1",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.DeleteCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			desc:           "bad user_id",
			store:          &bolt.EventRepositoryMock{},
			user_id:        "bad data",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse user_id", jsonErr.Details)
				assert.EqualValues(t, "strconv.Atoi: parsing \"bad data\": invalid syntax", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc:           "bad id",
			store:          &bolt.EventRepositoryMock{},
			user_id:        "3",
			id:             "bad data",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse id", jsonErr.Details)
				assert.EqualValues(t, "strconv.Atoi: parsing \"bad data\": invalid syntax", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "store server error",
			store: &bolt.EventRepositoryMock{
				DeleteFunc: func(user_id uint64, event_id uint64) error {
					return fmt.Errorf("can't delete record")
				},
			},
			user_id: "3",
			id:      "1",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.DeleteCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't delete event", jsonErr.Details)
				assert.EqualValues(t, "can't delete record", jsonErr.Error)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "event not found",
			store: &bolt.EventRepositoryMock{
				DeleteFunc: func(user_id uint64, event_id uint64) error {
					return event.ErrNotFound
				},
			},
			user_id: "3",
			id:      "1",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.DeleteCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't delete event", jsonErr.Details)
				assert.EqualValues(t, "your requested item is not found", jsonErr.Error)
				assert.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			api.eventStore = tC.store

			req = httptest.NewRequest("DELETE", "/delete_event", nil)
			q := req.URL.Query()
			q.Add("user_id", tC.user_id)
			q.Add("id", tC.id)
			req.URL.RawQuery = q.Encode()
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rec := httptest.NewRecorder()
			api.Delete(rec, req)

			tC.checkMockCalls(tC.store)

			tC.checkResponse(rec)
		})
	}
}

func TestGet(t *testing.T) {
	api := API{}
	req := new(http.Request)

	testCases := []struct {
		desc           string
		store          *bolt.EventRepositoryMock
		path           string
		user_id        string
		date           string
		checkMockCalls func(tr *bolt.EventRepositoryMock)
		checkResponse  func(rec *httptest.ResponseRecorder)
	}{
		{
			desc: "success events for day",
			store: &bolt.EventRepositoryMock{
				GetForDayFunc: func(user_id uint64, day time.Time) ([]event.Event, error) {
					return tEventForDay, nil
				},
			},
			path:    "/events_for_day",
			user_id: "3",
			date:    "2022-07-05T15:04:01Z",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.GetForDayCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				got := []event.Event{}
				err := json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				assert.EqualValues(t, tEventForDay, got)
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			desc: "success events for week",
			store: &bolt.EventRepositoryMock{
				GetForWeekFunc: func(user_id uint64, day time.Time) ([]event.Event, error) {
					return tEventForWeek, nil
				},
			},
			path:    "/events_for_week",
			user_id: "3",
			date:    "2022-07-05T15:04:01Z",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.GetForWeekCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				got := []event.Event{}
				err := json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				assert.EqualValues(t, tEventForWeek, got)
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			desc: "success events for month",
			store: &bolt.EventRepositoryMock{
				GetForMonthFunc: func(user_id uint64, day time.Time) ([]event.Event, error) {
					return tEventForMonth, nil
				},
			},
			path:    "/events_for_month",
			user_id: "3",
			date:    "2022-07-05T15:04:01Z",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.GetForMonthCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				got := []event.Event{}
				err := json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				assert.EqualValues(t, tEventForMonth, got)
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			desc:           "bad user_id",
			store:          &bolt.EventRepositoryMock{},
			path:           "/events_for_day",
			user_id:        "bad data",
			date:           "2022-07-05T15:04:01Z",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse user_id", jsonErr.Details)
				assert.EqualValues(t, "strconv.Atoi: parsing \"bad data\": invalid syntax", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc:           "bad date",
			store:          &bolt.EventRepositoryMock{},
			path:           "/events_for_day",
			user_id:        "3",
			date:           "bad date",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't parse date, use RFC3339 format", jsonErr.Details)
				assert.EqualValues(t, "parsing time \"bad date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"bad date\" as \"2006\"", jsonErr.Error)
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			desc: "server error",
			store: &bolt.EventRepositoryMock{
				GetForDayFunc: func(user_id uint64, day time.Time) ([]event.Event, error) {
					return []event.Event{}, fmt.Errorf("can't get events")
				},
			},
			path:    "/events_for_day",
			user_id: "3",
			date:    "2022-07-05T15:04:01Z",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.GetForDayCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				jsonErr := new(jsonError)
				err := json.NewDecoder(rec.Body).Decode(&jsonErr)
				require.NoError(t, err)
				assert.EqualValues(t, "can't get events", jsonErr.Details)
				assert.EqualValues(t, "can't get events", jsonErr.Error)
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			desc: "no events",
			store: &bolt.EventRepositoryMock{
				GetForDayFunc: func(user_id uint64, day time.Time) ([]event.Event, error) {
					return []event.Event{}, nil
				},
			},
			path:    "/events_for_day",
			user_id: "3",
			date:    "2022-07-05T15:04:01Z",
			checkMockCalls: func(tr *bolt.EventRepositoryMock) {
				calls := len(tr.GetForDayCalls())
				assert.Equal(t, 1, calls)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			api.eventStore = tC.store

			req = httptest.NewRequest("GET", tC.path, nil)
			q := req.URL.Query()
			q.Add("user_id", tC.user_id)
			q.Add("date", tC.date)
			req.URL.RawQuery = q.Encode()
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rec := httptest.NewRecorder()
			api.Get(rec, req)

			tC.checkMockCalls(tC.store)

			tC.checkResponse(rec)
		})
	}
}
