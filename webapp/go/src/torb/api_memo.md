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