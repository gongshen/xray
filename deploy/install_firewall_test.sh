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
ExecStart=/usr/local/bin/stat -port 56611 -traffic-db /var/lib/xray-stat/stat.db -collect-interval 5s
EOFSTAT
assert_eq "56611" "$(detect_stat_port "${tmp_dir}/stat.service")" "detect stat port"
assert_eq "203.0.113.10" "$(detect_stat_remote_ip "${tmp_dir}/stat.service")" "detect stat remote ip"

cat > "${tmp_dir}/config.yaml" <<'YAML'
system:
  env: public
  addr: 8888
YAML
assert_eq "8888" "$(detect_admin_port "${tmp_dir}/config.yaml")" "detect admin port"

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

echo "install_firewall_test passed"
