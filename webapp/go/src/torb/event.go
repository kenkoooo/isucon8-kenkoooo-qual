package main

import "database/sql"

func getEvent(eventID, loginUserID int64) (*Event, error) {
	var event Event
	if err := db.QueryRow("SELECT * FROM events WHERE id = ?", eventID).Scan(&event.ID, &event.Title, &event.PublicFg, &event.ClosedFg, &event.Price); err != nil {
		return nil, err
	}
	event.Sheets = map[string]*Sheets{
		"S": &Sheets{},
		"A": &Sheets{},
		"B": &Sheets{},
		"C": &Sheets{},
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
		event.Sheets[s.Rank].Price = event.Price + s.Price
		event.Total++
		event.Sheets[s.Rank].Total++

		if r.ReservedAt != nil {
			s.Mine = params.UserID.Valid && params.UserID.Int64 == loginUserID
			s.Reserved = true
			s.ReservedAtUnix = r.ReservedAt.Unix()
		} else {
			event.Remains++
			event.Sheets[s.Rank].Remains++
		}

		event.Sheets[s.Rank].Detail = append(event.Sheets[s.Rank].Detail, &s)
	}

	return &event, nil
}

func getEvents(all bool) ([]*Event, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT * FROM events ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.Title, &event.PublicFg, &event.ClosedFg, &event.Price); err != nil {
			return nil, err
		}
		if !all && !event.PublicFg {
			continue
		}
		events = append(events, &event)
	}
	for i, v := range events {
		event, err := getEvent(v.ID, -1)
		if err != nil {
			return nil, err
		}
		for k := range event.Sheets {
			event.Sheets[k].Detail = nil
		}
		events[i] = event
	}
	return events, nil
}
