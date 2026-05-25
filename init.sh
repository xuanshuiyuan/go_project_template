#!/bin/bash

set -e

read -p "请输入项目名称: " project_name
if [ -z "$project_name" ]; then
  echo "错误: 项目名称不能为空" >&2
  exit 1
fi

old_project_name="go_project_template"
old_port=15000

read -p "请输入端口号: " port
if [ -z "$port" ]; then
  echo "错误: 端口号不能为空" >&2
  exit 1
fi

# 修改端口号 - 兼容 Linux 和 macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' "s|$old_port|$port|g" cmd/conf/local/logic.toml
else
  sed -i "s|$old_port|$port|g" cmd/conf/local/logic.toml
fi

# 用 find + sed 替换手动递归遍历，排除 .git 目录和 init.sh 自身
# 替换文件内容
find . -type f \
  -not -path "./.git/*" \
  -not -name "init.sh" \
  -exec grep -l "$old_project_name" {} \; 2>/dev/null | while read -r file; do
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s|$old_project_name|$project_name|g" "$file"
  else
    sed -i "s|$old_project_name|$project_name|g" "$file"
  fi
done

# 重命名包含旧项目名的文件和目录
find . -depth -name "*$old_project_name*" -not -path "./.git/*" | while read -r item; do
  new_item=$(echo "$item" | sed "s|$old_project_name|$project_name|g")
  mv "$item" "$new_item"
done

# 同步配置文件到其他环境
cp cmd/conf/local/logic.toml cmd/conf/develop/logic.toml
cp cmd/conf/local/logic.toml cmd/conf/production/logic.toml
cp cmd/conf/local/logic.toml target/logic.toml

echo "初始化完成"
