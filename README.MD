[![Build Status](https://drone.pi.scjtqs.com:8443/api/badges/scjtqs/jd_cookie/status.svg)](https://drone.pi.scjtqs.com:8443/scjtqs/jd_cookie)
# 说明

## v1.x 版本(tag) 本地运行，本地提取
京东的cookie 本地提取工具(小白专用) 

双击exe文件，运行服务，后用本地的浏览器打开 http://127.0.0.1:29099

来辅助提取你的cookie

## v2.x 版本(tag) 服务器运行，客户端自行扫码提取,服务端自动更新cookie
> 已不受服务端和客户端在同ip的限制
> 
> 客户端通过浏览器打开， http://公网ip或者域名:29099/
> 
> 然后扫码即可提取。
>
> 
```bash
## linux/macos下
chmod +x jdcookie
## 加上-d 代表后台运行
./jdcookie -d
## windows 下
## 可以直接双击运行。也可以cmd中运行
## cmd 中运行 参数和linux一样
.\jdcookie.exe -d
```
v2.0.5
> 应网友 @king6207 的要求，前端页面，提取cookie成功后增加一个直接复制cookie的按钮。
> 


v2.0.4
> 调整更新的sql语法，将replace into替换为 INSERT .... ON DUPLICATE KEY UPDATE。修复更新cookie引发id变更的问题（replace更新是先删在insert)。因小伙伴要求id顺序不能变。 

v2.0.3
> 数据库做调整、主键变更为id字段，方便一些用户做排序。

v2.0.2
> 增加了cookie存入db的支持

v 2.0.1
> 环境变量 UPSAVE cookie的推送地址。
>
> 简单做了一个docker镜像 scjtqs/jd_cookie:latest 
> 
> 支持  arm/arm64/i386/amd64

## 编译方法 (For v2.x tag)
```bash
#安装编译工具，用于打包html到二进制
go get -u github.com/gobuffalo/packr/v2/packr2
#build和go build一样，交叉编译的方式也是一样的
CGO_ENABLED=0 GOOS=linux GOARCH=arm packr2 build  -o jdcookie
#简单点的
packr2 build  -o jdcookie
# 编译windows x86版本
CGO_ENABLED=0 GOOS=windows GOARCH=386 packr2 build  -o jdcookie_windows_x86.exe
# 编译windows x64版本
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 packr2 build  -o jdcookie_windows_x64.exe
# 编译mac arm64 (m1)版本
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 packr2 build  -o jdcookie_darwin_arm64
# 编译mac amd64 x64版本
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 packr2 build  -o jdcookie_darwin_x64
```

## v3.x版本(tag) 服务端运行，db多类型支持。同样支持直接下载release中的exe文件本地运行。
> DB_TYPE 可以选择 mysql、postgres、sqlite3、mssql
> 
> 因为增加了sqlite3的支持，引入了cgo。不方便直接交叉编译了。于是 v3.x版本开始，就主要docker版了。docker支持 linux/amd64、linux/arm/v7、linux/arm64、linux/386。这4种架构平台。
> 
> release中的[下载](https://github.com/scjtqs/jd_cookie/releases)二进制包，无法使用sqlite3数据库(为了交叉编译，关闭了cgo）(仅linux_amd64的预编译包可以使用sqlite3)，要使用sqlite3数据库的，请自行在需要运行的机器上编译。或者使用docker版本。

v3.0.7
> 修复ua
> 
> 感谢 @EthanWDC

v3.0.6
> 提供京东app方式登录后，页面上提取cookie。方便仅提取cookie的玩家。
> 
> 会话采用session的方式，避开了cookie的安全隐患。
> 
> 先唤起京东app，点击登录。。根据服务器网络情况，等个1-2秒，回到之前浏览器，点击 "获取通过京东app跳转登录的cookie" 即可拿到对应的cookie（仅能提取一次，提取成功后会立马销毁，数据缓存在内存中，不落地到服务器文件）。

v3.0.5
> 添加手机页面唤起京东app直接获取cookie。（因小伙伴想要玩这功能
> 
> 跳转jdapp方式登录，无法在web上直接拿到cookie，只能在后台控制台打印看到（跳转到了jdapp。无法跳回去），或者推送到指定url、或者写入db中。

v3.0.4
> 修复ua 远程登录问题

v3.0.3
> 使用golang 1.16新特性，静态文件打包采用 embed包处理，直接go build就可以了。不再依赖packr2。直接go build即可。不过需要安装 go 1.16+版本。

v3.0.2
> 修复ua失效问题

v3.0.1
> 除了mysql外，应hcl要求，增加了其他db类型支持: postgres、sqlite3、mssql
> 

## 编译方法 （For v3.x tag)
```bash
#安装golang环境
#安装gcc等编译环境
## for alpine
apk add --update gcc musl-dev sqlite-static sqlite-dev
## for Fedora
sudo yum groupinstall "Development Tools" "Development Libraries"
## for Ubuntu
sudo apt-get install build-essential
## for Mac OSX
brew install sqlite3
brew upgrade icu4c
go install github.com/mattn/go-sqlite3
#linux/mac/UNIX下
CGO_ENABLED=1 go build  -o jdcookie
#查看帮助
./jdcookie -h 
#windows下
CGO_ENABLED=1 go build  -o jdcookie.exe
#查看帮助
.\jdcookie.exe -h
```


## docker 使用示例

#### quick start
```bash
docker run -d -p 29099:29099  scjtqs/jd_cookie:latest
```

### with environment
```bash
## 更新镜像
docker pull jd_cookie scjtqs/jd_cookie:latest
## UPSAVE 默认是空，不会推送到其他地址；DB_ENABLE默认为false，不会记录cookie到db。
docker run -d \
-p 29099:29099 \
-e UPSAVE="http://192.168.0.12:8080/savecookie" \
-e DB_ENABLE="true" \
-e DB_HOST="192.168.0.13" \
-e DB_PORT="3306" \
-e DB_USER="root" \
-e DB_PASS="abcde" \
-e DB_DATABASE="cookies" \
-e DB_TYPE="mysql" \
--name jd_cookie scjtqs/jd_cookie:latest
```
### for sqlite3 usage
```bash
docker run -d \
-p 29099:29099 \
-v `pwd`/db/:/opt/app/db/ \
-e UPSAVE="http://192.168.0.12:8080/savecookie" \
-e DB_ENABLE="true" \
-e DB_HOST="./db/cookies.db" \
-e DB_TYPE="sqlite3" \
--name jd_cookie scjtqs/jd_cookie:latest
```

### 佛系更新，够用就行。欢迎pr和建议。欢迎start
