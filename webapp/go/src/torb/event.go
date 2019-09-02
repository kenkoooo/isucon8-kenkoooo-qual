package main

import (
	"database/sql"
	"fmt"
	"sort"
)

func getEvent(eventID, loginUserID int64) (*Event, error) {
	var event Event
	if err := db.QueryRow("SELECT * FROM events WHERE id = ?", eventID).Scan(&event.ID, &event.Title, &event.PublicFg, &event.ClosedFg, &event.Price); err != nil {
		return nil, err
	}

	query := `
	SELECT s.*, r.user_id, r.reserved_at
	FROM sheets s
	LEFT JOIN (
		SELECT sheet_id, user_id, reserved_at 
		FROM reservations 
		WHERE event_id = ? AND canceled_at IS NULL 
		GROUP BY event_id, sheet_id 
		HAVING reserved_at = MIN(reserved_at)
	) r ON r.sheet_id = s.id
	ORDER BY s.rank, s.num
	`

	rows, err := db.Query(query, event.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sheets []*Sheet
	var reservations []*Reservation
	var userIDs []*sql.NullInt64
	for rows.Next() {
		var s Sheet
		var r Reservation
		var params struct {
			UserID sql.NullInt64
		}

		if err := rows.Scan(
			&s.ID, &s.Rank, &s.Num, &s.Price, &params.UserID, &r.ReservedAt); err != nil {
			return nil, err
		}
		sheets = append(sheets, &s)
		reservations = append(reservations, &r)
		userIDs = append(userIDs, &params.UserID)
	}

	RefineEvent(reservations, sheets, userIDs, &event, loginUserID)

	return &event, nil
}

func RefineEvent(reservations []*Reservation, sheets []*Sheet, userIDs []*sql.NullInt64, event *Event, loginUserID int64) {
	event.Sheets = map[string]*Sheets{
		"S": &Sheets{},
		"A": &Sheets{},
		"B": &Sheets{},
		"C": &Sheets{},
	}
	for i, r := range reservations {
		s := sheets[i]
		userID := userIDs[i]
		event.Sheets[s.Rank].Price = event.Price + s.Price
		event.Total++
		event.Sheets[s.Rank].Total++

		if r.ReservedAt != nil {
			s.Mine = userID.Valid && userID.Int64 == loginUserID
			s.Reserved = true
			s.ReservedAtUnix = r.ReservedAt.Unix()
		} else {
			event.Remains++
			event.Sheets[s.Rank].Remains++
		}

		event.Sheets[s.Rank].Detail = append(event.Sheets[s.Rank].Detail, s)
	}
}

func getEvents(all bool) ([]*Event, error) {
	query := `
	SELECT s.*, r.user_id, r.reserved_at, e.id, e.title, e.public_fg, e.closed_fg, e.price
	FROM events e
	LEFT JOIN (
		SELECT sheet_id, user_id, reserved_at, event_id
		FROM reservations 
		WHERE canceled_at IS NULL 
		GROUP BY event_id, sheet_id 
		HAVING reserved_at = MIN(reserved_at)
	) r ON r.event_id = e.id
	LEFT JOIN sheets s ON s.id = r.sheet_id
	ORDER BY s.rank, s.num
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[int64]Event)
	sheets := make(map[int64][]*Sheet)
	reservations := make(map[int64][]*Reservation)
	userIDs := make(map[int64][]*sql.NullInt64)
	for rows.Next() {
		var e Event
		var r Reservation
		var userID sql.NullInt64
		var fakeS struct {
			ID    sql.NullInt64
			Rank  sql.NullString
			Num   sql.NullInt64
			Price sql.NullInt64
		}

		if err := rows.Scan(
			&fakeS.ID, &fakeS.Rank, &fakeS.Num, &fakeS.Price,
			&userID, &r.ReservedAt,
			&e.ID, &e.Title, &e.PublicFg, &e.ClosedFg, &e.Price); err != nil {
			return nil, err
		}
		if !all && !e.PublicFg {
			continue
		}

		eventMap[e.ID] = e
		if fakeS.ID.Valid {
			s := Sheet{
				ID:    fakeS.ID.Int64,
				Rank:  fakeS.Rank.String,
				Num:   fakeS.Num.Int64,
				Price: fakeS.Price.Int64,
			}

			sheets[e.ID] = append(sheets[e.ID], &s)
			reservations[e.ID] = append(reservations[e.ID], &r)
			userIDs[e.ID] = append(userIDs[e.ID], &userID)
		}
	}

	fmt.Println(eventMap)
	var events Events
	for _, event := range eventMap {
		s := sheets[event.ID]
		r := reservations[event.ID]
		userID := userIDs[event.ID]
		RefineEvent(r, s, userID, &event, -1)
		for k := range event.Sheets {
			event.Sheets[k].Detail = nil
		}
		events = append(events, &event)
	}

	//for i, v := range events {
	//	event, err := getEvent(v.ID, -1)
	//	if err != nil {
	//		return nil, err
	//	}
	//	for k := range event.Sheets {
	//		event.Sheets[k].Detail = nil
	//	}
	//	events[i] = event
	//}
	sort.Sort(events)
	return events, nil
}

type Events []*Event

func (e Events) Len() int {
	return len(e)
}

func (e Events) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Events) Less(i, j int) bool {
	return e[i].ID < e[j].ID
}
