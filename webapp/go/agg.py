import json
import re

with open("./echo.log") as f:
    lines = f.readlines()

results = [json.loads(s) for s in lines]
d = {}
for row in results:
    latency = row["latency"]
    path = row["uri"]
    if re.match(r"/api/events/\d+/actions/reserve", path):
        path = "/api/events/:id/actions/reserve"
    elif re.match(r"/admin/api/reports/events/\d+/sales", path):
        path = "/admin/api/reports/events/:id/sales"
    elif re.match(r"/api/users/\d+", path):
        path = "/api/users/:id"
    elif re.match(r"/api/events/\d+", path):
        path = "/api/events/:id"
    elif re.match(r"/admin/api/events/\d+", path):
        path = "/admin/api/events/:id"

    cur = d.get(path, 0)
    d[path] = cur + latency

result = [(t, path) for path, t in d.items()]
result = sorted(result, reverse=True)
for t, path in result:
    time_ms = t // 1000 // 1000
    print("{}\t{}".format(str(time_ms).rjust(10), path))
