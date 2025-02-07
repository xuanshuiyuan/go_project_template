#!/bin/sh

#项目名称
PROJECT_NAME=$1
#时间
TIME=$(date +%Y%m%d-%H%M%S)
BAKTIME=$(date +%Y%m%d-%H%M%S)
#版本
VERSION=${TIME}
#项目代码基础目录
BASE_PATH=/data/docker/code/go
#项目代码目录
PROJECT_PATH=${BASE_PATH}/${PROJECT_NAME}
#项目代码根目录
VERSION_PATH=${PROJECT_PATH}/${VERSION}

#go代码目录
GO_PATH=/home/go/src

#系统项目基础目录
PATH_WWW=/home/www
#项目代码下载目录
FILE_PATH=/data/file/${PROJECT_NAME}

###解压
cd ${FILE_PATH}
mv ${PROJECT_NAME}.tar ${TIME}.tar
mkdir -p ${VERSION_PATH}
tar -xvf ${TIME}.tar -C ${VERSION_PATH}
##

#将软连接指向新目录
rm -rf ${PATH_WWW}/${PROJECT_NAME}
ln -snf ${VERSION_PATH} ${PATH_WWW}/${PROJECT_NAME}

#只保留最新的2个文件夹，其余的删除
if [ -d "${PROJECT_PATH}" ]
then
cd ${PROJECT_PATH}
rm -rf `ls -t -d */ |tail -n +3`
fi

chmod -R 777 ${FILE_PATH}
#只保留最新的2个压缩文件，其余的删除
if [ -d "${FILE_PATH}" ]
then
cd ${FILE_PATH}
rm -rf `ls -t |tail -n +3`
fi

#代码备份目录
BAK_PATH=${GO_PATH}/bak/go_project_template
chmod -R 777 ${BAK_PATH}
#只保留最新的5个压缩文件，其余的删除
if [ -d "${BAK_PATH}" ]
then
cd ${BAK_PATH}
rm -rf `ls -t |tail -n +5`
fi

cd ${GO_PATH}
# rm -rf xcc_mall_admin.bak
mv ${PROJECT_NAME}/ bak/${PROJECT_NAME}/${PROJECT_NAME}.${BAKTIME}

mkdir -p ${PROJECT_NAME}
# rm -rf ${GO_PATH}/target/

cd ${PROJECT_NAME}
mkdir -p scripts/log
mkdir -p target/runtime_log
cd ${PATH_WWW}
# cp -r target/* ${GO_PATH}/target
cp -r ${PROJECT_NAME}/* ${GO_PATH}/${PROJECT_NAME}
cd ${GO_PATH}/${PROJECT_NAME}

make stop
sleep 1
make run

# nohup ${GO_PATH}/reload.sh 2>&1 > target/log/reload.log &
