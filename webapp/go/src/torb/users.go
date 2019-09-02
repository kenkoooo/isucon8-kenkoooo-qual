package main

import (
	"strconv"

	"github.com/labstack/echo"
)

func getUser(c echo.Context) error {
	contextUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	user, err := getLoginUser(c)
	if err != nil {
		return err
	}
	if user.ID != contextUserID {
		return resError(c, "forbidden", 403)
	}

	query := `
	SELECT r.*,
	s.rank AS sheet_rank, s.num AS sheet_num, s.price AS s_price,
	e.id AS e_id, e.title AS e_title, e.public_fg, e.closed_fg, e.price AS e_price
	FROM reservations r 
	INNER JOIN sheets s ON s.id = r.sheet_id
	INNER JOIN events e ON e.id = r.event_id
	WHERE r.user_id = ? 
	ORDER BY IFNULL(r.canceled_at, r.reserved_at) DESC 
	LIMIT 5
	`

	rows, err := db.Query(query, user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var recentReservations []Reservation
	for rows.Next() {
		var r Reservation
		var s Sheet
		var e Event
		if err := rows.Scan(
			&r.ID, &r.EventID, &r.SheetID, &r.UserID, &r.ReservedAt, &r.CanceledAt,
			&s.Rank, &s.Num, &s.Price,
			&e.ID, &e.Title, &e.PublicFg, &e.ClosedFg, &e.Price); err != nil {
			return err
		}

		price := e.Price + s.Price
		e.Sheets = nil
		e.Total = 0
		e.Remains = 0

		r.Event = &e
		r.SheetRank = s.Rank
		r.SheetNum = s.Num
		r.Price = price
		r.ReservedAtUnix = r.ReservedAt.Unix()
		if r.CanceledAt != nil {
			r.CanceledAtUnix = r.CanceledAt.Unix()
		}
		recentReservations = append(recentReservations, r)
	}
	if recentReservations == nil {
		recentReservations = make([]Reservation, 0)
	}

	var totalPrice int
	if err := db.QueryRow("SELECT IFNULL(SUM(e.price + s.price), 0) FROM reservations r INNER JOIN sheets s ON s.id = r.sheet_id INNER JOIN events e ON e.id = r.event_id WHERE r.user_id = ? AND r.canceled_at IS NULL", user.ID).Scan(&totalPrice); err != nil {
		return err
	}

	rows, err = db.Query("SELECT event_id FROM reservations WHERE user_id = ? GROUP BY event_id ORDER BY MAX(IFNULL(canceled_at, reserved_at)) DESC LIMIT 5", user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var recentEvents []*Event
	for rows.Next() {
		var eventID int64
		if err := rows.Scan(&eventID); err != nil {
			return err
		}
		event, err := getEvent(eventID, -1)
		if err != nil {
			return err
		}
		for k := range event.Sheets {
			event.Sheets[k].Detail = nil
		}
		recentEvents = append(recentEvents, event)
	}
	if recentEvents == nil {
		recentEvents = make([]*Event, 0)
	}

	return c.JSON(200, echo.Map{
		"id":                  user.ID,
		"nickname":            user.Nickname,
		"recent_reservations": recentReservations,
		"total_price":         totalPrice,
		"recent_events":       recentEvents,
	})
}
