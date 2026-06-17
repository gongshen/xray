#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
export XRAY_INSTALL_TEST_MODE=1

source "${ROOT_DIR}/deploy/install.sh"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "${tmp_dir}"' EXIT

assert_eq() {
  local expected="$1"
  local actual="$2"
  local message="$3"
  if [[ "${expected}" != "${actual}" ]]; then
    echo "FAIL: ${message}: expected '${expected}', got '${actual}'" >&2
    exit 1
  fi
}

assert_eq "10s" "${stat_collect_interval}" "default stat collect interval"

cat > "${tmp_dir}/xray.json" <<'JSON'
{
  "inbounds": [
    {"listen": "0.0.0.0", "port": 80, "protocol": "vmess"},
    {"listen": "::", "port": 443, "protocol": "vless"},
    {"listen": "127.0.0.1", "port": 11111, "protocol": "dokodemo-door", "tag": "api"}
  ]
}
JSON
assert_eq "80,443" "$(detect_xray_ports "${tmp_dir}/xray.json")" "detect public xray ports"

cat > "${tmp_dir}/stat.service" <<'EOFSTAT'
[Service]
Environment="REMOTE_IP=203.0.113.10"
ExecStart=/usr/local/bin/stat -port 56611 -traffic-db /var/lib/xray-stat/stat.db -collect-interval 10s
EOFSTAT
assert_eq "56611" "$(detect_stat_port "${tmp_dir}/stat.service")" "detect stat port"
assert_eq "203.0.113.10" "$(detect_stat_remote_ip "${tmp_dir}/stat.service")" "detect stat remote ip"

cat > "${tmp_dir}/config.yaml" <<'YAML'
system:
  env: public
  addr: 8888
YAML
assert_eq "8888" "$(detect_admin_port "${tmp_dir}/config.yaml")" "detect admin port"

cat > "${tmp_dir}/config-https.yaml" <<'YAML'
system:
  env: public
  addr: 443
YAML
assert_eq "443" "$(detect_admin_port "${tmp_dir}/config-https.yaml")" "detect https admin port"

rules_file="${tmp_dir}/iptables.rules"
generate_iptables_config "${rules_file}" "22" "80,443" "56611" "203.0.113.10" "8888"

grep -q -- "-A INPUT -p tcp -m tcp --dport 22 -j ACCEPT" "${rules_file}"
grep -q -- "-A INPUT -p tcp -m tcp --dport 80 -j ACCEPT" "${rules_file}"
grep -q -- "-A INPUT -p tcp -m tcp --dport 443 -j ACCEPT" "${rules_file}"
grep -q -- "-A INPUT -p tcp -s 203.0.113.10 -m tcp --dport 56611 -j ACCEPT" "${rules_file}"
grep -q -- "-A INPUT -p tcp -m tcp --dport 8888 -j ACCEPT" "${rules_file}"
if grep -q -- "--dport 3306" "${rules_file}"; then
  echo "FAIL: mysql port must not be opened by firewall config" >&2
  exit 1
fi

https_admin_rules_file="${tmp_dir}/iptables-https-admin.rules"
generate_iptables_config "${https_admin_rules_file}" "22" "80" "56611" "203.0.113.10" "443"
grep -q -- "-A INPUT -p tcp -m tcp --dport 443 -j ACCEPT" "${https_admin_rules_file}"

logrotate_file="${tmp_dir}/xray-logrotate"
create_xray_logrotate_config "${logrotate_file}"

grep -q -- "/var/log/xray/access.log /var/log/xray/error.log" "${logrotate_file}"
grep -q -- "daily" "${logrotate_file}"
grep -q -- "rotate 365" "${logrotate_file}"
grep -q -- "maxage 365" "${logrotate_file}"
grep -q -- "missingok" "${logrotate_file}"
grep -q -- "notifempty" "${logrotate_file}"
grep -q -- "compress" "${logrotate_file}"
grep -q -- "delaycompress" "${logrotate_file}"
grep -q -- "copytruncate" "${logrotate_file}"
grep -q -- "dateext" "${logrotate_file}"

menu_logrotate_file="${tmp_dir}/menu-xray-logrotate"
xray_logrotate_conf_dir="${menu_logrotate_file}"
printf '11\n' | menu >/dev/null
grep -q -- "/var/log/xray/access.log /var/log/xray/error.log" "${menu_logrotate_file}"

assert_eq "'abc''def'" "$(sqlite_quote "abc'def")" "sqlite quote escapes single quote"

traffic_where="$(build_traffic_event_where "8" "2026-06-17 00:00:00" "2026-06-18 00:00:00")"
grep -q -- "tag = '8'" <<<"${traffic_where}"
grep -q -- "collected_at >= strftime('%s', '2026-06-17 00:00:00')" <<<"${traffic_where}"
grep -q -- "collected_at <= strftime('%s', '2026-06-18 00:00:00')" <<<"${traffic_where}"

traffic_detail_sql="$(build_traffic_detail_sql "8" "2026-06-17 00:00:00" "2026-06-17 23:59:59" "10")"
grep -q -- "printf('%.2fM', (down + up) / 1048576.0) AS total" <<<"${traffic_detail_sql}"

traffic_top_windows_sql="$(build_traffic_top_windows_sql "8" "2026-06-17 00:00:00" "2026-06-17 23:59:59" "60" "10")"
grep -q -- "AS window_start" <<<"${traffic_top_windows_sql}"
grep -q -- "printf('%.2fM', SUM(down + up) / 1048576.0) AS total" <<<"${traffic_top_windows_sql}"
grep -q -- "ORDER BY SUM(down + up) DESC" <<<"${traffic_top_windows_sql}"

assert_eq "2026-06-17 00:00:00|2026-06-17 23:59:59" "$(analysis_date_to_time_range "2026-06-17")" "date expands to one day"
if analysis_date_to_time_range "2026-06-17 12:00:00" >/dev/null; then
  echo "FAIL: analysis date must reject datetime input" >&2
  exit 1
fi
if analysis_date_to_time_range "2026-02-30" >/dev/null; then
  echo "FAIL: analysis date must reject invalid date" >&2
  exit 1
fi

assert_eq "email: 8$" "$(build_access_log_email_pattern "8")" "access log numeric tag regex"
assert_eq "email: a\\.b\\+c$" "$(build_access_log_email_pattern "a.b+c")" "access log regex escapes tag"

target_summary="$(cat <<'ACCESSLOG' | summarize_access_log_targets 5
2026/06/17 10:00:00 1.1.1.1:10001 accepted tcp:googlevideo.com:443 email: 8
2026/06/17 10:00:01 1.1.1.1:10002 accepted tcp:googlevideo.com:443 email: 8
2026/06/17 10:00:02 1.1.1.1:10003 accepted udp:mtalk.google.com:5228 email: 8
2026/06/17 10:00:03 tcp:github.com:443 accepted email: 8
ACCESSLOG
)"
grep -Eq "2[[:space:]]+googlevideo\\.com" <<<"${target_summary}"
grep -Eq "1[[:space:]]+mtalk\\.google\\.com" <<<"${target_summary}"
grep -Eq "1[[:space:]]+github\\.com" <<<"${target_summary}"

touch "${tmp_dir}/access.log"
touch "${tmp_dir}/access.log-20260617.gz"
touch "${tmp_dir}/access.log-2026-06-17.gz"
touch "${tmp_dir}/access.log.1"
touch "${tmp_dir}/error.log"
touch "${tmp_dir}/access.log.backup"
access_files="$(find_xray_access_log_files "${tmp_dir}" | sed "s#${tmp_dir}/##" | sort | tr '\n' ',')"
assert_eq "access.log,access.log-2026-06-17.gz,access.log-20260617.gz," "${access_files}" "find date-based access logs only"

menu_log_dir="${tmp_dir}/menu-logs"
mkdir -p "${menu_log_dir}"
cat > "${menu_log_dir}/access.log" <<'ACCESSLOG'
2026/06/17 10:00:00 tcp:example.com:443 accepted email: 8
2026/06/17 10:01:00 tcp:other.example:443 accepted email: 18
ACCESSLOG
xray_log_dir="${menu_log_dir}"
menu_output="$(printf '12\n8\n%s\n2026-06-17\n1\n5\n' "${tmp_dir}/missing.db" | menu || true)"
grep -q -- "tcp:example.com:443 accepted email: 8" <<<"${menu_output}"
if grep -q -- "tcp:other.example:443 accepted email: 18" <<<"${menu_output}"; then
  echo "FAIL: user traffic menu matched another user's access log" >&2
  exit 1
fi

echo "install_firewall_test passed"
