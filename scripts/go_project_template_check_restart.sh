#!/bin/sh

ExNumber=$(ps -ef | grep -w target/go_project_template_logic | grep -v grep | wc -l)
#SplitLogNumber=$(ps -ef | grep -w scripts/split_log.sh | grep -v grep | wc -l)
if [ $ExNumber -le 0 ]; then
  cd /home/go/src/go_project_template/
  make stop
  #  sleep 1s
  pkill -f target/go_project_template_logic
  #  pkill -f scripts/split_log.sh
  #  pkill -f scripts/check_restart_abs.sh
  make run
  echo "target restart"
fi

SplitLogNumber=$(ps -ef | grep -w go_project_template_split_log.sh | grep -v grep | wc -l)
if [ $SplitLogNumber -le 0 ]; then
  cd /home/go/src/go_project_template/
  nohup scripts/go_project_template_split_log.sh >>scripts/log/split_log.log 2>&1 &
  echo "go_project_template_split_log.sh restart"
fi
