#!/usr/bin/env bash
# filename: airtasker_recon_v2.sh (Optimized Logging)
set -euo pipefail

DOMAIN="airtasker.com"
IPS=("13.54.63.107" "54.252.94.120" "54.252.94.121" "54.252.94.122")

TS="$(date +%Y%m%d_%H%M%S)"
OUTDIR="recon_${DOMAIN}_${TS}"
mkdir -p "$OUTDIR"

# Enhanced logging functions
bold(){ printf "\033[1m%s\033[0m\n" "$*"; }
green(){ printf "\033[32m%s\033[0m\n" "$*"; }
yellow(){ printf "\033[33m%s\033[0m\n" "$*"; }
red(){ printf "\033[31m%s\033[0m\n" "$*"; }
blue(){ printf "\033[34m%s\033[0m\n" "$*"; }
have(){ command -v "$1" >/dev/null 2>&1; }

# Progress tracking
log_step() {
    local step_name="$1"
    local current="$2"
    local total="$3"
    local timestamp=$(date '+%H:%M:%S')
    printf "\n\033[1;34m[%s] [%d/%d] %s\033[0m\n" "$timestamp" "$current" "$total" "$step_name"
}

log_substep() {
    local substep="$1"
    local status="${2:-INFO}"
    local timestamp=$(date '+%H:%M:%S')
    case "$status" in
        "SUCCESS") printf "\033[32m  ✓ [%s] %s\033[0m\n" "$timestamp" "$substep" ;;
        "ERROR")   printf "\033[31m  ✗ [%s] %s\033[0m\n" "$timestamp" "$substep" ;;
        "WARN")    printf "\033[33m  ⚠ [%s] %s\033[0m\n" "$timestamp" "$substep" ;;
        *)         printf "\033[37m  ▶ [%s] %s\033[0m\n" "$timestamp" "$substep" ;;
    esac
}

log_summary() {
    local action="$1"
    local target="$2"
    local result="$3"
    local duration="$4"
    printf "\033[36m    📊 %s on %s: %s (Duration: %s)\033[0m\n" "$action" "$target" "$result" "$duration"
}

start_timer() {
    echo $(date +%s)
}

end_timer() {
    local start_time="$1"
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    printf "%dm %ds" $((duration/60)) $((duration%60))
}

# Initialization
printf "\n\033[1;36m╔══════════════════════════════════════════════════════╗\033[0m\n"
printf "\033[1;36m║           🔍 AIRTASKER RECON TOOL v2.0               ║\033[0m\n"
printf "\033[1;36m╠══════════════════════════════════════════════════════╣\033[0m\n"
printf "\033[1;36m║ Target Domain: %-37s ║\033[0m\n" "$DOMAIN"
printf "\033[1;36m║ Target IPs: %-40s ║\033[0m\n" "${#IPS[@]} addresses"
printf "\033[1;36m║ Output Directory: %-32s ║\033[0m\n" "$OUTDIR"
printf "\033[1;36m║ Scan Start: %-37s ║\033[0m\n" "$(date '+%Y-%m-%d %H:%M:%S %Z')"
printf "\033[1;36m╚══════════════════════════════════════════════════════╝\033[0m\n"

# Tool dependency check
log_substep "Checking required tools..." "INFO"
missing_tools=()
for tool in nmap curl whois dig openssl; do 
    if have "$tool"; then
        log_substep "$tool: Available" "SUCCESS"
    else
        log_substep "$tool: MISSING" "ERROR"
        missing_tools+=("$tool")
    fi
done

if [ ${#missing_tools[@]} -gt 0 ]; then
    red "Missing required tools: ${missing_tools[*]}"
    exit 1
fi

# 选择更高效的扫描方式：root 用 -sS，非 root 用 -sT
if sudo -n true 2>/dev/null; then 
    SCAN="-sS"
    log_substep "Scan mode: TCP SYN (privileged)" "SUCCESS"
else 
    SCAN="-sT"
    log_substep "Scan mode: TCP Connect (unprivileged)" "WARN"
fi

# 0) DNS / PTR Analysis
TOTAL_STEPS=6
log_step "DNS Resolution & PTR Lookup" 1 $TOTAL_STEPS
step_start=$(start_timer)

log_substep "Resolving A/AAAA records for ${DOMAIN}" "INFO"
dns_start=$(start_timer)
{
  echo "# A/AAAA of ${DOMAIN}"
  a_records=$(dig +short A "${DOMAIN}" 2>/dev/null || echo "N/A")
  aaaa_records=$(dig +short AAAA "${DOMAIN}" 2>/dev/null || echo "N/A")
  echo "$a_records"
  echo "$aaaa_records"
  echo
  
  for ip in "${IPS[@]}"; do
    echo "# PTR of ${ip}"
    ptr_record=$(dig +short -x "${ip}" 2>/dev/null || echo "N/A")
    echo "$ptr_record"
  done
} > "${OUTDIR}/dns_overview.txt"

if [ $? -eq 0 ]; then
    log_substep "DNS resolution completed" "SUCCESS"
    if [ "$a_records" != "N/A" ] && [ -n "$a_records" ]; then
        log_substep "A records found: $(echo $a_records | wc -w) entries" "SUCCESS"
    else
        log_substep "No A records found" "WARN"
    fi
else
    log_substep "DNS resolution failed" "ERROR"
fi

dns_duration=$(end_timer $dns_start)
log_summary "DNS Lookup" "$DOMAIN + ${#IPS[@]} IPs" "Completed" "$dns_duration"

# 1) Top-200 Port Scan + Fingerprinting
log_step "Nmap Top-200 Port Scan + Service Detection" 2 $TOTAL_STEPS

SAFE_NSE='default and not broadcast and not external and not targets-asn and not hostmap-robtex and not http-robtex-shared-ns'
scan_results=()

for i in "${!IPS[@]}"; do
    ip="${IPS[$i]}"
    current_scan=$((i + 1))
    total_ips=${#IPS[@]}
    
    log_substep "Scanning $ip (${current_scan}/${total_ips}) - Top 200 ports" "INFO"
    scan_start=$(start_timer)
    
    # Run nmap with output suppression but capture results
    nmap_output="${OUTDIR}/${ip}_top200.nmap"
    if nmap -Pn ${SCAN} -T4 --top-ports 200 \
           --max-retries 2 --host-timeout 5m \
           --script "${SAFE_NSE}" \
           -oA "${OUTDIR}/${ip}_top200" "${ip}" > "$nmap_output" 2>&1; then
        
        # Parse results
        open_ports=$(grep -E "^[0-9]+/tcp.*open" "$nmap_output" | wc -l | xargs || echo "0")
        filtered_ports=$(grep -E "^[0-9]+/tcp.*filtered" "$nmap_output" | wc -l | xargs || echo "0")
        
        scan_duration=$(end_timer $scan_start)
        if [ "$open_ports" -gt 0 ]; then
            log_substep "$ip: Found $open_ports open ports, $filtered_ports filtered" "SUCCESS"
            scan_results+=("$ip:SUCCESS:$open_ports open, $filtered_ports filtered")
        elif [ "$filtered_ports" -gt 0 ]; then
            log_substep "$ip: $filtered_ports ports filtered (possible firewall)" "WARN"
            scan_results+=("$ip:FILTERED:$filtered_ports filtered")
        else
            log_substep "$ip: No open ports detected" "WARN"
            scan_results+=("$ip:CLOSED:All ports closed/filtered")
        fi
        log_summary "Port Scan" "$ip" "Completed" "$scan_duration"
    else
        scan_duration=$(end_timer $scan_start)
        log_substep "$ip: Scan failed or timed out" "ERROR"
        log_summary "Port Scan" "$ip" "Failed/Timeout" "$scan_duration"
        scan_results+=("$ip:ERROR:Scan failed")
    fi
done

# 2) Web Services Analysis (80,443)
log_step "Web Services & Security Analysis" 3 $TOTAL_STEPS

HTTP_SCRIPTS="http-title,http-headers,http-security-headers,http-methods,http-robots.txt,http-cookie-flags"
TLS_SCRIPTS="ssl-cert,ssl-enum-ciphers,tls-alpn,tls-nextprotoneg"
web_results=()

for i in "${!IPS[@]}"; do
    ip="${IPS[$i]}"
    current_web=$((i + 1))
    
    log_substep "Analyzing web services on $ip (${current_web}/${#IPS[@]})" "INFO"
    web_start=$(start_timer)
    
    web_output="${OUTDIR}/${ip}_web.nmap"
    if nmap -Pn ${SCAN} -T4 -p 80,443 -sV --version-all \
           --max-retries 2 --host-timeout 5m \
           --script "${SAFE_NSE},${HTTP_SCRIPTS},${TLS_SCRIPTS}" \
           --script-args http.useragent="Airtasker-Baseline-Recon/1.1" \
           -oA "${OUTDIR}/${ip}_web" "${ip}" > "$web_output" 2>&1; then
        
        # Parse web results
        http_open=$(grep -E "80/tcp.*open" "$web_output" | wc -l | xargs || echo "0")
        https_open=$(grep -E "443/tcp.*open" "$web_output" | wc -l | xargs || echo "0")
        http_filtered=$(grep -E "80/tcp.*filtered" "$web_output" | wc -l | xargs || echo "0")
        https_filtered=$(grep -E "443/tcp.*filtered" "$web_output" | wc -l | xargs || echo "0")
        
        web_duration=$(end_timer $web_start)
        
        status_msg=""
        if [ "$http_open" -gt 0 ] || [ "$https_open" -gt 0 ]; then
            [ "$http_open" -gt 0 ] && status_msg+="HTTP:open "
            [ "$https_open" -gt 0 ] && status_msg+="HTTPS:open "
            log_substep "$ip: $status_msg" "SUCCESS"
            web_results+=("$ip:OPEN:$status_msg")
        elif [ "$http_filtered" -gt 0 ] || [ "$https_filtered" -gt 0 ]; then
            [ "$http_filtered" -gt 0 ] && status_msg+="HTTP:filtered "
            [ "$https_filtered" -gt 0 ] && status_msg+="HTTPS:filtered "
            log_substep "$ip: $status_msg" "WARN"
            web_results+=("$ip:FILTERED:$status_msg")
        else
            log_substep "$ip: No web services detected" "WARN"
            web_results+=("$ip:NONE:No web services")
        fi
        
        log_summary "Web Analysis" "$ip" "Completed" "$web_duration"
    else
        web_duration=$(end_timer $web_start)
        log_substep "$ip: Web analysis failed" "ERROR"
        log_summary "Web Analysis" "$ip" "Failed" "$web_duration"
        web_results+=("$ip:ERROR:Analysis failed")
    fi
done

# 3) Full Port Scan (65535 ports) - CAUTION: This is slow!
log_step "Full Port Scan (All 65535 ports)" 4 $TOTAL_STEPS

yellow "⚠️  WARNING: Full port scan can take 5-10 minutes per IP"
read -p "Continue with full port scan? (y/N): " -n 1 -r -t 10
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    full_results=()
    
    for i in "${!IPS[@]}"; do
        ip="${IPS[$i]}"
        current_full=$((i + 1))
        
        log_substep "Full port scan on $ip (${current_full}/${#IPS[@]}) - This may take 5-10 minutes" "INFO"
        full_start=$(start_timer)
        
        full_output="${OUTDIR}/${ip}_full.nmap"
        if timeout 600 nmap -Pn --open -p- ${SCAN} -T3 \
               --max-retries 2 --host-timeout 5m \
               -oA "${OUTDIR}/${ip}_full" "${ip}" > "$full_output" 2>&1; then
            
            # Parse full scan results
            total_open=$(grep -E "^[0-9]+/tcp.*open" "$full_output" | wc -l | xargs || echo "0")
            
            full_duration=$(end_timer $full_start)
            if [ "$total_open" -gt 0 ]; then
                log_substep "$ip: Found $total_open total open ports" "SUCCESS"
                full_results+=("$ip:SUCCESS:$total_open open ports")
            else
                log_substep "$ip: No additional open ports found" "WARN"
                full_results+=("$ip:NONE:No open ports")
            fi
            log_summary "Full Scan" "$ip" "Completed" "$full_duration"
        else
            full_duration=$(end_timer $full_start)
            log_substep "$ip: Full scan timed out or failed" "ERROR"
            log_summary "Full Scan" "$ip" "Timeout/Failed" "$full_duration"
            full_results+=("$ip:TIMEOUT:Scan incomplete")
        fi
    done
else
    log_substep "Full port scan skipped by user" "WARN"
    full_results=("SKIPPED:User choice")
fi

# 4) HTTP Headers & SSL Certificate Collection
log_step "HTTP Headers & SSL Certificate Analysis" 5 $TOTAL_STEPS

http_results=()
cert_results=()

# Test direct domain access
log_substep "Testing direct domain access (${DOMAIN})" "INFO"
domain_start=$(start_timer)

if curl -sS -I --max-time 10 "https://${DOMAIN}/" > "${OUTDIR}/${DOMAIN}_https.headers" 2>/dev/null; then
    log_substep "HTTPS headers collected for ${DOMAIN}" "SUCCESS"
    http_results+=("${DOMAIN}:HTTPS:SUCCESS")
else
    log_substep "HTTPS connection to ${DOMAIN} failed" "WARN"
    http_results+=("${DOMAIN}:HTTPS:FAILED")
fi

if curl -sS -I --max-time 10 "http://${DOMAIN}/" > "${OUTDIR}/${DOMAIN}_http.headers" 2>/dev/null; then
    log_substep "HTTP headers collected for ${DOMAIN}" "SUCCESS"
    http_results+=("${DOMAIN}:HTTP:SUCCESS")
else
    log_substep "HTTP connection to ${DOMAIN} failed" "WARN"
    http_results+=("${DOMAIN}:HTTP:FAILED")
fi

# Test each IP with SNI
for i in "${!IPS[@]}"; do
    ip="${IPS[$i]}"
    current_cert=$((i + 1))
    
    log_substep "Testing IP $ip with SNI (${current_cert}/${#IPS[@]})" "INFO"
    
    # Collect headers with SNI
    if curl -sS -I --max-time 10 --resolve "${DOMAIN}:443:${ip}" "https://${DOMAIN}/" \
         > "${OUTDIR}/${DOMAIN}_on_${ip}.headers" 2>/dev/null; then
        log_substep "$ip: HTTPS headers with SNI collected" "SUCCESS"
        http_results+=("$ip:SNI_HTTPS:SUCCESS")
    else
        log_substep "$ip: HTTPS connection with SNI failed" "WARN"
        http_results+=("$ip:SNI_HTTPS:FAILED")
    fi
    
    # Collect SSL certificate
    cert_output="${OUTDIR}/${DOMAIN}_on_${ip}.cert.txt"
    if echo | timeout 10 openssl s_client -servername "${DOMAIN}" -connect "${ip}:443" 2>/dev/null \
        | openssl x509 -noout -issuer -subject -dates > "$cert_output" 2>/dev/null && [ -s "$cert_output" ]; then
        
        # Extract key certificate info
        issuer=$(grep "issuer=" "$cert_output" | sed 's/issuer=//' | head -1)
        expires=$(grep "notAfter=" "$cert_output" | sed 's/notAfter=//' | head -1)
        
        log_substep "$ip: SSL certificate collected" "SUCCESS"
        log_substep "    Issuer: ${issuer:0:50}..." "INFO"
        log_substep "    Expires: $expires" "INFO"
        cert_results+=("$ip:CERT:SUCCESS:$expires")
    else
        log_substep "$ip: SSL certificate collection failed" "WARN"
        cert_results+=("$ip:CERT:FAILED")
    fi
done

domain_duration=$(end_timer $domain_start)
log_summary "HTTP & SSL Analysis" "Domain + ${#IPS[@]} IPs" "Completed" "$domain_duration"

# 5) WHOIS Information Gathering
log_step "WHOIS Information Collection" 6 $TOTAL_STEPS

whois_results=()
whois_start=$(start_timer)

for i in "${!IPS[@]}"; do
    ip="${IPS[$i]}"
    current_whois=$((i + 1))
    
    log_substep "WHOIS lookup for $ip (${current_whois}/${#IPS[@]})" "INFO"
    
    if whois "${ip}" > "${OUTDIR}/${ip}.whois.txt" 2>/dev/null; then
        # Extract key WHOIS info
        org=$(grep -i "^org\|^organization" "${OUTDIR}/${ip}.whois.txt" | head -1 | cut -d: -f2- | xargs || echo "N/A")
        country=$(grep -i "^country" "${OUTDIR}/${ip}.whois.txt" | head -1 | cut -d: -f2- | xargs || echo "N/A")
        
        log_substep "$ip: WHOIS collected - $org, $country" "SUCCESS"
        whois_results+=("$ip:SUCCESS:$org|$country")
    else
        log_substep "$ip: WHOIS lookup failed" "WARN"
        whois_results+=("$ip:FAILED")
    fi
done

whois_duration=$(end_timer $whois_start)
log_summary "WHOIS Collection" "${#IPS[@]} IPs" "Completed" "$whois_duration"

# 6) Final Report Generation
printf "\n\033[1;35m╔════════════════════════════════════════════════╗\033[0m\n"
printf "\033[1;35m║                📋 FINAL REPORT                 ║\033[0m\n"
printf "\033[1;35m╚════════════════════════════════════════════════╝\033[0m\n"

SUMMARY="${OUTDIR}/SUMMARY.txt"
report_start=$(start_timer)

{
  echo "╔══════════════════════════════════════════════════════════════════╗"
  echo "║                    🎯 AIRTASKER RECON SUMMARY                     ║"
  echo "╠══════════════════════════════════════════════════════════════════╣"
  echo "║ Scan Time: $(date '+%Y-%m-%d %H:%M:%S %Z')                                ║"
  echo "║ Target Domain: ${DOMAIN}                               ║"
  echo "║ Scanned IPs: ${#IPS[@]} addresses                                     ║"
  echo "╚══════════════════════════════════════════════════════════════════╝"
  echo
  
  echo "🔍 [DNS RESOLUTION]"
  if [ -f "${OUTDIR}/dns_overview.txt" ]; then
      cat "${OUTDIR}/dns_overview.txt"
  else
      echo "DNS data not available"
  fi
  echo
  
  echo "📊 [SCAN RESULTS SUMMARY]"
  for result in "${scan_results[@]}"; do
      IFS=':' read -r ip status details <<< "$result"
      printf "  %-15s | %-10s | %s\n" "$ip" "$status" "$details"
  done
  echo
  
  echo "🌐 [WEB SERVICES SUMMARY]"  
  for result in "${web_results[@]}"; do
      IFS=':' read -r ip status details <<< "$result"
      printf "  %-15s | %-10s | %s\n" "$ip" "$status" "$details"
  done
  echo
  
  if [ "${full_results[0]}" != "SKIPPED:User choice" ]; then
      echo "🔎 [FULL PORT SCAN RESULTS]"
      for result in "${full_results[@]}"; do
          IFS=':' read -r ip status details <<< "$result"
          printf "  %-15s | %-10s | %s\n" "$ip" "$status" "$details"
      done
      echo
  fi
  
  echo "🔐 [SSL CERTIFICATE STATUS]"
  for result in "${cert_results[@]}"; do
      IFS=':' read -r ip type status expires <<< "$result"
      if [ "$status" = "SUCCESS" ]; then
          printf "  %-15s | %-8s | Expires: %s\n" "$ip" "$status" "$expires"
      else
          printf "  %-15s | %s\n" "$ip" "$status"
      fi
  done
  echo
  
  echo "📋 [DETAILED FINDINGS PER IP]"
  for ip in "${IPS[@]}"; do
    echo "─────────── $ip ───────────"
    
    [ -f "${OUTDIR}/${ip}_top200.gnmap" ] && {
        echo "📍 [Open Ports - Top 200]"
        grep -E "Ports: " "${OUTDIR}/${ip}_top200.gnmap" | sed 's/.*Ports: /  /'
        echo
    }
    
    [ -f "${OUTDIR}/${ip}_full.gnmap" ] && {
        echo "📍 [Open Ports - Full Scan]"
        grep -E "Ports: " "${OUTDIR}/${ip}_full.gnmap" | sed 's/.*Ports: /  /'
        echo
    }
    
    if [ -f "${OUTDIR}/${DOMAIN}_on_${ip}.headers" ]; then
      echo "🛡️  [Security Headers via IP+SNI]"
      awk 'BEGIN{IGNORECASE=1}
           /^strict-transport-security:/{print "  " $0}
           /^content-security-policy:/{print "  " $0}
           /^x-frame-options:/{print "  " $0}
           /^x-content-type-options:/{print "  " $0}
           /^referrer-policy:/{print "  " $0}
           /^permissions-policy:/{print "  " $0}
           /^set-cookie:/{print "  " $0}' "${OUTDIR}/${DOMAIN}_on_${ip}.headers"
      echo
    fi
    
    [ -f "${OUTDIR}/${DOMAIN}_on_${ip}.cert.txt" ] && {
        echo "📜 [SSL Certificate - SNI=${DOMAIN}]"
        sed 's/^/  /' "${OUTDIR}/${DOMAIN}_on_${ip}.cert.txt"
        echo
    }
  done
  
  echo "╔══════════════════════════════════════════════════════════════════╗"
  echo "║                        📋 SCAN COMPLETE                          ║"
  echo "╚══════════════════════════════════════════════════════════════════╝"
} > "$SUMMARY"

report_duration=$(end_timer $report_start)
log_summary "Report Generation" "Complete analysis" "Generated" "$report_duration"

# Final completion message
total_duration=$(end_timer $step_start)
printf "\n\033[1;32m🎉 RECON COMPLETED SUCCESSFULLY!\033[0m\n"
printf "\033[1;36m📁 Report Directory: %s\033[0m\n" "$OUTDIR"
printf "\033[1;36m📄 Summary Report: %s\033[0m\n" "$SUMMARY" 
printf "\033[1;36m⏱️  Total Duration: %s\033[0m\n" "$total_duration"

echo
green "💡 NEXT STEPS RECOMMENDATIONS:"
echo "  1. Review the summary report for key findings"
echo "  2. Check security headers for missing protections"
echo "  3. Verify SSL certificate expiration dates"
echo "  4. Investigate any unexpected open ports"
echo "  5. Consider deeper vulnerability scanning if authorized"
