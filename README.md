# Go-BticinoThermostatInflux

Add your Bticino Thermostat measurations to a Influx DB!

---
## Requirements
- Golang compiler
- `github.com/influxdata/influxdb-client-go` Go library
    - Install via `go get github.com/influxdata/influxdb-client-go`
- Legrand API credentials
    - Create your API credentials at <https://portal.developer.legrand.com/>products. The default API (Starter Kit, which is free) has 500 API calls/day. It could take a few days for your API to get approved.
        - Just a suggestion: if you don't have a website to set for the redirect, you can use something like <https://httpbin.org/get>, so that when you do the first login it'll redirect you to a page with the code instead of copying it from the address bar.
---
## Configuration
- Add the required data in the config/config.go file
    - InfluxDB Section (**MINIMUM SERVER VERSION: 1.8**)
        - InfluxHost -> The address of the InfluxDB server (address:port)
        - InfluxToken -> If using a version between 1.8 (included) and 2 (excluded) use username:password (if anonymous login is allowed, you can leave it empty), otherwise fill with the Token.
        - InfluxOrg -> If using a version between 1.8 (included) and 2 (excluded) leave it empty, otherwise fill with the Organization.
        - InfluxBucket -> If using a version between 1.8 (included) and 2 (excluded) use the format database/retention-policy (if you want to use the default retention policy or you don't know what a retention policy is, fill with the database name), otherwise fill with the Bucket.
        - InfluxMeasurementName -> I don't know how else I should explain it, I usually use `environment`.
    - Thermostat Section
        - ClientID -> The Client ID of your API.
        - ClientSecret -> The Client Secret of your API.
        - Redirect -> The website you specified when you created the API credentials (**IT CANNOT BE DIFFERENT, OTHERWISE IT'LL GIVE YOU AN ERROR**).
        - PlantID -> The ID of the location where your Thermostat is.
        - ModuleID -> The ID of the Thermostat.
    - Various Section
        - CalculateAbsolutePath -> Whether the software should make the file names absolute (eg. from refresh.txt to /opt/Thermostat/refresh.txt). Leave it on true if you're gonna build it, set false if you want to use it via `go run TempGraph.go`.
    - Advanced Section (EDIT ONLY IF YOU KNOW WHAT YOU'RE DOING)
        - RefreshFileName -> The name of the file that'll contain the refresh token.
        - RequestDelay -> The time between the requests, the default value is 182 so you'll make about 500 requests per day (with a little margin, just to be sure)

---
## Installation
- Install the golang compiler
- Configure the software
- Build or run it (watch out, you **MUST** remember to edit var>CalculateAbsolutePath in the config file accordingly to what you choose!)
___

That's all Folks!
For help just [ask me on Telegram](https://t.me/Knocks)!

This Source Code Form is subject to the terms of the Apache-v2.0 License. If a copy of the Apache-V2.0 License was not distributed with this
file, You can obtain one at <https://www.apache.org/licenses/LICENSE-2.0>.
