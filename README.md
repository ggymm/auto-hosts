# auto-hosts

### 流程

1. 根据域名，到 DNS 服务器解析出 A 类 IP 地址
2. 对解析出的 IP 地址使用 PING 命令检查延时
3. 选择延时最低的 IP 地址，最终生成 hosts 文件


### hosts

```text
140.82.113.26 alive.github.com
20.27.177.116 api.github.com
185.199.108.153 assets-cdn.github.com
185.199.110.133 avatars.githubusercontent.com
185.199.109.133 avatars0.githubusercontent.com
185.199.109.133 avatars1.githubusercontent.com
185.199.109.133 avatars2.githubusercontent.com
185.199.109.133 avatars3.githubusercontent.com
185.199.109.133 avatars4.githubusercontent.com
185.199.110.133 avatars5.githubusercontent.com
185.199.108.133 camo.githubusercontent.com
140.82.112.21 central.github.com
185.199.109.133 cloud.githubusercontent.com
20.27.177.114 codeload.github.com
140.82.112.21 collector.github.com
185.199.110.133 desktop.githubusercontent.com
185.199.110.133 favicons.githubusercontent.com
243.185.187.39 gist.github.com
52.217.235.225 github-cloud.s3.amazonaws.com
54.231.136.169 github-com.s3.amazonaws.com
16.182.72.57 github-production-release-asset-2e65be.s3.amazonaws.com
52.217.113.65 github-production-repository-file-5c1aeb.s3.amazonaws.com
54.231.137.161 github-production-user-asset-6210df.s3.amazonaws.com
192.0.66.2 github.blog
20.27.177.113 github.com
140.82.113.17 github.community
185.199.110.154 github.githubassets.com
151.101.37.194 github.global.ssl.fastly.net
185.199.110.153 github.io
185.199.108.133 github.map.fastly.net
185.199.110.153 githubstatus.com
140.82.113.26 live.github.com
185.199.109.133 media.githubusercontent.com
185.199.108.133 objects.githubusercontent.com
13.107.42.16 pipelines.actions.githubusercontent.com
185.199.109.133 raw.githubusercontent.com
185.199.110.133 user-images.githubusercontent.com
185.199.110.133 private-user-images.githubusercontent.com
13.107.246.64 vscode.dev
140.82.112.21 education.github.com
```


### 依赖库

[dns](https://github.com/miekg/dns)

[ping](https://github.com/prometheus-community/pro-bing)

[govcl](https://github.com/ying32/govcl)


### DNS服务器

```text
1.1.1.1
1.2.4.8
4.2.2.1
8.8.8.8
8.20.247.20
8.26.56.26
9.9.9.9
45.11.45.11
64.6.64.6
74.82.42.42
77.88.8.8
80.80.80.80
84.200.69.80
94.140.14.14
101.101.101.101
101.226.4.6
114.114.114.114
119.29.29.29
156.154.70.1
168.126.63.1
180.76.76.76
180.184.1.1
182.254.118.118
185.222.222.222
195.46.39.39
199.85.126.10
202.120.2.100
208.67.222.222
210.2.4.8
223.5.5.5
```


### 软件截图

![image](https://github.com/ggymm/auto-hosts/assets/20694755/7b6f0053-7c2e-496e-b0e9-297fd92d95f2)
