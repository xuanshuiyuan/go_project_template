# 车型参数更新/月初

0 0 1 * * curl http://127.0.0.1:15010/api/timer/car_model_extend_update

# 车辆异常告警/1分钟一次

*/1 * * * * curl http://127.0.0.1:15010/api/timer/vehicleAbnormalAlarm

# 车辆里程核验/每天23:59分

59 23 * * * curl http://127.0.0.1:15010/api/timer/vehicleMileageVerify

[comment]: <> (# 优惠券失效/每日凌晨)

[comment]: <> (0 0 * * * curl http://127.0.0.1:15010/api/timer/coupon_expiration)

[comment]: <> (建议用内网访问（提供ip）)

[comment]: <> (外网：dds-wz970a3a25620c941797-pub.mongodb.rds.aliyuncs.com:3717)

[comment]: <> (内网：dds-wz970a3a25620c941.mongodb.rds.aliyuncs.com:3717)

[comment]: <> (eobd库)

[comment]: <> (用户：eobd_ggc（只读）)

[comment]: <> (密码：eobd_ggc_20240509)

[comment]: <> (ip内: 172.29.223.60)

[comment]: <> (外：120.25.168.13)

[comment]: <> (port: 6379)

[comment]: <> (database: 10)

[comment]: <> (password: dxj@20200331)


#linux/centos安装最新版本chrome
wget http://dist.control.lth.se/public/CentOS-7/x86_64/google.x86_64/google-chrome-stable-124.0.6367.118-1.x86_64.rpm
yum install google-chrome-stable-124.0.6367.118-1.x86_64.rpm
#安装中文相关的字体
yum -y groupinstall Fonts
