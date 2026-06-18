#!/usr/bin/env bash

export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
if [ -t 0 ]; then
  stty erase ^?
fi

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
stat_data_dir="/var/lib/xray-stat"
stat_db_path="${stat_data_dir}/stat.db"
stat_log_level="info"
stat_collect_interval="10s"
stat_traffic_retention_months="12"
stat_api_traffic_tag="1"
stat_log_clean_dir="/root/log"
stat_log_retention_months="12"
xray_conf_dir="/usr/local/etc/xray"
iptables_conf_dir="/usr/local/etc/xray/iptables"
xray_logrotate_conf_dir="/etc/logrotate.d/xray"
xray_log_dir="/var/log/xray"
stat_xray_log_dir="/var/log/xray"
stat_xray_log_retention_months="12"
service_time_zone="Asia/Shanghai"
xray_timezone_conf_dir="/etc/systemd/system/xray.service.d"
xray_timezone_conf_file="${xray_timezone_conf_dir}/timezone.conf"
sqlite_utc_to_shanghai_modifier="+8 hours"
sqlite_shanghai_to_utc_modifier="-8 hours"

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
  ${INS} lsof tar wget curl unzip jq logrotate
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
  mkdir -p ${stat_data_dir} >/dev/null 2>&1
  chmod 755 ${stat_data_dir} >/dev/null 2>&1
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

function create_xray_logrotate_config() {
  local output_file=${1:-${xray_logrotate_conf_dir}}
  mkdir -p "$(dirname "${output_file}")"

  cat > "${output_file}" <<'LOGROTATEEOF'
/var/log/xray/access.log /var/log/xray/error.log {
    daily
    rotate 365
    maxage 365
    missingok
    notifempty
    compress
    delaycompress
    copytruncate
    dateext
    su root root
}
LOGROTATEEOF

  print_ok "Xray 日志轮转配置已创建: ${output_file}"
}

function ensure_logrotate_installed() {
  if command -v logrotate >/dev/null 2>&1; then
    return 0
  fi

  if [[ "${XRAY_INSTALL_TEST_MODE:-0}" == "1" ]]; then
    return 0
  fi

  if command -v apt-get >/dev/null 2>&1; then
    apt-get update
    apt-get install -y logrotate
  elif command -v yum >/dev/null 2>&1; then
    yum install -y logrotate
  elif command -v dnf >/dev/null 2>&1; then
    dnf install -y logrotate
  else
    print_error "未找到可用的包管理器，请先手动安装 logrotate"
    return 1
  fi
  judge "安装 logrotate"
}

function configure_xray_logrotate() {
  if [[ "${XRAY_INSTALL_TEST_MODE:-0}" != "1" ]]; then
    is_root
  fi
  ensure_logrotate_installed
  create_xray_logrotate_config
}

function create_xray_timezone_config() {
  local output_file=${1:-${xray_timezone_conf_file}}
  mkdir -p "$(dirname "${output_file}")"

  cat > "${output_file}" <<XRAYTIMEZONEEOF
[Service]
Environment="TZ=${service_time_zone}"
XRAYTIMEZONEEOF

  print_ok "Xray systemd 时区配置已创建: ${output_file}"
}

function configure_xray_timezone() {
  if [[ "${XRAY_INSTALL_TEST_MODE:-0}" != "1" ]]; then
    is_root
  fi
  create_xray_timezone_config
  if [[ "${XRAY_INSTALL_TEST_MODE:-0}" != "1" ]]; then
    systemctl daemon-reload
  fi
}

function ensure_sqlite3_installed() {
  if command -v sqlite3 >/dev/null 2>&1; then
    return 0
  fi

  if [[ "${XRAY_INSTALL_TEST_MODE:-0}" == "1" ]]; then
    return 0
  fi

  if command -v apt-get >/dev/null 2>&1; then
    apt-get update
    apt-get install -y sqlite3
  elif command -v yum >/dev/null 2>&1; then
    yum install -y sqlite
  elif command -v dnf >/dev/null 2>&1; then
    dnf install -y sqlite
  else
    print_error "未找到可用的包管理器，请先手动安装 sqlite3"
    return 1
  fi
  judge "安装 sqlite3"
}

function sqlite_quote() {
  local value=${1:-}
  value=${value//\'/\'\'}
  printf "'%s'" "${value}"
}

function normalize_positive_int() {
  local value=${1:-}
  local default_value=${2:-50}
  local max_value=${3:-1000}
  if [[ ! "${value}" =~ ^[0-9]+$ ]] || [[ "${value}" -le 0 ]]; then
    value="${default_value}"
  fi
  if [[ "${value}" -gt "${max_value}" ]]; then
    value="${max_value}"
  fi
  printf "%s" "${value}"
}

function normalize_analysis_date() {
  local analysis_date=${1:-}
  if [[ "${analysis_date}" =~ ^[0-9]{8}$ ]]; then
    analysis_date="${analysis_date:0:4}-${analysis_date:4:2}-${analysis_date:6:2}"
  fi
  if [[ ! "${analysis_date}" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}$ ]]; then
    echo "日期格式无效，请输入 YYYYMMDD，例如 20260617"
    return 1
  fi

  local normalized_date
  if ! normalized_date=$(date -d "${analysis_date}" "+%Y-%m-%d" 2>/dev/null); then
    echo "日期无效: ${analysis_date}"
    return 1
  fi
  if [[ "${normalized_date}" != "${analysis_date}" ]]; then
    echo "日期无效: ${analysis_date}"
    return 1
  fi

  printf "%s" "${analysis_date}"
}

function normalize_analysis_clock() {
  local clock_value=${1:-}
  local boundary=${2:-start}
  if [[ ! "${clock_value}" =~ ^([0-9]{1,2}):([0-9]{2})$ ]]; then
    echo "时间格式无效，请输入 小时:分钟，例如 8:10 或 09:00"
    return 1
  fi

  local hour=$((10#${BASH_REMATCH[1]}))
  local minute=$((10#${BASH_REMATCH[2]}))
  if (( hour < 0 || hour > 23 || minute < 0 || minute > 59 )); then
    echo "时间无效: ${clock_value}"
    return 1
  fi

  local second="00"
  if [[ "${boundary}" == "end" ]]; then
    second="59"
  fi
  printf "%02d:%02d:%s" "${hour}" "${minute}" "${second}"
}

function normalized_clock_to_minutes() {
  local clock_value=${1:-}
  local hour=$((10#${clock_value:0:2}))
  local minute=$((10#${clock_value:3:2}))
  printf "%s" $((hour * 60 + minute))
}

function analysis_datetime_range_from_parts() {
  local analysis_date
  local start_clock
  local end_clock
  if ! analysis_date=$(normalize_analysis_date "${1:-}"); then
    echo "${analysis_date}"
    return 1
  fi
  if ! start_clock=$(normalize_analysis_clock "${2:-}" "start"); then
    echo "${start_clock}"
    return 1
  fi
  if ! end_clock=$(normalize_analysis_clock "${3:-}" "end"); then
    echo "${end_clock}"
    return 1
  fi

  local start_time="${analysis_date} ${start_clock}"
  local end_time="${analysis_date} ${end_clock}"
  if [[ "${start_time}" > "${end_time}" ]]; then
    echo "结束时间不能早于开始时间"
    return 1
  fi

  local start_minutes
  local end_minutes
  start_minutes=$(normalized_clock_to_minutes "${start_clock}")
  end_minutes=$(normalized_clock_to_minutes "${end_clock}")
  if (( end_minutes - start_minutes > 120 )); then
    echo "开始时间和结束时间不能超过 2 小时"
    return 1
  fi

  printf "%s|%s" "${start_time}" "${end_time}"
}

function build_traffic_event_where() {
  local user_tag=$1
  local start_time=${2:-}
  local end_time=${3:-}
  local where="tag = $(sqlite_quote "${user_tag}")"
  if [[ -n "${start_time}" ]]; then
    where="${where} AND collected_at >= strftime('%s', $(sqlite_quote "${start_time}"), '${sqlite_shanghai_to_utc_modifier}')"
  fi
  if [[ -n "${end_time}" ]]; then
    where="${where} AND collected_at <= strftime('%s', $(sqlite_quote "${end_time}"), '${sqlite_shanghai_to_utc_modifier}')"
  fi
  printf "%s" "${where}"
}

function build_traffic_minute_summary_sql() {
  local where
  where=$(build_traffic_event_where "$1" "${2:-}" "${3:-}")
  cat <<SQL
SELECT
  strftime('%Y-%m-%d %H:%M', collected_at, 'unixepoch', '${sqlite_utc_to_shanghai_modifier}') AS minute,
  COUNT(*) AS events,
  COALESCE(SUM(down), 0) AS down,
  COALESCE(SUM(up), 0) AS up,
  COALESCE(SUM(down + up), 0) AS total
FROM traffic_event
WHERE ${where}
GROUP BY minute
ORDER BY minute ASC;
SQL
}

function format_bytes() {
  local bytes=${1:-0}
  awk -v bytes="${bytes}" 'BEGIN {
    split("B KiB MiB GiB TiB", units, " ");
    value = bytes + 0;
    unit = 1;
    while (value >= 1024 && unit < 5) {
      value = value / 1024;
      unit++;
    }
    if (unit == 1) {
      printf "%d %s", value, units[unit];
    } else {
      printf "%.2f %s", value, units[unit];
    }
  }'
}

function escape_grep_regex() {
  printf "%s" "${1:-}" | sed 's/[][(){}.^$*+?|\\]/\\&/g'
}

function build_access_log_email_pattern() {
  printf "email: %s$" "$(escape_grep_regex "$1")"
}

function summarize_access_log_targets_by_minute() {
  awk '
    function extract_target(line, fields, field_count, i, target) {
      field_count = split(line, fields, /[[:space:]]+/)
      for (i = 1; i <= field_count; i++) {
        if (fields[i] ~ /^(tcp|udp):/) {
          target = fields[i]
          sub(/^(tcp|udp):/, "", target)
          sub(/:[0-9]+$/, "", target)
          return target
        }
      }
      return ""
    }
    function is_filtered_ipv4(target, octets, a, b, c, d) {
      if (target !~ /^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$/) {
        return 0
      }
      split(target, octets, ".")
      a = octets[1] + 0
      b = octets[2] + 0
      c = octets[3] + 0
      d = octets[4] + 0
      return a == 0 || a == 10 || a == 127 || a >= 224 || \
        (a == 172 && b >= 16 && b <= 31) || \
        (a == 192 && b == 168) || \
        (a == 169 && b == 254) || \
        (a == 100 && b >= 64 && b <= 127) || \
        (a == 255 && b == 255 && c == 255 && d == 255)
    }
    function is_filtered_ipv6(target) {
      return target == "::" || target == "::1" || \
        target ~ /^fc/ || target ~ /^fd/ || target ~ /^ff/ || \
        target ~ /^fe[89ab][0-9a-f]*:/
    }
    function normalize_target(target, labels, label_count) {
      target = tolower(target)
      gsub(/^\[/, "", target)
      gsub(/\]$/, "", target)
      sub(/\.$/, "", target)
      if (is_filtered_ipv4(target) || is_filtered_ipv6(target)) {
        return ""
      }
      if (target ~ /^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$/ || target ~ /:/) {
        return target
      }
      label_count = split(target, labels, ".")
      if (label_count >= 2) {
        return labels[label_count - 1] "." labels[label_count]
      }
      return target
    }
    {
      log_minute = substr($0, 1, 16)
      if (log_minute !~ /^[0-9][0-9][0-9][0-9]\/[0-9][0-9]\/[0-9][0-9] [0-9][0-9]:[0-9][0-9]$/) {
        next
      }
      gsub("/", "-", log_minute)
      target = extract_target($0)
      if (target == "") {
        next
      }
      target = normalize_target(target)
      if (target == "") {
        next
      }
      key = log_minute SUBSEP target
      if (!seen_target[key]++) {
        if (target_list[log_minute] == "") {
          target_list[log_minute] = target
        } else {
          target_list[log_minute] = target_list[log_minute] ", " target
        }
      }
      minutes[log_minute] = 1
    }
    END {
      for (minute in minutes) {
        print minute "|" target_list[minute]
      }
    }
  ' | sort
}

function is_xray_access_log_file() {
  local name
  name=$(basename "$1")
  [[ "${name}" == "access.log" ]] && return 0
  [[ "${name}" =~ ^access\.log-([0-9]{8}|[0-9]{4}-[0-9]{2}-[0-9]{2})(\.gz)?$ ]]
}

function find_xray_access_log_files() {
  local log_dir=${1:-${xray_log_dir}}
  [[ -d "${log_dir}" ]] || return 0
  find "${log_dir}" -maxdepth 1 -type f \( -name "access.log" -o -name "access.log-*" \) | while IFS= read -r file; do
    if is_xray_access_log_file "${file}"; then
      printf "%s\n" "${file}"
    fi
  done
}

function normalize_access_log_time() {
  local value=${1:-}
  value=${value//-/\/}
  printf "%s" "${value}"
}

function filter_access_log_time_range() {
  local start_time
  local end_time
  start_time=$(normalize_access_log_time "${1:-}")
  end_time=$(normalize_access_log_time "${2:-}")
  awk -v start_time="${start_time}" -v end_time="${end_time}" '
    {
      log_time = substr($0, 1, 19)
      if (log_time ~ /^[0-9][0-9][0-9][0-9]\/[0-9][0-9]\/[0-9][0-9] [0-9][0-9]:[0-9][0-9]:[0-9][0-9]$/) {
        if (start_time != "" && log_time < start_time) next
        if (end_time != "" && log_time > end_time) next
      }
      print
    }
  '
}

function emit_access_log_matches() {
  local user_tag=$1
  local log_dir=${2:-${xray_log_dir}}
  local pattern
  pattern=$(build_access_log_email_pattern "${user_tag}")
  find_xray_access_log_files "${log_dir}" | sort | while IFS= read -r file; do
    if [[ "${file}" == *.gz ]]; then
      gzip -cd -- "${file}" 2>/dev/null | grep -hE -- "${pattern}" || true
    else
      grep -hE -- "${pattern}" "${file}" 2>/dev/null || true
    fi
  done
}

function collect_traffic_minute_summary() {
  local db_path=$1
  local user_tag=$2
  local start_time=${3:-}
  local end_time=${4:-}
  local output_file=$5

  : > "${output_file}"
  if [[ ! -f "${db_path}" ]]; then
    print_error "SQLite 数据库不存在: ${db_path}，仅展示 access.log 聚合"
    return 0
  fi
  ensure_sqlite3_installed || return 1
  if ! sqlite3 "${db_path}" "SELECT 1 FROM traffic_event LIMIT 1;" >/dev/null 2>&1; then
    print_error "SQLite 中未找到 traffic_event 表: ${db_path}，仅展示 access.log 聚合"
    return 0
  fi

  sqlite3 -noheader -separator "|" "${db_path}" "$(build_traffic_minute_summary_sql "${user_tag}" "${start_time}" "${end_time}")" > "${output_file}"
}

function collect_access_minute_summary() {
  local log_dir=$1
  local user_tag=$2
  local start_time=${3:-}
  local end_time=${4:-}
  local matched_log_file=$5
  local output_file=$6

  : > "${matched_log_file}"
  : > "${output_file}"
  if [[ ! -d "${log_dir}" ]]; then
    print_error "Xray 日志目录不存在: ${log_dir}，仅展示 SQLite 流量聚合"
    return 0
  fi

  emit_access_log_matches "${user_tag}" "${log_dir}" | filter_access_log_time_range "${start_time}" "${end_time}" > "${matched_log_file}"
  summarize_access_log_targets_by_minute < "${matched_log_file}" > "${output_file}"
}

function print_minute_traffic_access_table() {
  local traffic_file=$1
  local access_file=$2
  local rows

  rows=$(awk -F'|' '
    FILENAME == ARGV[1] {
      events[$1] = $2
      down[$1] = $3
      up[$1] = $4
      total[$1] = $5
      seen[$1] = 1
      next
    }
    FILENAME == ARGV[2] {
      targets[$1] = $2
      seen[$1] = 1
      next
    }
    END {
      for (minute in seen) {
        printf "%s|%s|%s|%s|%s|%s\n", minute, events[minute] + 0, down[minute] + 0, up[minute] + 0, total[minute] + 0, targets[minute]
      }
    }
  ' "${traffic_file}" "${access_file}" | sort)

  if [[ -z "${rows}" ]]; then
    echo "该时间段没有匹配的流量采集或访问日志"
    return 0
  fi

  printf "%-16s %-6s %-12s %-12s %-12s %s\n" "minute" "events" "down" "up" "total" "targets"
  printf "%-16s %-6s %-12s %-12s %-12s %s\n" "----------------" "------" "------------" "------------" "------------" "-------"
  while IFS="|" read -r minute events down up total targets; do
    [[ -z "${targets}" ]] && targets="-"
    printf "%-16s %-6s %-12s %-12s %-12s %s\n" "${minute}" "${events}" "$(format_bytes "${down}")" "$(format_bytes "${up}")" "$(format_bytes "${total}")" "${targets}"
  done <<< "${rows}"
}

function analyze_user_minute_report() {
  local db_path=$1
  local log_dir=$2
  local user_tag=$3
  local start_time=$4
  local end_time=$5
  local traffic_file
  local matched_log_file
  local access_file
  traffic_file=$(mktemp)
  matched_log_file=$(mktemp)
  access_file=$(mktemp)

  collect_traffic_minute_summary "${db_path}" "${user_tag}" "${start_time}" "${end_time}" "${traffic_file}" || {
    rm -f "${traffic_file}" "${matched_log_file}" "${access_file}"
    return 1
  }
  collect_access_minute_summary "${log_dir}" "${user_tag}" "${start_time}" "${end_time}" "${matched_log_file}" "${access_file}"

  echo
  echo -e "${Blue}===== 用户分钟级流量和访问目标聚合 =====${Font}"
  echo "用户 tag/id: ${user_tag}"
  echo "时间范围: ${start_time} ~ ${end_time}"
  echo "access.log 匹配行数: $(wc -l < "${matched_log_file}" | tr -d " ")"
  echo
  print_minute_traffic_access_table "${traffic_file}" "${access_file}"

  rm -f "${traffic_file}" "${matched_log_file}" "${access_file}"
}

function analyze_user_traffic_detail() {
  read -rp "请输入分析日期(格式 20260617): " analysis_date
  read -rp "请输入开始时间(格式 小时:分钟，例如 8:10): " start_clock
  read -rp "请输入结束时间(格式 小时:分钟，例如 9:00): " end_clock
  local normalized_range
  if ! normalized_range=$(analysis_datetime_range_from_parts "${analysis_date}" "${start_clock}" "${end_clock}"); then
    print_error "${normalized_range}"
    return 1
  fi
  IFS="|" read -r start_time end_time <<<"${normalized_range}"

  read -rp "请输入用户 tag/id（例如 8）: " user_tag
  if [[ -z "${user_tag}" ]]; then
    print_error "用户 tag/id 不能为空"
    return 1
  fi

  read -rp "SQLite 数据库路径(默认 ${stat_db_path}): " db_path
  [[ -z "${db_path}" ]] && db_path="${stat_db_path}"

  analyze_user_minute_report "${db_path}" "${xray_log_dir}" "${user_tag}" "${start_time}" "${end_time}"
}

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
        "clients": []
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

  print_ok "Xray 配置文件已创建"
  echo -e "${Yellow}Xray 初始 clients 为空，用户配置将由 xray-admin 数据库绑定生成。${Font}"
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
Environment="TZ=Asia/Shanghai"
ExecStart=/usr/local/bin/stat -port __STAT_PORT__ -level __STAT_LOG_LEVEL__ -traffic-db __STAT_DB_PATH__ -collect-interval __COLLECT_INTERVAL__ -traffic-retention-months __TRAFFIC_RETENTION_MONTHS__ -stat-api-traffic-tag __STAT_API_TRAFFIC_TAG__ -log-clean-dir __LOG_CLEAN_DIR__ -log-retention-months __LOG_RETENTION_MONTHS__ -xray-log-dir __XRAY_LOG_DIR__ -xray-log-retention-months __XRAY_LOG_RETENTION_MONTHS__
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
traffic_collect_interval: 1h
sysinfo_collect_interval: 5m
traffic-meter:
  enable: true
  stat-url: http://127.0.0.1:56611
  tag: "1"
  flush-interval: 10s
silicon-flow:
  api-key: ""
  base-url: https://api.siliconflow.cn
  model: deepseek-ai/DeepSeek-V3.2
  timeout: 30s
ADMINCONFIGEOF
  print_ok "xray_admin config.yaml 文件已创建"
}

function is_yes() {
  [[ "$1" == "y" || "$1" == "Y" || "$1" == "yes" || "$1" == "YES" ]]
}

function ask_yes_no() {
  local prompt=$1
  local default_answer=${2:-N}
  local answer
  local suffix

  if [[ "$default_answer" == "Y" || "$default_answer" == "y" ]]; then
    suffix="[Y/n]"
  else
    suffix="[y/N]"
  fi

  read -rp "${prompt} ${suffix}: " answer
  if [[ -z "$answer" ]]; then
    answer="$default_answer"
  fi
  is_yes "$answer"
}

function valid_ipv4() {
  local ip=$1
  local a b c d
  [[ "$ip" =~ ^[0-9]+(\.[0-9]+){3}$ ]] || return 1
  IFS=. read -r a b c d <<< "$ip"
  for part in "$a" "$b" "$c" "$d"; do
    [[ "$part" =~ ^[0-9]+$ ]] || return 1
    ((part >= 0 && part <= 255)) || return 1
  done
}

function normalize_port_list() {
  local input=$1
  local item
  local ports=()
  local seen=","

  input="${input//，/,}"
  input="${input//;/,}"
  input="${input// /,}"
  input="${input//$'\t'/,}"

  IFS=',' read -ra raw_ports <<< "$input"
  for item in "${raw_ports[@]}"; do
    item="$(echo "$item" | tr -d '[:space:]')"
    [[ -z "$item" ]] && continue
    [[ "$item" =~ ^[0-9]+$ ]] || return 1
    ((item >= 1 && item <= 65535)) || return 1
    if [[ "$seen" != *",$item,"* ]]; then
      ports+=("$item")
      seen="${seen}${item},"
    fi
  done

  [[ ${#ports[@]} -gt 0 ]] || return 1
  (IFS=,; echo "${ports[*]}")
}

function merge_port_lists() {
  local merged=""
  local item
  for item in "$@"; do
    [[ -z "$item" ]] && continue
    if [[ -z "$merged" ]]; then
      merged="$item"
    else
      merged="${merged},${item}"
    fi
  done
  normalize_port_list "$merged"
}

function port_list_contains() {
  local ports=",$1,"
  local port=$2
  [[ "$ports" == *",$port,"* ]]
}

function require_port_list() {
  local label=$1
  local detected=$2
  local default_value=$3
  local input
  local normalized

  if [[ -n "$detected" ]]; then
    if ask_yes_no "检测到${label}: ${detected}，是否使用" "Y"; then
      echo "$detected"
      return 0
    fi
  fi

  while true; do
    read -rp "请输入${label}(默认${default_value}，多个端口用逗号分隔): " input
    [[ -z "$input" ]] && input="$default_value"
    if normalized="$(normalize_port_list "$input")"; then
      echo "$normalized"
      return 0
    fi
    print_error "${label}格式不正确，请输入 1-65535 的端口，多个端口用逗号分隔"
  done
}

function require_ipv4() {
  local label=$1
  local detected=$2
  local input

  if [[ -n "$detected" ]] && valid_ipv4 "$detected"; then
    if ask_yes_no "检测到${label}: ${detected}，是否使用" "Y"; then
      echo "$detected"
      return 0
    fi
  fi

  while true; do
    read -rp "请输入${label}: " input
    if valid_ipv4 "$input"; then
      echo "$input"
      return 0
    fi
    print_error "${label}必须是 IPv4 地址"
  done
}

function detect_ssh_port() {
  local file
  local port
  for file in /etc/ssh/sshd_config /etc/ssh/sshd_config.d/*.conf; do
    [[ -f "$file" ]] || continue
    port="$(awk 'tolower($1) == "port" && $2 ~ /^[0-9]+$/ {print $2}' "$file" | tail -n 1)"
    [[ -n "$port" ]] && echo "$port" && return 0
  done
  if command -v ss >/dev/null 2>&1; then
    ss -tlnp 2>/dev/null | awk '/sshd/ {n=split($4,a,":"); print a[n]; exit}'
  fi
}

function detect_xray_ports() {
  local config_file=${1:-${xray_conf_dir}/config.json}
  [[ -f "$config_file" ]] || return 0
  if command -v jq >/dev/null 2>&1; then
    jq -r '
      .inbounds[]?
      | select((.port // empty) != "")
      | select((.tag // "") != "api")
      | select((.listen // "0.0.0.0") != "127.0.0.1")
      | select((.listen // "0.0.0.0") != "localhost")
      | select((.listen // "0.0.0.0") != "::1")
      | .port
    ' "$config_file" 2>/dev/null | grep -E '^[0-9]+$' | awk '!seen[$0]++' | paste -sd, -
    return 0
  fi

  awk '
    function trim_value(v) {
      gsub(/^[[:space:]"]+|[[:space:]",]+$/, "", v)
      return v
    }
    function flush_port() {
      if (port != "" && listen != "127.0.0.1" && listen != "localhost" && listen != "::1" && tag != "api") {
        if (!seen[port]++) {
          if (out != "") out = out ","
          out = out port
        }
      }
      port = ""
      listen = "0.0.0.0"
      tag = ""
    }
    /"inbounds"[[:space:]]*:/ {in_inbounds = 1}
    in_inbounds && /"outbounds"[[:space:]]*:/ {in_inbounds = 0}
    in_inbounds {
      line = $0
      if (line ~ /"port"[[:space:]]*:/) {
        value = line
        sub(/^.*"port"[[:space:]]*:[[:space:]]*"?/, "", value)
        sub(/[^0-9].*$/, "", value)
        value = trim_value(value)
        if (value ~ /^[0-9]+$/) port = value
      }
      if (line ~ /"listen"[[:space:]]*:/) {
        value = line
        sub(/^.*"listen"[[:space:]]*:[[:space:]]*"/, "", value)
        sub(/".*$/, "", value)
        listen = trim_value(value)
      }
      if (line ~ /"tag"[[:space:]]*:/) {
        value = line
        sub(/^.*"tag"[[:space:]]*:[[:space:]]*"/, "", value)
        sub(/".*$/, "", value)
        tag = trim_value(value)
      }
      if (line ~ /}/ && port != "") {
        flush_port()
      }
    }
    END { print out }
  ' "$config_file"
  return 0
}

function detect_stat_port() {
  local service_file=${1:-${stat_service_dir}}
  [[ -f "$service_file" ]] || return 0
  sed -nE 's/.*(^|[[:space:]])-port[= ]+([0-9]+).*/\2/p' "$service_file" | head -n 1
}

function detect_stat_remote_ip() {
  local service_file=${1:-${stat_service_dir}}
  [[ -f "$service_file" ]] || return 0
  grep -Eo 'REMOTE_IP=[^"[:space:]]+' "$service_file" | head -n 1 | cut -d= -f2
}

function detect_admin_port() {
  local config_file=${1:-${xray_admin_conf_dir}/config.yaml}
  [[ -f "$config_file" ]] || return 0
  awk '
    /^[[:space:]]*system:/ {in_system=1; next}
    in_system && /^[^[:space:]]/ {in_system=0}
    in_system && /^[[:space:]]*addr:/ {
      gsub(/"/, "", $2)
      if ($2 ~ /^[0-9]+$/) {
        print $2
        exit
      }
    }
  ' "$config_file"
}

function append_tcp_accept_rules() {
  local ports=$1
  local source_ip=${2:-}
  local port
  [[ -n "$ports" ]] || return 0
  IFS=',' read -ra port_array <<< "$ports"
  for port in "${port_array[@]}"; do
    [[ -z "$port" ]] && continue
    if [[ -n "$source_ip" ]]; then
      echo "-A INPUT -p tcp -s ${source_ip} -m tcp --dport ${port} -j ACCEPT"
    else
      echo "-A INPUT -p tcp -m tcp --dport ${port} -j ACCEPT"
    fi
  done
}

function generate_iptables_config() {
  local output_file=$1
  local ssh_ports=$2
  local xray_ports=$3
  local stat_port=$4
  local stat_source_ip=$5
  local admin_ports=${6:-}
  local extra_ports=${7:-}

  mkdir -p "$(dirname "$output_file")"
  {
    echo "*filter"
    echo ":INPUT DROP [0:0]"
    echo ":FORWARD DROP [0:0]"
    echo ":OUTPUT ACCEPT [0:0]"
    echo ""
    echo "-A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT"
    echo "-A INPUT -i lo -j ACCEPT"
    echo "-A INPUT -p icmp -j ACCEPT"
    append_tcp_accept_rules "$ssh_ports"
    append_tcp_accept_rules "$admin_ports"
    if [[ -n "$stat_port" && -n "$stat_source_ip" ]]; then
      append_tcp_accept_rules "$stat_port" "$stat_source_ip"
    fi
    append_tcp_accept_rules "$xray_ports"
    append_tcp_accept_rules "$extra_ports"
    echo ""
    echo "COMMIT"
  } > "$output_file"
}

function disable_ipv6() {
  print_ok "禁用 IPv6"
  cat > /etc/sysctl.d/99-xray-disable-ipv6.conf <<'SYSCTLEOF'
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6 = 1
SYSCTLEOF
  sysctl -w net.ipv6.conf.all.disable_ipv6=1 >/dev/null 2>&1 || true
  sysctl -w net.ipv6.conf.default.disable_ipv6=1 >/dev/null 2>&1 || true
  sysctl -w net.ipv6.conf.lo.disable_ipv6=1 >/dev/null 2>&1 || true
  sysctl -p /etc/sysctl.d/99-xray-disable-ipv6.conf >/dev/null 2>&1 || true
}

function save_iptables_rules() {
  if [[ -f /etc/os-release ]]; then
    source /etc/os-release
  fi

  if [[ "${ID:-}" == "centos" || "${ID:-}" == "ol" ]]; then
    service iptables save
    return $?
  fi

  if command -v netfilter-persistent >/dev/null 2>&1; then
    netfilter-persistent save
    return $?
  fi

  if [[ -d /etc/iptables ]]; then
    iptables-save > /etc/iptables/rules.v4
    return $?
  fi

  print_error "未找到持久化 iptables 规则的工具，请安装 iptables-persistent 或手动保存规则"
  return 1
}

function apply_iptables_with_rollback() {
  local config_file=$1
  local backup_file="${xray_conf_dir}/iptables.backup.$(date +%Y%m%d%H%M%S)"
  local marker="/tmp/xray-iptables-rollback.$$"
  local rollback_pid

  command -v iptables-save >/dev/null 2>&1 || {
    print_error "未找到 iptables-save"
    return 1
  }
  command -v iptables-restore >/dev/null 2>&1 || {
    print_error "未找到 iptables-restore"
    return 1
  }

  iptables-save > "$backup_file" || {
    print_error "备份当前 iptables 规则失败"
    return 1
  }
  print_ok "当前 iptables 已备份到 ${backup_file}"

  iptables-restore < "$config_file" || {
    print_error "应用 iptables 规则失败，正在恢复旧规则"
    iptables-restore < "$backup_file" 2>/dev/null || true
    return 1
  }

  touch "$marker"
  (
    sleep 120
    if [[ -f "$marker" ]]; then
      iptables-restore < "$backup_file" 2>/dev/null || true
      rm -f "$marker"
    fi
  ) &
  rollback_pid=$!

  echo -e "${Yellow}防火墙规则已临时应用。请立刻新开一个 SSH 窗口确认可以登录。${Font}"
  echo -e "${Yellow}如果 120 秒内没有确认，脚本会自动恢复旧规则。${Font}"
  if ask_yes_no "已确认 SSH 和必要服务可访问，是否保存为持久规则" "N"; then
    rm -f "$marker"
    kill "$rollback_pid" >/dev/null 2>&1 || true
    save_iptables_rules || return 1
    print_ok "iptables 规则已保存"
    return 0
  fi

  rm -f "$marker"
  kill "$rollback_pid" >/dev/null 2>&1 || true
  iptables-restore < "$backup_file" 2>/dev/null || true
  print_error "未确认保存，已恢复旧规则"
  return 1
}

function create_iptables_config() {
  generate_iptables_config "$iptables_conf_dir" "$FIREWALL_SSH_PORTS" "$FIREWALL_XRAY_PORTS" "$FIREWALL_STAT_PORT" "$FIREWALL_ADMIN_IP" "$FIREWALL_ADMIN_PORTS" "$FIREWALL_EXTRA_PORTS"
  print_ok "iptables 配置文件已创建: ${iptables_conf_dir}"
}

# ==================== 安装函数 ====================

function xray_install() {
  print_ok "安装 Xray"
  curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh | bash -s -- install --version ${xray_version}
  judge "Xray 安装"
}

function configure_xray() {
  create_xray_config
  configure_xray_logrotate
  create_xray_timezone_config
  systemctl daemon-reload
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
  mkdir -p ${stat_data_dir}
  chmod 755 ${stat_data_dir}
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

  read -rp "请输入Stat接口自身流量归属Tag(默认1)：" statApiTrafficTag
  [ -z "$statApiTrafficTag" ] && statApiTrafficTag="${stat_api_traffic_tag}"
  sed -i "s|__STAT_API_TRAFFIC_TAG__|${statApiTrafficTag}|" ${stat_service_dir}

  sed -i "s|__STAT_LOG_LEVEL__|${stat_log_level}|" ${stat_service_dir}
  sed -i "s|__STAT_DB_PATH__|${stat_db_path}|" ${stat_service_dir}
  sed -i "s|__COLLECT_INTERVAL__|${stat_collect_interval}|" ${stat_service_dir}
  sed -i "s|__TRAFFIC_RETENTION_MONTHS__|${stat_traffic_retention_months}|" ${stat_service_dir}
  sed -i "s|__LOG_CLEAN_DIR__|${stat_log_clean_dir}|" ${stat_service_dir}
  sed -i "s|__LOG_RETENTION_MONTHS__|${stat_log_retention_months}|" ${stat_service_dir}
  sed -i "s|__XRAY_LOG_DIR__|${stat_xray_log_dir}|" ${stat_service_dir}
  sed -i "s|__XRAY_LOG_RETENTION_MONTHS__|${stat_xray_log_retention_months}|" ${stat_service_dir}
  
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

  is_root

  local is_admin_node=0
  local is_proxy_node=0
  local detected

  FIREWALL_SSH_PORTS="$(require_port_list "SSH端口" "$(detect_ssh_port)" "22")"
  FIREWALL_ADMIN_PORTS=""
  FIREWALL_XRAY_PORTS=""
  FIREWALL_STAT_PORT=""
  FIREWALL_ADMIN_IP=""
  FIREWALL_EXTRA_PORTS=""

  if ask_yes_no "当前服务器是否存在 xray-admin 管理端" "N"; then
    is_admin_node=1
    detected="$(detect_admin_port)"
    FIREWALL_ADMIN_PORTS="$(require_port_list "xray-admin管理端端口" "$detected" "8888")"
  fi

  if ask_yes_no "当前服务器是否是代理节点" "Y"; then
    is_proxy_node=1
    detected="$(detect_xray_ports)"
    FIREWALL_XRAY_PORTS="$(require_port_list "Xray代理端口" "$detected" "80")"

    detected="$(detect_stat_port)"
    FIREWALL_STAT_PORT="$(require_port_list "Stat端口" "$detected" "56611")"

    detected="$(detect_stat_remote_ip)"
    FIREWALL_ADMIN_IP="$(require_ipv4 "允许访问Stat的管理端IP" "$detected")"
  fi

  if [[ "$is_admin_node" -eq 0 && "$is_proxy_node" -eq 0 ]]; then
    print_error "当前服务器既不是管理端也不是代理节点，不应配置此防火墙模板"
    return 1
  fi

  if ! port_list_contains "$(merge_port_lists "$FIREWALL_SSH_PORTS" "$FIREWALL_ADMIN_PORTS" "$FIREWALL_XRAY_PORTS" "$FIREWALL_STAT_PORT")" "80"; then
    if ask_yes_no "是否额外放行 80 端口(用于 HTTP/ACME standalone 等)" "N"; then
      FIREWALL_EXTRA_PORTS="$(merge_port_lists "$FIREWALL_EXTRA_PORTS" "80")"
    fi
  fi

  echo ""
  echo -e "${Blue}即将生成的最小放行规则:${Font}"
  echo "  SSH端口: ${FIREWALL_SSH_PORTS}"
  [[ -n "$FIREWALL_ADMIN_PORTS" ]] && echo "  xray-admin端口: ${FIREWALL_ADMIN_PORTS}"
  [[ -n "$FIREWALL_XRAY_PORTS" ]] && echo "  Xray代理端口: ${FIREWALL_XRAY_PORTS}"
  [[ -n "$FIREWALL_STAT_PORT" ]] && echo "  Stat端口: ${FIREWALL_STAT_PORT}，仅允许 ${FIREWALL_ADMIN_IP} 访问"
  [[ -n "$FIREWALL_EXTRA_PORTS" ]] && echo "  额外放行端口: ${FIREWALL_EXTRA_PORTS}"
  echo "  MySQL/Redis/数据库端口: 不放行"
  echo "  IPv6: 禁用"
  echo ""

  if ! ask_yes_no "确认生成并应用以上防火墙规则" "N"; then
    print_ok "已取消 iptables 配置"
    return 0
  fi

  confirm_overwrite "${iptables_conf_dir}" "iptables 配置文件" || return 0
  create_iptables_config
  disable_ipv6
  apply_iptables_with_rollback "$iptables_conf_dir"
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
  echo -e "${Green}11.${Font} 配置 Xray 日志轮转"
  echo -e "${Green}12.${Font} 分析用户流量明细"
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
  11)
    configure_xray_logrotate
    configure_xray_timezone
    if [[ "${XRAY_INSTALL_TEST_MODE:-0}" != "1" ]]; then
      systemctl restart xray 2>/dev/null && print_ok "xray 重启完成"
    fi
    ;;
  12)
    analyze_user_traffic_detail
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

if [[ "${XRAY_INSTALL_TEST_MODE:-0}" != "1" ]]; then
  menu "$@"
fi
