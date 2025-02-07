# split_log.sh 日志分割脚本

# check_restart.sh 检测服务是否挂掉，挂掉就重启

# check_restart_abs.sh 三秒跑一次check_restart.sh

# crontab_check_restart.sh 1分钟跑一次检测check_restart_abs.sh split_log.sh是否挂掉，挂掉就重启

# crontab配置，自动启动

*/1 * * * * bash /home/go/src/go_project_template/scripts/go_project_template_crontab_check_restart.sh