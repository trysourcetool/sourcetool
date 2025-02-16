#!/bin/sh

set -e

login() {
  curl -v -c cookies.txt -X POST http://auth.local.trysourcetool.com:8080/api/v1/users/signin \
    -H "Content-Type: application/json" \
    -d '{"email":"john.doe@acme.com", "password":"password"}'
}

refresh_token() {
  curl -v -b cookies.txt -X POST http://acme.local.trysourcetool.com:8080/api/v1/users/refreshToken
}

users_me() {
  curl -v -b cookies.txt http://acme.local.trysourcetool.com:8080/api/v1/users/me
}

list_environments() {
  curl -v -b cookies.txt http://acme.local.trysourcetool.com:8080/api/v1/environments
}

create_api_key() {
  curl -v -b cookies.txt http://acme.local.trysourcetool.com:8080/api/v1/apiKeys \
    -H "Content-Type: application/json" \
    -d '{"environmentId":"0193dab0-2f93-7d13-b725-97afa7f028a3", "name": "API Key 1"}'
}

list_api_keys() {
  curl -v -b cookies.txt http://acme.local.trysourcetool.com:8080/api/v1/apiKeys
}

list_pages() {
  curl -v -b cookies.txt http://acme.local.trysourcetool.com:8080/api/v1/pages
}

case "$1" in
  "login") login ;;
  "refresh") refresh_token ;;
  "users_me") users_me ;;
  "list_environments") list_environments ;;
  "create_api_key") create_api_key ;;
  "list_api_keys") list_api_keys ;;
  "list_pages") list_pages ;;
esac
