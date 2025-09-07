#!/usr/bin/env bash
# filename: airtasker_recon_burp.sh (Enhanced with Burp Suite Integration)
set -euo pipefail

# Cleanup function for proper exit
cleanup() {
    if [ -n "${BURP_PID:-}" ]; then
        echo "Cleaning up Burp Suite process..."
        kill -TERM "$BURP_PID" 2>/dev/null || true
        sleep 2
        kill -KILL "$BURP_PID" 2>/dev/null || true
    fi
}

# Set up trap for cleanup on exit
trap cleanup EXIT INT TERM

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

# Security disclaimer
printf "\n\033[1;31m╔══════════════════════════════════════════════════════╗\033[0m\n"
printf "\033[1;31m║                    ⚠️  SECURITY NOTICE                ║\033[0m\n"
printf "\033[1;31m╠══════════════════════════════════════════════════════╣\033[0m\n"
printf "\033[1;31m║ This tool performs network reconnaissance and        ║\033[0m\n" 
printf "\033[1;31m║ vulnerability scanning. Only use on systems you     ║\033[0m\n"
printf "\033[1;31m║ own or have explicit authorization to test.         ║\033[0m\n"
printf "\033[1;31m║                                                      ║\033[0m\n"
printf "\033[1;31m║ Unauthorized scanning is illegal and unethical.     ║\033[0m\n"
printf "\033[1;31m╚══════════════════════════════════════════════════════╝\033[0m\n"

read -p "Do you have authorization to scan $DOMAIN? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    red "❌ Scanning aborted - Authorization required"
    exit 1
fi

printf "\033[32m✓ Authorization confirmed - Proceeding with scan\033[0m\n"

# Initialization
printf "\n\033[1;36m╔══════════════════════════════════════════════════════╗\033[0m\n"
printf "\033[1;36m║        🔍 AIRTASKER RECON TOOL v2.0 + BURP          ║\033[0m\n"
printf "\033[1;36m╠══════════════════════════════════════════════════════╣\033[0m\n"
printf "\033[1;36m║ Target Domain: %-37s ║\033[0m\n" "$DOMAIN"
printf "\033[1;36m║ Target IPs: %-40s ║\033[0m\n" "${#IPS[@]} addresses"
printf "\033[1;36m║ Output Directory: %-32s ║\033[0m\n" "$OUTDIR"
printf "\033[1;36m║ Scan Start: %-37s ║\033[0m\n" "$(date '+%Y-%m-%d %H:%M:%S %Z')"
printf "\033[1;36m╚══════════════════════════════════════════════════════╝\033[0m\n"

# Tool dependency check
log_substep "Checking required tools..." "INFO"
missing_tools=()
optional_tools=()

# Essential tools
for tool in nmap curl whois dig openssl; do 
    if have "$tool"; then
        log_substep "$tool: Available" "SUCCESS"
    else
        log_substep "$tool: MISSING" "ERROR"
        missing_tools+=("$tool")
    fi
done

# Check for jq (needed for Burp Suite API)
if have "jq"; then
    log_substep "jq: Available (needed for Burp Suite)" "SUCCESS"
else
    log_substep "jq: MISSING (install for Burp Suite support: brew install jq)" "WARN"
fi

# Burp Suite Professional check
BURP_JAR=""
BURP_AVAILABLE=false

# Common Burp Suite locations
BURP_PATHS=(
    "/Applications/Burp Suite Professional.app/Contents/java/app/burpsuite_pro.jar"
    "/opt/BurpSuitePro/burpsuite_pro.jar"
    "$HOME/BurpSuitePro/burpsuite_pro.jar"
    "/usr/local/BurpSuitePro/burpsuite_pro.jar"
    "$(find /Applications -name "burpsuite_pro.jar" 2>/dev/null | head -1)"
)

for burp_path in "${BURP_PATHS[@]}"; do
    if [ -f "$burp_path" ] && [ -n "$burp_path" ]; then
        BURP_JAR="$burp_path"
        BURP_AVAILABLE=true
        log_substep "Burp Suite Pro: Available at $burp_path" "SUCCESS"
        break
    fi
done

if [ "$BURP_AVAILABLE" = false ]; then
    log_substep "Burp Suite Pro: NOT FOUND (optional - advanced web scanning disabled)" "WARN"
    log_substep "    Install Burp Pro for enhanced web vulnerability scanning" "INFO"
fi

# Java check for Burp Suite
if [ "$BURP_AVAILABLE" = true ]; then
    if have "java"; then
        JAVA_VERSION=$(java -version 2>&1 | head -1 | cut -d'"' -f2)
        log_substep "Java: Available (version $JAVA_VERSION)" "SUCCESS"
    else
        log_substep "Java: MISSING (required for Burp Suite)" "ERROR"
        BURP_AVAILABLE=false
    fi
fi

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

# Burp Suite helper functions
start_burp_headless() {
    local burp_port="$1"
    local project_file="$2"
    
    if [ "$BURP_AVAILABLE" = false ]; then
        return 1
    fi
    
    log_substep "Starting Burp Suite in headless mode (port $burp_port)" "INFO"
    
    # Create Burp project file
    java -jar "$BURP_JAR" --project-file="$project_file" \
         --headless --listen-port="$burp_port" \
         > "${OUTDIR}/burp_startup.log" 2>&1 &
    
    BURP_PID=$!
    sleep 10  # Allow Burp to start
    
    # Test if Burp is responding
    if curl -s "http://127.0.0.1:$burp_port" > /dev/null 2>&1; then
        log_substep "Burp Suite started successfully (PID: $BURP_PID)" "SUCCESS"
        return 0
    else
        log_substep "Burp Suite failed to start" "ERROR"
        return 1
    fi
}

stop_burp() {
    if [ -n "$BURP_PID" ]; then
        log_substep "Stopping Burp Suite (PID: $BURP_PID)" "INFO"
        kill -TERM "$BURP_PID" 2>/dev/null || true
        sleep 3
        kill -KILL "$BURP_PID" 2>/dev/null || true
    fi
}

# Burp REST API functions
burp_spider_scan() {
    local target_url="$1"
    local burp_port="$2"
    
    log_substep "Starting spider scan on $target_url" "INFO"
    
    # Start spider via REST API
    local task_id=$(curl -s -X POST "http://127.0.0.1:$burp_port/v0.1/spider" \
        -H "Content-Type: application/json" \
        -d "{\"baseUrl\": \"$target_url\"}" | jq -r '.taskId' 2>/dev/null || echo "")
    
    if [ -n "$task_id" ] && [ "$task_id" != "null" ]; then
        log_substep "Spider scan started (Task ID: $task_id)" "SUCCESS"
        
        # Wait for spider completion (max 5 minutes)
        local timeout=300
        local elapsed=0
        
        while [ $elapsed -lt $timeout ]; do
            local status=$(curl -s "http://127.0.0.1:$burp_port/v0.1/spider/$task_id" | jq -r '.status' 2>/dev/null || echo "running")
            
            if [ "$status" = "finished" ]; then
                log_substep "Spider scan completed for $target_url" "SUCCESS"
                return 0
            fi
            
            sleep 10
            elapsed=$((elapsed + 10))
        done
        
        log_substep "Spider scan timed out for $target_url" "WARN"
        return 1
    else
        log_substep "Failed to start spider scan for $target_url" "ERROR"
        return 1
    fi
}

burp_active_scan() {
    local target_url="$1"
    local burp_port="$2"
    
    log_substep "Starting active vulnerability scan on $target_url" "INFO"
    
    # Start active scan via REST API
    local task_id=$(curl -s -X POST "http://127.0.0.1:$burp_port/v0.1/scan" \
        -H "Content-Type: application/json" \
        -d "{\"baseUrl\": \"$target_url\"}" | jq -r '.taskId' 2>/dev/null || echo "")
    
    if [ -n "$task_id" ] && [ "$task_id" != "null" ]; then
        log_substep "Active scan started (Task ID: $task_id)" "SUCCESS"
        
        # Wait for scan completion (max 10 minutes)
        local timeout=600
        local elapsed=0
        
        while [ $elapsed -lt $timeout ]; do
            local status=$(curl -s "http://127.0.0.1:$burp_port/v0.1/scan/$task_id" | jq -r '.status' 2>/dev/null || echo "running")
            
            if [ "$status" = "finished" ]; then
                log_substep "Active scan completed for $target_url" "SUCCESS"
                
                # Get scan results
                local issues=$(curl -s "http://127.0.0.1:$burp_port/v0.1/scan/$task_id/issues" | jq length 2>/dev/null || echo "0")
                log_substep "Found $issues potential vulnerabilities" "INFO"
                
                return 0
            fi
            
            sleep 15
            elapsed=$((elapsed + 15))
        done
        
        log_substep "Active scan timed out for $target_url" "WARN"
        return 1
    else
        log_substep "Failed to start active scan for $target_url" "ERROR"
        return 1
    fi
}

generate_burp_report() {
    local burp_port="$1"
    local output_file="$2"
    
    log_substep "Generating Burp Suite security report" "INFO"
    
    # Get all issues via REST API
    local issues_json="${OUTDIR}/burp_issues.json"
    if curl -s "http://127.0.0.1:$burp_port/v0.1/issues" > "$issues_json"; then
        local issue_count=$(jq length "$issues_json" 2>/dev/null || echo "0")
        
        if [ "$issue_count" -gt 0 ]; then
            # Generate HTML report
            {
                echo "<html><head><title>Burp Suite Security Report - $DOMAIN</title></head><body>"
                echo "<h1>🛡️ Burp Suite Security Scan Report</h1>"
                echo "<h2>Target: $DOMAIN</h2>"
                echo "<h3>Scan Date: $(date '+%Y-%m-%d %H:%M:%S')</h3>"
                echo "<h3>Total Issues Found: $issue_count</h3>"
                echo "<hr>"
                
                # Parse and display issues
                jq -r '.[] | "
                <div style=\"border:1px solid #ccc; margin:10px; padding:10px;\">
                <h4 style=\"color:red;\">🔴 " + .name + "</h4>
                <p><strong>Severity:</strong> " + .severity + "</p>
                <p><strong>URL:</strong> " + .origin + "</p>
                <p><strong>Description:</strong> " + .description + "</p>
                </div>"' "$issues_json" 2>/dev/null
                
                echo "</body></html>"
            } > "$output_file"
            
            log_substep "Security report generated: $output_file ($issue_count issues)" "SUCCESS"
        else
            echo "No security issues detected by Burp Suite" > "$output_file"
            log_substep "No vulnerabilities found" "SUCCESS"
        fi
        
        return 0
    else
        log_substep "Failed to generate Burp report" "ERROR"
        return 1
    fi
}

# 0) DNS / PTR Analysis
TOTAL_STEPS=8
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

# 5) Burp Suite Web Application Security Testing
if [ "$BURP_AVAILABLE" = true ]; then
    log_step "Burp Suite Web Application Security Scanning" 6 $TOTAL_STEPS
    
    burp_results=()
    BURP_PORT=8080
    BURP_PROJECT="${OUTDIR}/airtasker_burp.burp"
    
    # Start Burp Suite in headless mode
    if start_burp_headless "$BURP_PORT" "$BURP_PROJECT"; then
        
        # Test HTTPS targets first
        working_targets=()
        
        # Test direct domain access
        log_substep "Testing $DOMAIN accessibility for Burp scanning" "INFO"
        if curl -s -I --max-time 10 "https://${DOMAIN}/" > /dev/null 2>&1; then
            working_targets+=("https://${DOMAIN}")
            log_substep "Added https://${DOMAIN} to scan targets" "SUCCESS"
        fi
        
        if curl -s -I --max-time 10 "http://${DOMAIN}/" > /dev/null 2>&1; then
            working_targets+=("http://${DOMAIN}")
            log_substep "Added http://${DOMAIN} to scan targets" "SUCCESS"
        fi
        
        # Test working IPs
        for ip in "${IPS[@]}"; do
            if curl -s -I --max-time 10 --resolve "${DOMAIN}:443:${ip}" "https://${DOMAIN}/" > /dev/null 2>&1; then
                working_targets+=("https://${ip}")
                log_substep "Added https://$ip to scan targets" "SUCCESS"
            fi
        done
        
        if [ ${#working_targets[@]} -gt 0 ]; then
            log_substep "Found ${#working_targets[@]} accessible web targets for scanning" "SUCCESS"
            
            # Phase 1: Spider/Crawler scan
            for target in "${working_targets[@]}"; do
                spider_start=$(start_timer)
                log_substep "Web crawling: $target" "INFO"
                
                if burp_spider_scan "$target" "$BURP_PORT"; then
                    spider_duration=$(end_timer $spider_start)
                    log_summary "Web Crawler" "$target" "Completed" "$spider_duration"
                    burp_results+=("$target:SPIDER:SUCCESS")
                else
                    spider_duration=$(end_timer $spider_start)
                    log_summary "Web Crawler" "$target" "Failed" "$spider_duration"
                    burp_results+=("$target:SPIDER:FAILED")
                fi
            done
            
            # Phase 2: Active vulnerability scanning (ask user for permission)
            printf "\n"
            yellow "⚠️  ACTIVE VULNERABILITY SCANNING WARNING"
            yellow "    This will perform intrusive security tests against $DOMAIN"
            yellow "    Only proceed if you have explicit authorization to test"
            yellow "    Active scanning can take 10-15 minutes per target"
            echo
            read -p "Proceed with active vulnerability scanning? (y/N): " -n 1 -r -t 15
            echo
            
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                log_substep "User authorized active scanning" "SUCCESS"
                
                for target in "${working_targets[@]}"; do
                    vuln_start=$(start_timer)
                    log_substep "Vulnerability scanning: $target" "INFO"
                    
                    if burp_active_scan "$target" "$BURP_PORT"; then
                        vuln_duration=$(end_timer $vuln_start)
                        log_summary "Vulnerability Scan" "$target" "Completed" "$vuln_duration"
                        burp_results+=("$target:VULN_SCAN:SUCCESS")
                    else
                        vuln_duration=$(end_timer $vuln_start)
                        log_summary "Vulnerability Scan" "$target" "Failed/Timeout" "$vuln_duration"
                        burp_results+=("$target:VULN_SCAN:TIMEOUT")
                    fi
                done
            else
                log_substep "Active scanning skipped by user" "WARN"
                burp_results+=("ACTIVE_SCAN:SKIPPED:User choice")
            fi
            
            # Generate security report
            burp_report_file="${OUTDIR}/burp_security_report.html"
            generate_burp_report "$BURP_PORT" "$burp_report_file"
            
        else
            log_substep "No accessible web targets found for Burp scanning" "WARN"
            burp_results+=("NO_TARGETS:All targets inaccessible")
        fi
        
        # Clean up Burp Suite
        stop_burp
        
    else
        log_substep "Failed to start Burp Suite - skipping web security tests" "ERROR"
        burp_results+=("BURP_STARTUP:FAILED")
    fi
    
else
    log_step "Burp Suite Web Security Testing (SKIPPED - Not Available)" 6 $TOTAL_STEPS
    log_substep "Install Burp Suite Professional for advanced web security testing" "WARN"
    burp_results=("BURP_NOT_AVAILABLE:Install Burp Suite Pro")
fi

# 6) WHOIS Information Gathering
log_step "WHOIS Information Collection" 7 $TOTAL_STEPS

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

# 7) Final Report Generation
log_step "Comprehensive Security Report Generation" 8 $TOTAL_STEPS

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
  
  echo "🛡️  [BURP SUITE SECURITY SCAN RESULTS]"
  if [ ${#burp_results[@]} -gt 0 ]; then
      for result in "${burp_results[@]}"; do
          IFS=':' read -r target action status details <<< "$result"
          if [ "$action" = "SPIDER" ]; then
              printf "  🕷️  Spider: %-25s | %s\n" "$target" "$status"
          elif [ "$action" = "VULN_SCAN" ]; then
              printf "  🔍 Vuln Scan: %-21s | %s\n" "$target" "$status"
          else
              printf "  %-15s | %-10s | %s\n" "$target" "$action" "$status"
          fi
      done
  else
      echo "  Burp Suite not available or skipped"
  fi
  
  # Add link to detailed Burp report if it exists
  if [ -f "${OUTDIR}/burp_security_report.html" ]; then
      echo "  📄 Detailed security report: ${OUTDIR}/burp_security_report.html"
  fi
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
echo "  1. Review the comprehensive summary report for all findings"
echo "  2. Check security headers for missing protections (HSTS, CSP, etc.)"
echo "  3. Verify SSL certificate expiration dates and issuers"
echo "  4. Investigate any unexpected open ports or services"

if [ -f "${OUTDIR}/burp_security_report.html" ]; then
    echo "  5. 🛡️  PRIORITY: Review Burp Suite security report for vulnerabilities"
    echo "     Report location: ${OUTDIR}/burp_security_report.html"
    echo "  6. Address any HIGH/CRITICAL vulnerabilities immediately"
    echo "  7. Consider implementing missing security controls"
else
    echo "  5. Install Burp Suite Professional for comprehensive web security testing"
    echo "  6. Consider automated security scanning in CI/CD pipeline"
fi

echo "  8. Set up monitoring for configuration changes"
echo "  9. Schedule regular security assessments"
echo "  10. Document and track remediation efforts"

if [ "$BURP_AVAILABLE" = true ]; then
    echo
    green "🔒 BURP SUITE SECURITY RECOMMENDATIONS:"
    echo "  • Enable passive scanning in production environments"
    echo "  • Integrate security testing into development workflow"
    echo "  • Set up automated vulnerability monitoring"
    echo "  • Consider Burp Enterprise for continuous scanning"
fi
