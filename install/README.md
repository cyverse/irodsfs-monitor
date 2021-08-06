# Setup irodsfs-monitor systemd service

Copy the irodsfs-monitor binary `bin/irodsfs-monitor` to `/usr/bin/`.

Copy the systemd service `irodsfs-monitor.service` to `/usr/lib/systemd/system/`.

Create a service user `irodsfsmonitor`.
```bash
sudo adduser -r -d /dev/null -s /sbin/nologin irodsfsmonitor
```

Copy the irodsfs-monitor configuration `irodsfs-monitor.conf` to `/etc/irodsfs-monitor/`.
Be sure that this file must be only accessible by the `irodsfsmonitor` user.

Start the service.
```bash
sudo service irodsfs-monitor start
```