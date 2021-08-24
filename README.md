# servicenow-reporter-middleware

Simple tool, which sends requests to the REST API of Service-Now by url, mentioned in config.json, compares it with prev data from Redis and if there is a difference - sends a notification to the Telegram chat.

config.json should contain 4 fields:

url: Url of Service-Now API, such as https://instance.service-now.com/api/now/table/incident?active=true
login: Service-Now internal Login
password: Service-Now internal Password
period: period in seconds, how often checks will be made
