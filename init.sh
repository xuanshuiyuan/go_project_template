#!/bin/bash

# 提示用户输入要替换的新字符串
read -p "请输入项目名称: " project_name

if [ -z "$project_name" ]; then
  echo "项目名称不能为空"
  exit 0
fi

# 使用sed命令替换文件中的字符串
# -i选项表示直接修改文件内容
# s/old_string/new_string/g表示将所有的old_string替换为new_string
#如果你在 macOS 上遇到此问题，请尝试使用 sed -i '' 's/old/new/g' filename（注意两个单引号之间的空字符串）来避免备份扩展名。

#更新Makefile里面的项目名称
#echo "1：Makefile 更新中..."
#sed -i '' "s|go_project_template|$project_name|g" Makefile
#echo "Makefile 更新完成"
###

#修改配置文件
#echo "1：配置文件更新中..."
#sed -i '' "s|go_project_template|$project_name|g" cmd/conf/local/logic.toml
#read -p "请输入端口号: " port
#if [ -z "$port" ]; then
#  echo "端口号不能为空"
#  exit 0
#else
#  sed -i '' "s|15000|$port|g" cmd/conf/local/logic.toml
#fi
#cp cmd/conf/local/logic.toml cmd/conf/develop/logic.toml
#cp cmd/conf/local/logic.toml cmd/conf/production/logic.toml
#echo "配置文件更新完成"
###
