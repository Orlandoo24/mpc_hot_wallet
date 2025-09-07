#!/usr/bin/env bash
# filename: airtasker_recon_v2.sh
set -euo pipefail

DOMAIN="airtasker.com"
IPS=("13.54.63.107" "54.252.94.120" "54.252.94.121" "54.252.94.122")

TS="$(date +%Y%m%d_%H%M%S)"
OUTDIR="recon_${DOMAIN}_${TS}"
mkdir -p "$OUTDIR"

bold(){ printf "\033[1m%s\033[0m\n" "$*"; }
have(){ command -v "$1" >/dev/null 2>&1; }

bold "[*] Output dir: $OUTDIR"
for b in nmap curl whois dig openssl; do have "$b" || { echo "!! missing $b"; exit 1; }; done

# 选择更高效的扫描方式：root 用 -sS，非 root 用 -sT
if sudo -n true 2>/dev/null; then SCAN="-sS"; else SCAN="-sT"; fi

# 0) DNS / PTR
bold "[*] DNS quick look…"
{
  echo "# A/AAAA of ${DOMAIN}"
  dig +short A   "${DOMAIN}"
  dig +short AAAA "${DOMAIN}"
  echo
  for ip in "${IPS[@]}"; do
    echo "# PTR of ${ip}"
    dig +short -x "${ip}"
  done
} > "${OUTDIR}/dns_overview.txt" || true

# 1) Top-200 端口 + 指纹（过滤会吵的 NSE）
bold "[*] Nmap top-ports(200) + fingerprint…"
SAFE_NSE='default and not broadcast and not external and not targets-asn and not hostmap-robtex and not http-robtex-shared-ns'
for ip in "${IPS[@]}"; do
  nmap -Pn ${SCAN} -T4 --top-ports 200 \
       --max-retries 2 --host-timeout 5m \
       --script "${SAFE_NSE}" \
       -oA "${OUTDIR}/${ip}_top200" "${ip}" || true
done

# 2) Web(80,443) 信息脚本（不含 http2）
bold "[*] Nmap web info (80,443)…"
HTTP_SCRIPTS="http-title,http-headers,http-security-headers,http-methods,http-robots.txt,http-cookie-flags"
TLS_SCRIPTS="ssl-cert,ssl-enum-ciphers,tls-alpn,tls-nextprotoneg"
for ip in "${IPS[@]}"; do
  nmap -Pn ${SCAN} -T4 -p 80,443 -sV --version-all \
       --max-retries 2 --host-timeout 5m \
       --script "${SAFE_NSE},${HTTP_SCRIPTS},${TLS_SCRIPTS}" \
       --script-args http.useragent="Airtasker-Baseline-Recon/1.1" \
       -oA "${OUTDIR}/${ip}_web" "${ip}" || true
done

# 3) 全端口“开放速览”（温和参数）
bold "[*] Nmap full-port open summary…"
for ip in "${IPS[@]}"; do
  nmap -Pn --open -p- ${SCAN} -T3 \
       --max-retries 2 --host-timeout 5m \
       -oA "${OUTDIR}/${ip}_full" "${ip}" || true
done

# 4) VHost/SNI 头部与证书（只读）
bold "[*] VHost-aware HTTPS headers & certs…"
curl -sS -I --max-time 10 "https://${DOMAIN}/" > "${OUTDIR}/${DOMAIN}_https.headers" || true
curl -sS -I --max-time 10 "http://${DOMAIN}/"  > "${OUTDIR}/${DOMAIN}_http.headers"  || true
for ip in "${IPS[@]}"; do
  curl -sS -I --max-time 10 --resolve "${DOMAIN}:443:${ip}" "https://${DOMAIN}/" \
       > "${OUTDIR}/${DOMAIN}_on_${ip}.headers" || true
  { echo | openssl s_client -servername "${DOMAIN}" -connect "${ip}:443" 2>/dev/null \
      | openssl x509 -noout -issuer -subject -dates; } \
      > "${OUTDIR}/${DOMAIN}_on_${ip}.cert.txt" || true
done

# 5) WHOIS
bold "[*] WHOIS…"
for ip in "${IPS[@]}"; do whois "${ip}" > "${OUTDIR}/${ip}.whois.txt" || true; done

# 6) 汇总
bold "[*] Compose SUMMARY…"
SUMMARY="${OUTDIR}/SUMMARY.txt"
{
  echo "=== Recon Summary (${TS}) for ${DOMAIN} ==="
  echo
  echo "[DNS]"
  cat "${OUTDIR}/dns_overview.txt" 2>/dev/null || true
  echo
  for ip in "${IPS[@]}"; do
    echo "----- ${ip} -----"
    [ -f "${OUTDIR}/${ip}_top200.gnmap" ] && { echo "[Open (top200)]"; grep -E "Ports: " "${OUTDIR}/${ip}_top200.gnmap" | sed 's/.*Ports: //'; }
    [ -f "${OUTDIR}/${ip}_full.gnmap"   ] && { echo "[Open (full)]";   grep -E "Ports: " "${OUTDIR}/${ip}_full.gnmap"   | sed 's/.*Ports: //'; }
    if [ -f "${OUTDIR}/${DOMAIN}_on_${ip}.headers" ]; then
      echo "[Security Headers (via IP+SNI)]"
      awk 'BEGIN{IGNORECASE=1}
           /^strict-transport-security:/{print}
           /^content-security-policy:/{print}
           /^x-frame-options:/{print}
           /^x-content-type-options:/{print}
           /^referrer-policy:/{print}
           /^permissions-policy:/{print}
           /^set-cookie:/{print}' "${OUTDIR}/${DOMAIN}_on_${ip}.headers"
    fi
    [ -f "${OUTDIR}/${DOMAIN}_on_${ip}.cert.txt" ] && { echo "[Cert (SNI=${DOMAIN})]"; cat "${OUTDIR}/${DOMAIN}_on_${ip}.cert.txt"; }
    echo
  done
  echo "=== END ==="
} > "$SUMMARY"

bold "[*] Done."
echo "Report dir : $OUTDIR"
echo "Quick read : $SUMMARY"
