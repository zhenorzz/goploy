# Server

## Configure

Add target server in goploy and get the server id.

Run [goploy-agent](https://github.com/zhenorzz/goploy-agent) in target server.


## Target illustrate

- CPU usage rate: /proc/stat add all columns and minis the 4th column
- RAM usage rate: /proc/meminfo Memfree/MemTotal
- Loadavg:  /proc/loadavg 
- TCP: /proc/net/tcp established、total
- Public network bandwidth: /proc/net/dev collect the 1st(in) and 9th(out) column the eth column
- Intranet bandwidth: /proc/net/dev collect the 1st(in) and 9th(out) column the lo column
- Disk usage rate: df --output=pcent,ipcent,source
- Disk IO: /sys/block/<dev>/stat collect the 0th(read iops) and 4th(write iops) column