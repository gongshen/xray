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
xray_timezone_conf_file="${tmp_dir}/menu-xray-timezone.conf"
printf '11\n' | menu >/dev/null
grep -q -- "/var/log/xray/access.log /var/log/xray/error.log" "${menu_logrotate_file}"
grep -q -- 'Environment="TZ=Asia/Shanghai"' "${xray_timezone_conf_file}"

assert_eq "'abc''def'" "$(sqlite_quote "abc'def")" "sqlite quote escapes single quote"

traffic_where="$(build_traffic_event_where "8" "2026-06-17 00:00:00" "2026-06-18 00:00:00")"
grep -q -- "tag = '8'" <<<"${traffic_where}"
grep -q -- "collected_at >= strftime('%s', '2026-06-17 00:00:00', '-8 hours')" <<<"${traffic_where}"
grep -q -- "collected_at <= strftime('%s', '2026-06-18 00:00:00', '-8 hours')" <<<"${traffic_where}"

traffic_minute_sql="$(build_traffic_minute_summary_sql "8" "2026-06-17 08:10:00" "2026-06-17 09:00:59")"
grep -q -- "strftime('%Y-%m-%d %H:%M', collected_at, 'unixepoch', '+8 hours') AS minute" <<<"${traffic_minute_sql}"
grep -q -- "GROUP BY minute" <<<"${traffic_minute_sql}"
grep -q -- "ORDER BY minute ASC" <<<"${traffic_minute_sql}"

stat_service_output="${tmp_dir}/stat-service-generated.service"
stat_service_dir="${stat_service_output}"
printf 'y\n' | create_stat_service >/dev/null
grep -q -- 'Environment="TZ=Asia/Shanghai"' "${stat_service_output}"

xray_timezone_output="${tmp_dir}/xray-timezone.conf"
create_xray_timezone_config "${xray_timezone_output}" >/dev/null
grep -q -- 'Environment="TZ=Asia/Shanghai"' "${xray_timezone_output}"

assert_eq "2026-06-17" "$(normalize_analysis_date "20260617")" "compact date normalizes"
assert_eq "2026-06-17" "$(normalize_analysis_date "2026-06-17")" "hyphen date normalizes"
assert_eq "08:10:00" "$(normalize_analysis_clock "8:10" "start")" "start clock normalizes"
assert_eq "09:00:59" "$(normalize_analysis_clock "9:00" "end")" "end clock includes full minute"
assert_eq "2026-06-17 08:10:00|2026-06-17 09:00:59" "$(analysis_datetime_range_from_parts "20260617" "8:10" "9:00")" "date and clock expand to range"
assert_eq "2026-06-17 08:10:00|2026-06-17 10:10:59" "$(analysis_datetime_range_from_parts "20260617" "8:10" "10:10")" "two hour range is allowed"
if normalize_analysis_date "2026-02-30" >/dev/null; then
  echo "FAIL: analysis date must reject invalid date" >&2
  exit 1
fi
if analysis_datetime_range_from_parts "20260617" "9:00" "8:10" >/dev/null; then
  echo "FAIL: analysis range must reject end before start" >&2
  exit 1
fi
if analysis_datetime_range_from_parts "20260617" "8:10" "10:11" >/dev/null; then
  echo "FAIL: analysis range must reject ranges longer than 2 hours" >&2
  exit 1
fi

assert_eq "email: 8$" "$(build_access_log_email_pattern "8")" "access log numeric tag regex"
assert_eq "email: a\\.b\\+c$" "$(build_access_log_email_pattern "a.b+c")" "access log regex escapes tag"

minute_target_summary="$(cat <<'ACCESSLOG' | summarize_access_log_targets_by_minute
2026/06/17 10:00:00 1.1.1.1:10001 accepted tcp:rr2---sn-3pm7dne6.googlevideo.com:443 email: 8
2026/06/17 10:00:10 1.1.1.1:10002 accepted tcp:r3---sn-3pm7dn7r.googlevideo.com:443 email: 8
2026/06/17 10:01:00 1.1.1.1:10003 accepted udp:android.clients.google.com:5228 email: 8
2026/06/17 10:01:10 1.1.1.1:10004 accepted udp:mtalk.google.com:5228 email: 8
2026/06/17 10:01:20 1.1.1.1:10005 accepted tcp:8.8.8.8:443 email: 8
2026/06/17 10:01:30 1.1.1.1:10006 accepted tcp:8.8.8.8:443 email: 8
ACCESSLOG
)"
grep -q -- "2026-06-17 10:00|googlevideo.com" <<<"${minute_target_summary}"
if grep -q -- "googlevideo.com(" <<<"${minute_target_summary}"; then
  echo "FAIL: access target summary must not include counts" >&2
  exit 1
fi
if grep -q -- "rr2---sn-3pm7dne6.googlevideo.com" <<<"${minute_target_summary}"; then
  echo "FAIL: googlevideo subdomain must be normalized" >&2
  exit 1
fi
grep -q -- "2026-06-17 10:01|google.com" <<<"${minute_target_summary}"
grep -q -- "8.8.8.8" <<<"${minute_target_summary}"
if grep -q -- "android.clients.google.com" <<<"${minute_target_summary}"; then
  echo "FAIL: google subdomain must be normalized" >&2
  exit 1
fi
if grep -q -- "mtalk.google.com" <<<"${minute_target_summary}"; then
  echo "FAIL: mtalk google subdomain must be normalized" >&2
  exit 1
fi

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
2026/06/17 10:00:15 tcp:example.com:443 accepted email: 8
2026/06/17 10:01:00 tcp:github.com:443 accepted email: 8
2026/06/17 10:01:00 tcp:other.example:443 accepted email: 18
ACCESSLOG
xray_log_dir="${menu_log_dir}"
menu_output="$(printf '12\n20260617\n10:00\n10:01\n8\n%s\n' "${tmp_dir}/missing.db" | menu || true)"
grep -q -- "2026-06-17 10:00" <<<"${menu_output}"
grep -q -- "example.com" <<<"${menu_output}"
grep -q -- "github.com" <<<"${menu_output}"
if grep -q -- "example.com(" <<<"${menu_output}"; then
  echo "FAIL: menu target aggregation must not include counts" >&2
  exit 1
fi
if grep -q -- "tcp:example.com:443 accepted email: 8" <<<"${menu_output}"; then
  echo "FAIL: user traffic menu must not print raw access logs" >&2
  exit 1
fi
if grep -q -- "tcp:other.example:443 accepted email: 18" <<<"${menu_output}"; then
  echo "FAIL: user traffic menu matched another user's access log" >&2
  exit 1
fi

echo "install_firewall_test passed"
