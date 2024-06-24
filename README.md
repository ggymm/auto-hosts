# auto-hosts

### 流程

1. 根据域名，自动解析出对应的 IP 地址
2. 对解析出的 IP 地址使用 PING 命令检查延时
3. 选择延时最低的 IP 地址，最终生成 hosts 文件


### 引用库

[dns](https://github.com/miekg/dns)
[ping](https://github.com/prometheus-community/pro-bing)
[govcl](https://github.com/ying32/govcl)


### 软件截图