#!/bin/sh

LOG_ROOT=/home/go/src/go_project_template/target/runtime_log
TARGET_ROOT=$LOG_ROOT

SERVER_ALL=(
  logic
)

while :; do
  #这里必须一分钟执行一次，否则天数据里面的内容会重复
  sleep 60

  YEAR=$(date +%Y)
  MONTH=$(date +%m)
  DAY=$(date +%d)

  #存放文件的路径
  LOG_PATH_ALL=(
    $LOG_ROOT/logic
    $LOG_ROOT/logic/$YEAR/$MONTH
  )

  #创建目录
  #
  for log_path in ${LOG_PATH_ALL[@]}; do
    if [ ! -d $log_path ]; then
      mkdir -p $log_path
    fi
  done

  for server in ${SERVER_ALL[@]}; do
    #这里用追加的形式
    cat $LOG_ROOT/${server}.log >>$LOG_ROOT/${server}/${server}_$(date +%Y%m%d%H%M).log
    #  #清空旧日志，这里不能用删除，删除或者mv后，不会产生新文件
    #  #注意这种方式清空，脚本日志重定向的时候要>> 不要> ，>不会重置文件大小 例如nohup ./a.sh > b & 改成 nohup ./a.sh >> b &
    cat /dev/null >$LOG_ROOT/${server}.log
    #  #将每分钟的日志汇入每日日志
    cat $LOG_ROOT/${server}/${server}_$(date -d "-1 minute" +"%Y%m%d%H%M").log >>$LOG_ROOT/${server}/${YEAR}/${MONTH}/${server}_$(date +%Y%m%d).log
    #    echo "\r" >>$LOG_ROOT/${server}/${YEAR}/${MONTH}/${server}_$(date +%Y%m%d).log
    #  #删除10分钟以前的分钟日志
    rm -f $LOG_ROOT/${server}/${server}_$(date -d "-3 minute" +"%Y%m%d%H%M").log
  done
done
