PKG=github.com/cyverse/irodsfs-monitor
GO111MODULE=on
GOPROXY=direct
GOPATH=$(shell go env GOPATH)

.EXPORT_ALL_VARIABLES:

.PHONY: build
build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -o bin/irodsfs-monitor ./cmd/

.PHONY: install_centos
install_centos:
	cp bin/irodsfs-monitor /usr/bin
	cp install/irodsfs-monitor.service /usr/lib/systemd/system/
	adduser -r -d /dev/null -s /sbin/nologin irodsfsmonitor
	mkdir -p /etc/irodsfs-monitor
	cp install/irodsfs-monitor.conf /etc/irodsfs-monitor
	chown irodsfsmonitor /etc/irodsfs-monitor/irodsfs-monitor.conf
	chmod 660 /etc/irodsfs-monitor/irodsfs-monitor.conf

.PHONY: install_ubuntu
install_ubuntu:
	cp bin/irodsfs-monitor /usr/bin
	cp install/irodsfs-monitor.service /etc/systemd/system/
	adduser --system --home /dev/null --shell /sbin/nologin irodsfsmonitor
	mkdir -p /etc/irodsfs-monitor
	cp install/irodsfs-monitor.conf /etc/irodsfs-monitor
	chown irodsfsmonitor /etc/irodsfs-monitor/irodsfs-monitor.conf
	chmod 660 /etc/irodsfs-monitor/irodsfs-monitor.conf

.PHONY: uninstall
uninstall:
	rm -f /usr/bin/irodsfs-monitor
	rm -f /etc/systemd/system/irodsfs-monitor.service
	rm -f /usr/lib/systemd/system/irodsfs-monitor.service
	userdel irodsfsmonitor | true
	rm -rf /etc/irodsfs-monitor
