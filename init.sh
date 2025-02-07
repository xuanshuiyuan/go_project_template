#!/bin/bash

# 提示用户输入要替换的新字符串
read -p "请输入项目名称: " project_name

if [ -z "$project_name" ]; then
  echo "项目名称不能为空"
  exit 0
fi

# 要过滤的目录列表
exclude_dirs=(".git") # 将此处替换为你想要过滤的目录名

# 文件名中包含的指定字符串及替换后的新字符串
filename_search_string="go_project_template"
filename_replace_string=$project_name

# 文件内容中要替换的字符串
content_search_string="go_project_template"
content_replace_string=$project_name

# 使用sed命令替换文件中的字符串
# -i选项表示直接修改文件内容
# s/old_string/new_string/g表示将所有的old_string替换为new_string
#如果你在 macOS 上遇到此问题，请尝试使用 sed -i '' 's/old/new/g' filename（注意两个单引号之间的空字符串）来避免备份扩展名。

#修改配置文件
read -p "请输入端口号: " port
if [ -z "$port" ]; then
  echo "端口号不能为空"
  exit 0
else
  sed -i '' "s|15000|$port|g" cmd/conf/local/logic.toml
fi
#echo "配置文件更新完成"
##

# 遍历当前目录和子目录的函数
function traverse_dir() {
  local cur_dir=$1
  for element in "$cur_dir"/*; do
    dir_or_file=$element
    # 检查是否为目录
    if [ -d "$dir_or_file" ]; then
      # 检查是否在排除列表中
      exclude=false
      for exclude_dir in "${exclude_dirs[@]}"; do
        if [[ "$dir_or_file" == *"/$exclude_dir"* ]]; then
          exclude=true
          break
        fi
      done

      # 如果不在排除列表中，则递归遍历子目录
      if [ "$exclude" = false ]; then
        traverse_dir "$dir_or_file"
      fi
    else
      # 处理文件
      # 检查文件名是否包含指定字符串
      if [[ "$dir_or_file" == *"$filename_search_string"* ]]; then
        # 替换文件名中的字符串
        new_name=$(echo "$dir_or_file" | sed "s|$filename_search_string|$filename_replace_string|g")
        mv "$dir_or_file" "$new_name"
        dir_or_file=$new_name # 更新文件名变量，以便后续操作
      fi

      if [[ "$dir_or_file" == *init.sh ]]; then
        echo "$dir_or_file"
      else
        echo "$dir_or_file"
        #      dir_or_file_=${dir_or_file:2}
        # 替换文件内容中的字符串
        sed -i '' "s|$content_search_string|$content_replace_string|g" "$dir_or_file"
      fi
    fi
  done
}

# 从当前目录开始遍历
traverse_dir "."

cp cmd/conf/local/logic.toml cmd/conf/develop/logic.toml
cp cmd/conf/local/logic.toml cmd/conf/production/logic.toml

echo "初始化完成"
