#!/usr/bin/env bash

export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
stty erase ^?

cd "$(
  cd "$(dirname "$0")" || exit
  pwd
)" || exit

# 字体颜色配置
Green="\033[32m"
Red="\033[31m"
Yellow="\033[33m"
Blue="\033[36m"
Font="\033[0m"
GreenBG="\033[42;37m"
RedBG="\033[41;37m"
OK="${Green}[OK]${Font}"
ERROR="${Red}[ERROR]${Font}"

# 变量
xray_version="v25.3.6"
project_version=""  # 由用户输入
github_repo="gongshen/xray"
stat_dir="/usr/local/bin/stat"
xray_admin_dir="/usr/local/bin/xray-admin"
xray_admin_conf_dir="/usr/local/etc/xray-admin"
xray_admin_service_dir="/etc/systemd/system/xray_admin.service"
stat_service_dir="/etc/systemd/system/stat.service"
xray_conf_dir="/usr/local/etc/xray"
iptables_conf_dir="/usr/local/etc/xray/iptables"

# GitHub 下载地址 (在获取版本号后设置)
github_release_url=""

# 获取项目版本号
function get_project_version() {
  if [ -z "$project_version" ]; then
    echo -e "${Blue}请输入要安装的版本号 (如 v1.0.0)${Font}"
    echo -e "${Yellow}提示: 可在 https://github.com/${github_repo}/releases 查看可用版本${Font}"
    read -rp "版本号: " project_version
    if [ -z "$project_version" ]; then
      print_error "版本号不能为空"
      return 1
    fi
    github_release_url="https://github.com/${github_repo}/releases/download/${project_version}"
  fi
  return 0
}

function print_ok() {
  echo -e "${OK} ${Blue} $1 ${Font}"
}

function print_error() {
  echo -e "${ERROR} ${RedBG} $1 ${Font}"
}

function is_root() {
  if [[ 0 == "$UID" ]]; then
    print_ok "当前用户是 root 用户，开始安装流程"
  else
    print_error "当前用户不是 root 用户，请切换到 root 用户后重新执行脚本"
    exit 1
  fi
}

judge() {
  if [[ 0 -eq $? ]]; then
    print_ok "$1 完成"
    sleep 1
  else
    print_error "$1 失败"
    exit 1
  fi
}

function system_check() {
  source '/etc/os-release'

  if [[ "${ID}" == "centos" && ${VERSION_ID} -ge 7 ]]; then
    print_ok "当前系统为 Centos ${VERSION_ID} ${VERSION}"
    INS="yum install -y"
  elif [[ "${ID}" == "ol" ]]; then
    print_ok "当前系统为 Oracle Linux ${VERSION_ID} ${VERSION}"
    INS="yum install -y"
  elif [[ "${ID}" == "debian" && ${VERSION_ID} -ge 9 ]]; then
    print_ok "当前系统为 Debian ${VERSION_ID} ${VERSION}"
    INS="apt install -y"
    apt update
  elif [[ "${ID}" == "ubuntu" && $(echo "${VERSION_ID}" | cut -d '.' -f1) -ge 18 ]]; then
    print_ok "当前系统为 Ubuntu ${VERSION_ID} ${UBUNTU_CODENAME}"
    INS="apt install -y"
    apt update
  else
    print_error "当前系统为 ${ID} ${VERSION_ID} 不在支持的系统列表内"
    exit 1
  fi

  $INS dbus

  # 关闭各类防火墙
  systemctl stop firewalld 2>/dev/null
  systemctl disable firewalld 2>/dev/null
  systemctl stop nftables 2>/dev/null
  systemctl disable nftables 2>/dev/null
  systemctl stop ufw 2>/dev/null
  systemctl disable ufw 2>/dev/null
}

function dependency_install() {
  ${INS} lsof tar wget curl unzip jq
  judge "安装基础依赖"

  if [[ "${ID}" == "centos" || "${ID}" == "ol" ]]; then
    ${INS} crontabs iptables-services
    touch /var/spool/cron/root && chmod 600 /var/spool/cron/root
    systemctl start crond && systemctl enable crond
  else
    ${INS} cron iptables-persistent
    touch /var/spool/cron/crontabs/root && chmod 600 /var/spool/cron/crontabs/root
    systemctl start cron && systemctl enable cron
  fi
  judge "crontab 自启动配置"

  ${INS} systemd
  judge "安装/升级 systemd"

  mkdir -p /usr/local/bin >/dev/null 2>&1
  mkdir -p ${xray_conf_dir} >/dev/null 2>&1
}

function basic_optimization() {
  sed -i '/^\*\ *soft\ *nofile\ *[[:digit:]]*/d' /etc/security/limits.conf
  sed -i '/^\*\ *hard\ *nofile\ *[[:digit:]]*/d' /etc/security/limits.conf
  echo '* soft nofile 65536' >>/etc/security/limits.conf
  echo '* hard nofile 65536' >>/etc/security/limits.conf

  if [[ "${ID}" == "centos" || "${ID}" == "ol" ]]; then
    sed -i 's/^SELINUX=.*/SELINUX=disabled/' /etc/selinux/config
    setenforce 0 2>/dev/null
  fi
}

# ==================== 通用函数 ====================

# 检查文件是否存在，如果存在则询问是否覆盖
# 参数: $1 = 文件路径, $2 = 文件描述
# 返回: 0 = 继续创建, 1 = 跳过创建
function confirm_overwrite() {
  local file_path=$1
  local file_desc=$2
  
  if [ -f "$file_path" ]; then
    echo -e "${Yellow}[警告]${Font} ${file_desc} 已存在: ${file_path}"
    read -rp "是否覆盖? (y/n): " overwrite
    if [[ "$overwrite" != "y" && "$overwrite" != "Y" ]]; then
      print_ok "跳过 ${file_desc}"
      return 1
    fi
  fi
  return 0
}

# ==================== 下载函数 ====================

function download_stat() {
  print_ok "从 GitHub 下载 stat..."
  
  # 获取版本号
  get_project_version || return 1
  
  # 先停止正在运行的服务，避免 "Text file busy" 错误
  if systemctl is-active --quiet stat 2>/dev/null; then
    print_ok "停止 stat 服务..."
    systemctl stop stat
    sleep 1
  fi
  
  # 检查是否需要覆盖
  confirm_overwrite "${stat_dir}" "stat 二进制文件" || return 0
  
  wget -O ${stat_dir} ${github_release_url}/stat
  if [ $? -ne 0 ]; then
    print_error "下载 stat 失败"
    return 1
  fi
  chmod +x ${stat_dir}
  print_ok "stat 下载完成"
}

function download_xray_admin() {
  print_ok "从 GitHub 下载 xray-admin..."
  
  # 获取版本号
  get_project_version || return 1
  
  # 检查是否需要覆盖
  confirm_overwrite "${xray_admin_dir}" "xray-admin 二进制文件" || return 0
  
  # 停止正在运行的服务，避免 "Text file busy" 错误
  systemctl stop xray_admin 2>/dev/null
  
  wget -O ${xray_admin_dir} ${github_release_url}/xray-admin
  if [ $? -ne 0 ]; then
    print_error "下载 xray-admin 失败"
    return 1
  fi
  chmod +x ${xray_admin_dir}
  print_ok "xray-admin 下载完成"
}

# ==================== 配置文件创建函数 ====================

function create_xray_config() {
  mkdir -p ${xray_conf_dir}
  mkdir -p /var/log/xray
  
  # 检查是否需要覆盖
  confirm_overwrite "${xray_conf_dir}/config.json" "Xray 配置文件" || return 0
  
  cat > ${xray_conf_dir}/config.json << 'XRAYEOF'
{
  "policy": {
    "levels": {
      "0": {
        "statsUserUplink": true,
        "statsUserDownlink": true
      }
    },
    "system": {
      "statsInboundUplink": true,
      "statsInboundDownlink": true,
      "statsOutboundUplink": true,
      "statsOutboundDownlink": true
    }
  },
  "stats": {},
  "log": {
    "access": "/var/log/xray/access.log",
    "error": "/var/log/xray/error.log",
    "loglevel": "warning"
  },
  "routing": {
    "domainStrategy": "AsIs",
    "rules": [
      {
        "type": "field",
        "inboundTag": ["api"],
        "outboundTag": "api"
      }
    ]
  },
  "inbounds": [
    {
      "listen": "0.0.0.0",
      "port": 80,
      "protocol": "vmess",
      "settings": {
        "clients": [
          {
            "id": "",
            "alterId": 64,
            "level": 0,
            "email": "admin"
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "none",
        "tcpSettings": {
          "header": {
            "type": "http",
            "request": {
              "method": "GET",
              "path": ["/"],
              "headers": {}
            },
            "response": {
              "version": "1.1",
              "status": "200",
              "reason": "OK",
              "headers": {}
            }
          }
        }
      }
    },
    {
      "listen": "127.0.0.1",
      "port": 11111,
      "protocol": "dokodemo-door",
      "settings": {
        "address": "127.0.0.1"
      },
      "tag": "api"
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    }
  ],
  "api": {
    "tag": "api",
    "services": ["StatsService"]
  }
}
XRAYEOF

  # 生成UUID并写入配置
  UUID=$(cat /proc/sys/kernel/random/uuid)
  sed -i "s/\"id\": \"\"/\"id\": \"${UUID}\"/" ${xray_conf_dir}/config.json
  
  print_ok "Xray 配置文件已创建"
  echo -e "${Blue}UUID: ${UUID}${Font}"
}

function create_stat_service() {
  # 检查是否需要覆盖
  confirm_overwrite "${stat_service_dir}" "stat.service 文件" || return 0
  
  cat > ${stat_service_dir} << 'STATEOF'
[Unit]
Description=Stat Service
After=network.target nss-lookup.target

[Service]
User=root
Environment="REMOTE_IP=__REMOTE_IP__"
ExecStart=/usr/local/bin/stat -port __STAT_PORT__
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target
STATEOF
  print_ok "stat.service 文件已创建"
}

function create_xray_admin_service() {
  # 检查是否需要覆盖
  confirm_overwrite "${xray_admin_service_dir}" "xray_admin.service 文件" || return 0
  
  cat > ${xray_admin_service_dir} << 'ADMINSERVICEEOF'
[Unit]
Description=xray_admin Service
After=network.target nss-lookup.target

[Service]
User=root
WorkingDirectory=/root
ExecStart=/usr/local/bin/xray-admin -c /usr/local/etc/xray-admin/config.yaml
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target
ADMINSERVICEEOF
  print_ok "xray_admin.service 文件已创建"
}

function create_xray_admin_config() {
  mkdir -p ${xray_admin_conf_dir}
  
  # 检查是否需要覆盖
  confirm_overwrite "${xray_admin_conf_dir}/config.yaml" "xray_admin 配置文件" || return 0
  
  cat > ${xray_admin_conf_dir}/config.yaml << 'ADMINCONFIGEOF'
jwt:
  signing-key: qmPlus
  expires-time: 7d
  buffer-time: 1d
  issuer: qmPlus

zap:
  level: info
  format: console
  prefix: "[xray-admin]"
  director: log
  show-line: true
  encode-level: LowercaseColorLevelEncoder
  stacktrace-key: stacktrace
  log-in-console: true

system:
  env: public
  addr: 8888
  db-type: mysql
  oss-type: local
  use-redis: false
  use-multipoint: false
  iplimit-count: 500
  iplimit-time: 3600
  router-prefix: ""

captcha:
  key-long: 6
  img-width: 240
  img-height: 80
  open-captcha: 0
  open-captcha-timeout: 3600

mysql:
  path: ""
  port: ""
  config: ""
  db-name: ""
  username: ""
  password: ""
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: ""
  log-zap: false

local:
  path: uploads/file
  store-path: uploads/file

stat_port: 56611
ADMINCONFIGEOF
  print_ok "xray_admin config.yaml 文件已创建"
}

function create_iptables_config() {
  mkdir -p ${xray_conf_dir}
  
  # 检查是否需要覆盖
  confirm_overwrite "${iptables_conf_dir}" "iptables 配置文件" || return 0
  
  cat > ${iptables_conf_dir} << 'IPTABLESEOF'
*filter
:INPUT DROP [0:0]
:FORWARD DROP [0:0]
:OUTPUT ACCEPT [0:0]

-A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
-A INPUT -i lo -j ACCEPT
-A INPUT -p tcp -m state --state NEW -m tcp --dport __SSHPORT__ -j ACCEPT
###-A INPUT -p tcp -m state --state NEW -m tcp --dport __ADMINPORT__ -j ACCEPT
-A INPUT -p tcp -s __ADMINIP__ --dport __STATPORT__ -j ACCEPT
-A INPUT -p tcp --dport __XRAYPORT__ -j ACCEPT

COMMIT
IPTABLESEOF
  print_ok "iptables 配置文件已创建"
}

# ==================== 安装函数 ====================

function xray_install() {
  print_ok "安装 Xray"
  curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh | bash -s -- install --version ${xray_version}
  judge "Xray 安装"
}

function configure_xray() {
  create_xray_config
  systemctl enable xray
  systemctl restart xray
  judge "Xray 启动"
}

function install_xray() {
  is_root
  system_check
  dependency_install
  basic_optimization
  xray_install
  configure_xray
}

function bbr_boost_sh() {
  [ -f "tcp.sh" ] && rm -rf ./tcp.sh
  wget -N --no-check-certificate "https://raw.githubusercontent.com/ylx2016/Linux-NetSpeed/master/tcp.sh" && chmod +x tcp.sh && ./tcp.sh
}

function install_stat() {
  print_ok "安装 Stat 服务"
  
  # 从 GitHub 下载 stat
  download_stat || return 1
  
  # 创建service文件
  create_stat_service
  
  read -rp "请输入管理端IP地址：" remoteIp
  sed -i "s|__REMOTE_IP__|${remoteIp}|" ${stat_service_dir}
  
  read -rp "请输入Stat监听端口(默认56611)：" statPort
  [ -z "$statPort" ] && statPort="56611"
  sed -i "s|__STAT_PORT__|${statPort}|" ${stat_service_dir}
  
  systemctl daemon-reload
  systemctl enable stat
  systemctl restart stat
  judge "Stat 启动"
}

function install_admin() {
  print_ok "安装 xray-admin 管理端"
  
  # 从 GitHub 下载 xray-admin
  download_xray_admin || return 1
  
  # 创建配置文件和service文件
  create_xray_admin_config
  create_xray_admin_service
  
  systemctl daemon-reload
  systemctl enable xray_admin
  systemctl restart xray_admin
  judge "xray_admin 启动"
}

function install_mysql() {
  print_ok "安装 MySQL 8.0"
  
  if [[ "${ID}" == "centos" || "${ID}" == "ol" ]]; then
    yum update -y
    rpm -Uvh https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm
    rpm --import https://repo.mysql.com/RPM-GPG-KEY-mysql-2022
    yum install mysql-server -y
  else
    ${INS} mysql-server
  fi
  
  systemctl start mysqld
  systemctl enable mysqld
  
  MYSQL_INIT_PASSWORD=$(grep 'temporary password' /var/log/mysqld.log 2>/dev/null | awk '{print $NF}')
  
  if [ -n "$MYSQL_INIT_PASSWORD" ]; then
    echo -e "${Yellow}MySQL 初始密码: ${MYSQL_INIT_PASSWORD}${Font}"
    read -rp "请输入新的MySQL密码: " NEW_MYSQL_PASSWORD
    mysql --connect-expired-password -uroot -p"$MYSQL_INIT_PASSWORD" <<EOF
ALTER USER 'root'@'localhost' IDENTIFIED BY '$NEW_MYSQL_PASSWORD';
EOF
    print_ok "MySQL 密码修改完成"
  fi
  
  judge "MySQL 安装"
}

function init_admin_db() {
  read -rp "请输入数据库类型(默认：mysql)：" ADMIN_DB_TYPE
  [ -z "$ADMIN_DB_TYPE" ] && ADMIN_DB_TYPE="mysql"
  read -rp "请输入数据库地址(默认：127.0.0.1)：" ADMIN_DB_HOST
  [ -z "$ADMIN_DB_HOST" ] && ADMIN_DB_HOST="127.0.0.1"
  read -rp "请输入数据库端口(默认：3306)：" ADMIN_DB_PORT
  [ -z "$ADMIN_DB_PORT" ] && ADMIN_DB_PORT="3306"
  read -rp "请输入数据库用户名(默认：root)：" ADMIN_DB_USERNAME
  [ -z "$ADMIN_DB_USERNAME" ] && ADMIN_DB_USERNAME="root"
  read -rp "请输入数据库密码：" ADMIN_DB_PASSWORD
  read -rp "请输入数据库库名(默认：gva)：" ADMIN_DB_NAME
  [ -z "$ADMIN_DB_NAME" ] && ADMIN_DB_NAME="gva"
  
  headers='Content-Type: application/json'
  data="{\"dbType\": \"$ADMIN_DB_TYPE\",\"host\": \"$ADMIN_DB_HOST\",\"port\": \"$ADMIN_DB_PORT\", \"userName\": \"$ADMIN_DB_USERNAME\", \"password\": \"$ADMIN_DB_PASSWORD\", \"dbName\": \"$ADMIN_DB_NAME\"}"
  curl -X POST -H "$headers" -d "$data" http://127.0.0.1:8888/init/initdb
  judge "数据库初始化"
}

function config_iptables() {
  print_ok "配置 iptables 防火墙"
  
  # 创建iptables配置文件
  create_iptables_config
  
  # SSH端口
  while true; do
    read -rp "请输入SSH端口号：" SSHPORT
    if [[ -n "$SSHPORT" ]]; then
      break
    else
      print_error "SSH端口不能为空，请重新输入"
    fi
  done
  sed -i "s|__SSHPORT__|${SSHPORT}|" ${iptables_conf_dir}
  
  # 管理端IP
  read -rp "请输入管理端IP地址：" remoteIp
  sed -i "s|__ADMINIP__|${remoteIp}|" ${iptables_conf_dir}
  
  # Xray端口
  read -rp "请输入Xray端口(默认80)：" xray_port
  [ -z "$xray_port" ] && xray_port="80"
  sed -i "s|__XRAYPORT__|${xray_port}|" ${iptables_conf_dir}
  
  # Stat端口
  read -rp "请输入Stat端口(默认56611)：" stat_port
  [ -z "$stat_port" ] && stat_port="56611"
  sed -i "s|__STATPORT__|${stat_port}|" ${iptables_conf_dir}
  
  # 是否是管理员服务器
  read -rp "是否是管理员服务器? (y/n): " is_admin
  if [[ "$is_admin" == "y" || "$is_admin" == "Y" ]]; then
    read -rp "请输入管理员端口号(默认8888)：" ADMINPORT
    [ -z "$ADMINPORT" ] && ADMINPORT="8888"
    sed -i "s|###||" ${iptables_conf_dir}
    sed -i "s|__ADMINPORT__|${ADMINPORT}|" ${iptables_conf_dir}
  fi
  
  iptables-restore < ${iptables_conf_dir}
  
  # 保存iptables规则
  if [[ "${ID}" == "centos" || "${ID}" == "ol" ]]; then
    service iptables save
  else
    netfilter-persistent save
  fi
  
  judge "iptables 配置"
}

# ==================== ACME 证书相关 ====================

function install_acme() {
  print_ok "安装 acme.sh"
  
  # 安装依赖
  ${INS} curl socat
  judge "安装 acme 依赖"
  
  # 安装 acme.sh
  curl https://get.acme.sh | sh -s email=admin@example.com
  judge "安装 acme.sh"
  
  # 设置默认CA
  ~/.acme.sh/acme.sh --set-default-ca --server letsencrypt
  
  print_ok "acme.sh 安装完成"
  echo ""
  echo -e "${Blue}使用方法:${Font}"
  echo "1. 申请证书: ~/.acme.sh/acme.sh --issue -d your-domain.com --standalone"
  echo "2. 续期证书: 选择菜单 17"
}

function renew_cert() {
  if [ ! -f ~/.acme.sh/acme.sh ]; then
    print_error "acme.sh 未安装，请先选择菜单 16 安装"
    return 1
  fi
  
  read -rp "请输入要续期的域名: " DOMAIN
  if [ -z "$DOMAIN" ]; then
    print_error "域名不能为空"
    return 1
  fi
  
  print_ok "开始续期证书: ${DOMAIN}"
  
  ~/.acme.sh/acme.sh --renew -d "${DOMAIN}" --force
  
  if [ $? -eq 0 ]; then
    print_ok "证书续期成功: ${DOMAIN}"
    
    read -rp "是否重启 xray 服务? (y/n): " restart_service
    if [[ "$restart_service" == "y" || "$restart_service" == "Y" ]]; then
      systemctl restart xray 2>/dev/null && print_ok "xray 重启完成"
    fi
  else
    print_error "证书续期失败"
    echo "可能需要重新申请: ~/.acme.sh/acme.sh --issue -d ${DOMAIN} --standalone --force"
  fi
}

# ==================== 菜单 ====================

menu() {
  echo -e "\t ${Green}Xray 安装管理脚本${Font}"
  echo -e "—————————————— 安装向导 ——————————————"
  echo -e "${Green}1.${Font}  安装 Xray (含依赖和优化)"
  echo -e "${Green}2.${Font}  配置 Xray"
  echo -e "${Green}3.${Font}  安装 BBR 加速脚本"
  echo -e "${Green}4.${Font}  安装 Stat 服务"
  echo -e "${Green}5.${Font}  安装 xray-admin 管理端"
  echo -e "${Green}6.${Font}  安装 MySQL"
  echo -e "${Green}7.${Font}  初始化管理端数据库"
  echo -e "${Green}8.${Font}  配置 iptables 防火墙"
  echo -e "${Green}9.${Font}  安装 acme.sh (SSL证书工具)"
  echo -e "${Green}10.${Font} 续期 SSL 证书"
  echo -e "${Green}99.${Font} 退出"
  echo -e "————————————————————————————————————"
  read -rp "请输入数字：" menu_num
  case $menu_num in
  1)
    install_xray
    ;;
  2)
    configure_xray
    ;;
  3)
    bbr_boost_sh
    ;;
  4)
    install_stat
    ;;
  5)
    install_admin
    ;;
  6)
    install_mysql
    ;;
  7)
    init_admin_db
    ;;
  8)
    config_iptables
    ;;
  9)
    install_acme
    ;;
  10)
    renew_cert
    ;;
  99)
    exit 0
    ;;
  *)
    print_error "请输入正确的数字"
    menu
    ;;
  esac
}

menu "$@"
