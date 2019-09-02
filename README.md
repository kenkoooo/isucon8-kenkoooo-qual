# MariaDB in Docker

```bash
docker run --name mariadb -e MYSQL_ROOT_PASSWORD=isucon --cpuset-cpus=0 -d -p 43306:3306 mariadb
mysql -h 127.0.0.1 -uroot --port 43306 -pisucon
```

# 結果
- スコア 40707

# 学びなど

- echo のエラーハンドリングは `e.HTTPErrorHandler = customHTTPErrorHandler` これだけでデバッグがクソほど楽になる
- Scan で nullable な値を採る時は `sql.Null...` が便利
- TRIGGER で同期した集計用テーブルを別に作ったら速くなった。
