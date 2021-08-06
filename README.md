# iRODS FUSE Lite Monitor
A monitoring service for iRODS FUSE Lite

## Build
Build an executable using `Makefile`. The executable will be created under `bin`.
```shell script
make build
```

## Run
Run `bin/irodsfs-monitor` to start the service. Additionally user can specify a service port number.

Available arguments are:

- `-p`: service port number
- `-f`: run the service in foreground



## APIs
Available REST/HTTP APIs are:

HTTP Method | API URL           | Description
------------|-------------------|-------------------------------------------
`GET`       | `/instances`      | list all iRODS FUSE Lite instances running
`GET`       | `/instances/<id>` | get an iRODS FUSE Lite instance
`POST`      | `/instances`      | report a new iRODS FUSE Lite instance
`GET`       | `/transfers`      | list all data transfers
`GET`       | `/transfers/<id>` | list all data transfers of an iRODS FUSE Lite instance
`POST`      | `/transfers`      | report a new data transfer performed by an iRODS FUSE Lite instance


