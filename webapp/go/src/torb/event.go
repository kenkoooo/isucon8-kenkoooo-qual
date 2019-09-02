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

	rows, err := db.Query("SELECT * FROM sheets ORDER BY `rank`, num")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sheet Sheet
		if err := rows.Scan(&sheet.ID, &sheet.Rank, &sheet.Num, &sheet.Price); err != nil {
			return nil, err
		}
		event.Sheets[sheet.Rank].Price = event.Price + sheet.Price
		event.Total++
		event.Sheets[sheet.Rank].Total++

		var reservation Reservation
		err := db.QueryRow("SELECT * FROM reservations WHERE event_id = ? AND sheet_id = ? AND canceled_at IS NULL GROUP BY event_id, sheet_id HAVING reserved_at = MIN(reserved_at)", event.ID, sheet.ID).Scan(&reservation.ID, &reservation.EventID, &reservation.SheetID, &reservation.UserID, &reservation.ReservedAt, &reservation.CanceledAt)
		if err == nil {
			sheet.Mine = reservation.UserID == loginUserID
			sheet.Reserved = true
			sheet.ReservedAtUnix = reservation.ReservedAt.Unix()
		} else if err == sql.ErrNoRows {
			event.Remains++
			event.Sheets[sheet.Rank].Remains++
		} else {
			return nil, err
		}

		event.Sheets[sheet.Rank].Detail = append(event.Sheets[sheet.Rank].Detail, &sheet)
	}

	return &event, nil
}
