package main

import (
	"database/sql"
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

	totalSheetCount, err := GetTotalSheetByRank()
	if err != nil {
		return nil, err
	}

	event.Total = 0
	event.Sheets = map[string]*Sheets{
		"S": &Sheets{},
		"A": &Sheets{},
		"B": &Sheets{},
		"C": &Sheets{},
	}
	for _, partialSheet := range totalSheetCount {
		event.Total += partialSheet.Total
		event.Remains += partialSheet.Total
		event.Sheets[partialSheet.Rank].Total = partialSheet.Total
		event.Sheets[partialSheet.Rank].Remains = partialSheet.Total
		event.Sheets[partialSheet.Rank].Price = event.Price + partialSheet.Price
	}

	RefineEvent(reservations, sheets, userIDs, &event, loginUserID)

	return &event, nil
}

func RefineEvent(reservations []*Reservation, sheets []*Sheet, userIDs []*sql.NullInt64, event *Event, loginUserID int64) {
	for i, r := range reservations {
		s := sheets[i]
		userID := userIDs[i]

		if r.ReservedAt != nil {
			s.Mine = userID.Valid && userID.Int64 == loginUserID
			s.Reserved = true
			s.ReservedAtUnix = r.ReservedAt.Unix()
			event.Remains--
			event.Sheets[s.Rank].Remains--
		}

		event.Sheets[s.Rank].Detail = append(event.Sheets[s.Rank].Detail, s)
	}
}

type PartialSheet struct {
	Rank  string
	Total int
	Price int64
}

func GetTotalSheetByRank() (map[string]PartialSheet, error) {
	rows, err := db.Query("SELECT rank, COUNT(*), MAX(price) FROM sheets GROUP BY rank")
	if err != nil {
		return nil, err
	}
	s := make(map[string]PartialSheet)
	for rows.Next() {
		var p PartialSheet
		if err := rows.Scan(&p.Rank, &p.Total, &p.Price); err != nil {
			return nil, err
		}
		s[p.Rank] = p
	}
	return s, nil
}

func getEvents(all bool) ([]*Event, error) {
	totalSheetCount, err := GetTotalSheetByRank()
	if err != nil {
		return nil, err
	}

	eventMap := make(map[int64]*Event)
	{
		rows, err := db.Query("SELECT * FROM events")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var e Event
			if err := rows.Scan(&e.ID, &e.Title, &e.PublicFg, &e.ClosedFg, &e.Price); err != nil {
				return nil, err
			}

			e.Total = 0
			e.Sheets = map[string]*Sheets{
				"S": &Sheets{},
				"A": &Sheets{},
				"B": &Sheets{},
				"C": &Sheets{},
			}
			for _, partialSheet := range totalSheetCount {
				e.Total += partialSheet.Total
				e.Remains += partialSheet.Total
				e.Sheets[partialSheet.Rank].Total = partialSheet.Total
				e.Sheets[partialSheet.Rank].Remains = partialSheet.Total
				e.Sheets[partialSheet.Rank].Price = e.Price + partialSheet.Price
			}
			eventMap[e.ID] = &e
		}
	}
	query := "SELECT sheet_rank, sold_count, event_id FROM sold"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rank string
		var sold int
		var eventId int64
		if err := rows.Scan(&rank, &sold, &eventId); err != nil {
			return nil, err
		}
		eventMap[eventId].Sheets[rank].Remains -= sold
		eventMap[eventId].Remains -= sold
	}

	var events Events
	for _, e := range eventMap {
		if !all && !e.PublicFg {
			continue
		}
		events = append(events, e)
	}
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
