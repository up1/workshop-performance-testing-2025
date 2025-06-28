# Workshop :: Performance Testing


## 1. Start Database with MySQL
* Insert data for test (10,000 records)
```
$docker compose up -d db
$docker compose ps
$docker compose logs --follow
```

## 1. Bad API
* Golang and Gin
* MySQL driver

```
$docker compose build api-bad
$docker compose up -d api-bad
$docker compose ps
$docker compose logs --follow
```

## 2. Performance testing with K6
```
$cd k6
$K6_WEB_DASHBOARD=true k6 run loadtest.js
```

Open K6 Dashboard
* http://localhost:5665/

## 3. Visualize K6 results
* Prometheus
* Grafana

### Start Prometheus
```
$docker compose up -d prometheus
$docker compose ps
$docker compose logs --follow
```

### Start Grafana
```
$docker compose up -d grafana
$docker compose ps
$docker compose logs --follow
```

Access to Grafana with url = http://localhost:3000
* user=admin
* password=admin

## 4. Performance testing with K6 and Docker
```
$docker compose build k6
$docker compose up k6
$docker compose ps
```

Open K6 Dashboard
* http://localhost:5665/

## 5. Monitoring MySQL Database
* [MySQL Exporter](https://github.com/prometheus/mysqld_exporter)
* [Grafana dashboard](https://grafana.com/grafana/dashboards/14057-mysql/)

```
$docker compose up -d prom_mysql_exporter
$docker compose ps
```

Access to metric of mysql
* http://localhost:9104