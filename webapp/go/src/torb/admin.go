package main

import (
	"github.com/labstack/echo"
	"strconv"
)

func getSalesById(c echo.Context) error {
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return resError(c, "not_found", 404)
	}

	e, err := getEvent(eventID, -1)
	if err != nil {
		return err
	}

	query := `
	SELECT r.*, s.rank AS sheet_rank, s.num AS sheet_num, s.price AS sheet_price, e.price AS event_price 
	FROM reservations r 
	INNER JOIN sheets s ON s.id = r.sheet_id 
	INNER JOIN events e ON e.id = r.event_id 
	WHERE r.event_id = ? 
	ORDER BY reserved_at ASC
	`
	rows, err := db.Query(query, e.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var r Reservation
		var s Sheet
		if err := rows.Scan(&r.ID, &r.EventID, &r.SheetID, &r.UserID, &r.ReservedAt, &r.CanceledAt, &s.Rank, &s.Num, &s.Price, &e.Price); err != nil {
			return err
		}
		report := Report{
			ReservationID: r.ID,
			EventID:       e.ID,
			Rank:          s.Rank,
			Num:           s.Num,
			UserID:        r.UserID,
			SoldAt:        r.ReservedAt.Format("2006-01-02T15:04:05.000000Z"),
			Price:         e.Price + s.Price,
		}
		if r.CanceledAt != nil {
			report.CanceledAt = r.CanceledAt.Format("2006-01-02T15:04:05.000000Z")
		}
		reports = append(reports, report)
	}
	return renderReportCSV(c, reports)
}

func getSales(c echo.Context) error {
	query := `
	SELECT 
	r.*, s.rank AS sheet_rank, s.num AS sheet_num, s.price AS sheet_price, e.id AS event_id, e.price AS event_price 
	FROM reservations r 
	INNER JOIN sheets s ON s.id = r.sheet_id 
	INNER JOIN events e ON e.id = r.event_id 
	ORDER BY reserved_at ASC 
	`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var r Reservation
		var s Sheet
		var e Event
		if err := rows.Scan(&r.ID, &r.EventID, &r.SheetID, &r.UserID, &r.ReservedAt, &r.CanceledAt, &s.Rank, &s.Num, &s.Price, &e.ID, &e.Price); err != nil {
			return err
		}
		report := Report{
			ReservationID: r.ID,
			EventID:       e.ID,
			Rank:          s.Rank,
			Num:           s.Num,
			UserID:        r.UserID,
			SoldAt:        r.ReservedAt.Format("2006-01-02T15:04:05.000000Z"),
			SoldAtUnix:    r.ReservedAt.Unix(),
			Price:         e.Price + s.Price,
		}
		if r.CanceledAt != nil {
			report.CanceledAt = r.CanceledAt.Format("2006-01-02T15:04:05.000000Z")
		}
		reports = append(reports, report)
	}
	return renderReportCSV(c, reports)
}

func adminTop(c echo.Context) error {
	var events []*Event
	administrator := c.Get("administrator")
	if administrator != nil {
		var err error
		if events, err = getEvents(true); err != nil {
			return err
		}
	}
	return c.Render(200, "admin.tmpl", echo.Map{
		"events":        events,
		"administrator": administrator,
		"origin":        c.Scheme() + "://" + c.Request().Host,
	})
}
