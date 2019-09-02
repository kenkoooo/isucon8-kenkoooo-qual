# MariaDB in Docker

```bash
# MariaDB を起動
docker run --name mariadb -e MYSQL_ROOT_PASSWORD=isucon --cpuset-cpus=0 -d -p 43306:3306 mariadb
./db/init-user.sh

# アプリを起動
cd webapp/go/
make deps
./run_local.sh

# ベンチマーク
./bench.sh
```

# やったこと
- スコア 40707
- JOIN に使われている reservations の各フィールドにインデックスを張った。
- 残席の集計用テーブルを作った。
- getEvents を潰した平たくした。

# 学びなど

- echo のエラーハンドリングは `e.HTTPErrorHandler = customHTTPErrorHandler` これだけでデバッグがクソほど楽になる
- Scan で nullable な値を採る時は `sql.Null...` が便利
- TRIGGER で同期した集計用テーブルを別に作ったら速くなった。
