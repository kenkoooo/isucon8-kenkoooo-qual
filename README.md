# MariaDB in Docker

```bash
docker run --name mariadb -e MYSQL_ROOT_PASSWORD=isucon --cpuset-cpus=0 -d -p 43306:3306 mariadb
mysql -h 127.0.0.1 -uroot --port 43306 -pisucon
```
