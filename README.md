# Go Expense Tracker

[![Go Report Card](https://goreportcard.com/badge/github.com/mhkarimi1383/goExpenseTracker)](https://goreportcard.com/report/github.com/mhkarimi1383/goExpenseTracker)

<!-- 
## Demo -->
[![Demo](./demo.gif)](./demo.gif)

> Already deployed at [tracker.karimi.dev](https://tracker.karimi.dev) just contact me to get access

## Powered by

* GoLang
* GoTemplate
* OpenIDConnect (like Keycloak or Google)
* MongoDB
* Bootstrap
* Flat Remix CSS
* Based on [goAPIBaseProject](https://github.com/mhkarimi1383/goAPIBaseProject)

## How to bring it up?

since I'm using onedev in my lab you can see that [.onedev-buildspec.yml](.onedev-buildspec.yml) contains configuration needed to build and deploy the project

or

set required variables and run golang project directly

as you can see there are some configurations you need to set in order to build and deploy the project

you can use config file called `config.yml` or set Environment variables

set environment variables like below and bring up the project

Environment Variables | Config File Key | Description | Default Value
---------------------|----------------|-------------|-------------
METRIC_ADDRESS | metric_address | Bind address for prometheus exporter | :9090
API_ADDRESS | api_address | Bind address for main server | :8080
SENTRY_DSN | sentry_dsn | DSN address for sending logs and errors to sentry | -
LOG_FORMAT | log_format | Format for logging can be either `text` or `json` | text
APPLICATION_TITLE | application_title | Title of application used for Database, etc. | goExpenseTracker
APPLICATION_DESCRIPTION | application_description | Description of application (not used in this project) | -
MONGODB_CONNECTION_URI | mongodb_connection_uri | URI for MongoDB connection | -
OPENID_ADDRESS | openid_address | Endpoint of OpenID Provider | -
OPENID_CLIENT_SECRET | openid_client_secret | Secret from OpenID Connect Provider | -
OPENID_CLIENT_ID | openid_client_id | ID for OpenID Connect Provider Client | -
OPENID_USERNAME_KEY | openid_username_key | Unique Key from OpenID Connect Provider to use as username (refer to provider documentation about jwt Token Spec) | preferred_username
BASE_URL | base_url | Base URL for application with protocol (`http` or `https` and no trailing slash [`/`]) | `http://127.0.0.1:8080`
