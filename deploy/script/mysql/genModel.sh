#!/usr/bin/env bash

# 使用方法:                             # 注释：脚本使用说明
# ./genModel.sh usercenter user         # 示例：生成单张表 user
# ./genModel.sh usercenter user_auth    # 示例：生成单张表 user_auth
# 再将./genModel下的文件剪切到对应服务的model目录里面，记得改package   # 操作提示


# 生成的表名                             # 说明：第二个参数为要生成的表名
tables=$2                               # 将脚本的第二个参数赋值给 tables 变量


# 表生成的 genmodel 目录                                 # 说明：模型文件输出目录
modeldir=../../../app/payment/model             # 定义输出目录为当前目录下的 genModel

# 数据库配置                             # 说明：以下为数据库连接信息
host=127.0.0.1                          # MySQL 主机地址
port=3306                               # MySQL 端口号
dbname=$1                               # 将脚本的第一个参数赋给 dbname（数据库名）
username=root                           # 数据库用户名
passwd=123456                           # 数据库密码


# 输出提示信息
echo "开始创建库: $dbname 的表: $tables"   # 打印将要操作的数据库和表名



# 调用 goctl 工具生成 model 代码   1.7.3: E:\CodeIDE\goctl\1.7.3
# 可以使用-home指定模版路径，默认是~/.goctl
goctl model mysql datasource --url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" --table="${tables}" --dir="${modeldir}" --home="E:\CodeIDE\goctl\1.7.3" -c




