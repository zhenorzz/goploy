# 服务器Agent

## 配置

goploy添加目标机器，获取服务器ID

目标机器运行[goploy-agent](https://github.com/zhenorzz/goploy-agent)


## 指标计算方式

- CPU使用率: /proc/stat 所有列加起来减去第四列（空闲时间）
- RAM使用率: /proc/meminfo Memfree/MemTotal
- Loadavg:  /proc/loadavg 
- TCP: /proc/net/tcp 统计established和total
- 外网带宽: /proc/net/dev 周期性采集第0列(eth)的第1列(in)和第9列(out)
- 内网带宽: /proc/net/dev 周期性采集第0列(lo)的第1列(in)和第9列(out)
- 硬盘使用率: df --output=pcent,ipcent,source
- 硬盘IO: /sys/block/<dev>/stat 周期性采集第0列(read iops)和第四列(write iops)