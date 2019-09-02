# GET /
- getEvents
```
{
  "events": [],
  "user": ...,
  "origin": ...
}
```


# POST /api/users
ユーザー作成?
```
{
    "nickname": string,
    "login_name": string,
    "password": string
}
```
- `SELECT * FROM users WHERE login_name = ?`
- `INSERT INTO users (login_name, pass_hash, nickname) VALUES (?, SHA2(?, 256), ?)`
```
{
  "id": int64,
  "nickname": string
}
```

# GET /api/user/:id
- `:id` = user id
- `SELECT id, nickname FROM users WHERE id = ?`
- getLoginUser
- `SELECT r.*, s.rank AS sheet_rank, s.num AS sheet_num FROM reservations r INNER JOIN sheets s ON s.id = r.sheet_id WHERE r.user_id = ? ORDER BY IFNULL(r.canceled_at, r.reserved_at) DESC LIMIT 5`
- for getEvent
- `SELECT IFNULL(SUM(e.price + s.price), 0) FROM reservations r INNER JOIN sheets s ON s.id = r.sheet_id INNER JOIN events e ON e.id = r.event_id WHERE r.user_id = ? AND r.canceled_at IS NULL`
- for getEvent
```
{
    "id": user id,
    "nickname": nickname,
    "recent_reservations": [Resevation],
    "total_price": ...,
    "recent_events": [Event]
}
```

# POST /api/events/:id/actions/reserve
- `:id` event id

```
{
    "sheet_rank": ""
}
```

- getLoginUser
- getEvent
- `SELECT * FROM sheets WHERE id NOT IN (SELECT sheet_id FROM reservations WHERE event_id = ? AND canceled_at IS NULL FOR UPDATE) AND `rank` = ? ORDER BY RAND() LIMIT 1`
- `INSERT INTO reservations (event_id, sheet_id, user_id, reserved_at) VALUES (?, ?, ?, ?)`

```
{
    "id": resevation id,
    "sheet_rank": string,
    "sheet_num": int64
}
```

# POST /admin/api/events
...
- `INSERT INTO events (title, public_fg, closed_fg, price) VALUES (?, ?, 0, ?)`
...

# DELETE /api/events/:id/sheets/:rank/:num/reservation
...
- `UPDATE reservations SET canceled_at = ? WHERE id = ?`
...

# POST /admin/api/events/:id/actions/edit
...
- `UPDATE events SET public_fg = ?, closed_fg = ? WHERE id = ?`
...
