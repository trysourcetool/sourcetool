# Sourcetool API

## Setup hosts
```
127.0.0.1 auth.local.trysourcetool.com
127.0.0.1 acme.local.trysourcetool.com
```

## Setup server
```
$ docker network create sourcetool_server_default
$ make dc-build
$ make dc-up
```

## Reset server
```
$ make remove-postgres-data
$ make dc-build
$ make dc-up
```

## Cleanup Docker
```
$ make remove-docker-images
$ make remove-docker-builder
```


## Swagger
```
$ make open-swagger
```