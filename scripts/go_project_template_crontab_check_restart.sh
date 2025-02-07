#!/bin/sh

source /etc/profile
source ~/.bash_profile

AbsNumber=$(ps -ef | grep -w go_project_template_check_restart_abs.sh | grep -v grep | wc -l)
if [ $AbsNumber -le 0 ]; then
  cd /home/go/src/go_project_template/
  nohup scripts/go_project_template_check_restart_abs.sh >>scripts/log/check_restart_abs.log 2>&1 &
  echo "go_project_template_check_restart_abs.sh restart"
fi
