#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
多账号/多区域资产枚举（零噪音）
- Organizations 多账号遍历（可关）
- 区域自动发现
- 资源覆盖：Route53, ELB/ALB/NLB, API GW v1/v2, CloudFront, Global Accelerator,
            S3 静态站, EC2 EIP/公网, ACM 证书, 通过 ELB 标签识别 EKS 暴露型 LB
输出：
  - /tmp/asi_assets.ndjson    # 逐资源 NDJSON（便于 Athena/SIEM）
  - /tmp/asi_hosts.txt        # 可解析域名/IP 汇总（供后续 TLS/HTTP 只读探测）
  - /tmp/asi_summary.csv      # 关键字段汇总
"""
import os, json, csv, sys, time, boto3, botocore
from botocore.config import Config
from datetime import datetime, timezone

ASSUME_ROLE = os.getenv("ROLE_TO_ASSUME")           # e.g. "arn:aws:iam::<member-account-id>:role/SecAuditReadOnly"
ENUM_ORGS   = os.getenv("ENUM_ORGS", "false") == "true"
OUT_NDJSON  = "/tmp/asi_assets.ndjson"
OUT_HOSTS   = "/tmp/asi_hosts.txt"
OUT_CSV     = "/tmp/asi_summary.csv"

cfg = Config(retries={"max_attempts": 6}, read_timeout=10, connect_timeout=5)

def assume(role_arn):
    sts = boto3.client("sts", config=cfg)
    resp = sts.assume_role(RoleArn=role_arn, RoleSessionName="asi-enum")
    creds = resp["Credentials"]
    return boto3.Session(
        aws_access_key_id=creds["AccessKeyId"],
        aws_secret_access_key=creds["SecretAccessKey"],
        aws_session_token=creds["SessionToken"],
    )

def list_accounts(sess):
    if not ENUM_ORGS:
        return [sess]
    org = sess.client("organizations", config=cfg)
    sessions = []
    paginator = org.get_paginator("list_accounts")
    for page in paginator.paginate():
        for acc in page["Accounts"]:
            if acc["Status"] == "SUSPENDED": continue
            role = ASSUME_ROLE.replace("<member-account-id>", acc["Id"]) if "<member-account-id>" in ASSUME_ROLE else ASSUME_ROLE
            try:
                sessions.append(assume(role))
            except botocore.exceptions.ClientError:
                continue
    return sessions

def all_regions(sess):
    ec2 = sess.client("ec2", config=cfg)
    return [r["RegionName"] for r in ec2.describe_regions(AllRegions=True)["Regions"] if r["OptInStatus"] in ("opt-in-not-required", "opted-in")]

def write_ndjson(obj, f):
    f.write(json.dumps(obj, ensure_ascii=False) + "\n")

def append_host(h, bag):
    if h and isinstance(h, str):
        bag.add(h.strip().rstrip("."))

def route53(sess, nd, hosts):
    r53 = sess.client("route53", config=cfg)
    for hz in r53.list_hosted_zones()["HostedZones"]:
        if hz["Config"].get("PrivateZone"): continue
        zid = hz["Id"]
        pag = r53.get_paginator("list_resource_record_sets")
        for page in pag.paginate(HostedZoneId=zid):
            for rr in page["ResourceRecordSets"]:
                if rr["Type"] in ("A","AAAA","CNAME","ALIAS"):
                    name = rr["Name"].rstrip(".")
                    append_host(name, hosts)
                    nd.write(json.dumps({"kind":"route53", "name":name, "type":rr["Type"], "ttl":rr.get("TTL")})+"\n")

def elb_all(sess, region, nd, hosts):
    elb = sess.client("elbv2", config=cfg, region_name=region)
    for page in elb.get_paginator("describe_load_balancers").paginate():
        for lb in page["LoadBalancers"]:
            dns = lb["DNSName"]
            append_host(dns, hosts)
            # 标记可能来自 EKS 的 Service
            tags = {}
            try:
                t = elb.describe_tags(ResourceArns=[lb["LoadBalancerArn"]])["TagDescriptions"][0]["Tags"]
                tags = {i["Key"]: i["Value"] for i in t}
            except: pass
            nd.write(json.dumps({"kind":"elb", "dns":dns, "scheme":lb["Scheme"], "type":lb["Type"], "region":region, "tags":tags})+"\n")

def api_gateway(sess, region, nd, hosts):
    apig = sess.client("apigateway", config=cfg, region_name=region)     # v1
    try:
        for dom in apig.get_domain_names().get("items", []):
            append_host(dom["domainName"], hosts)
            write_ndjson({"kind":"apigw_v1_domain","domain":dom["domainName"],"region":region}, nd)
    except: pass
    apiv2 = sess.client("apigatewayv2", config=cfg, region_name=region)  # v2
    try:
        for page in apiv2.get_paginator("get_domain_names").paginate():
            for item in page.get("Items", []):
                append_host(item["DomainName"], hosts)
                write_ndjson({"kind":"apigw_v2_domain","domain":item["DomainName"],"region":region}, nd)
    except: pass

def cloudfront(sess, nd, hosts):
    cf = sess.client("cloudfront", config=cfg)
    try:
        for page in cf.get_paginator("list_distributions").paginate():
            for d in page.get("DistributionList", {}).get("Items", []):
                dom = d["DomainName"]                          # dxxxxx.cloudfront.net
                append_host(dom, hosts)
                for alt in d.get("Aliases", {}).get("Items", []):
                    append_host(alt, hosts)
                write_ndjson({"kind":"cloudfront","domain":dom,"alts":d.get("Aliases",{}).get("Items",[]),"oac":bool(d.get("OriginAccessControlId"))}, nd)
    except: pass

def global_acc(sess, nd, hosts):
    ga = sess.client("globalaccelerator", config=cfg)
    try:
        for arn in [a["AcceleratorArn"] for a in ga.list_accelerators().get("Accelerators", [])]:
            d = ga.describe_accelerator(AcceleratorArn=arn)["Accelerator"]
            dns = d["DnsName"]
            append_host(dns, hosts)
            write_ndjson({"kind":"global_accelerator","dns":dns,"enabled":d["Enabled"]}, nd)
    except: pass

def s3_web(sess, nd, hosts):
    s3 = sess.client("s3", config=cfg)
    for b in s3.list_buckets().get("Buckets", []):
        name = b["Name"]
        try:
            loc = s3.get_bucket_location(Bucket=name)["LocationConstraint"] or "us-east-1"
            # 静态网站端点（按区拼接）
            web_ep = f"{name}.s3-website-{loc}.amazonaws.com" if loc!="us-east-1" else f"{name}.s3-website-us-east-1.amazonaws.com"
            # 若配置了网站托管，则纳入
            s3.get_bucket_website(Bucket=name)
            append_host(web_ep, hosts)
            write_ndjson({"kind":"s3_website","bucket":name,"endpoint":web_ep,"region":loc}, nd)
        except: continue

def ec2_public(sess, region, nd, hosts):
    ec2 = sess.client("ec2", config=cfg, region_name=region)
    paginator = ec2.get_paginator("describe_instances")
    for page in paginator.paginate(Filters=[{"Name":"instance-state-name","Values":["running"]}]):
        for r in page["Reservations"]:
            for i in r["Instances"]:
                pip = i.get("PublicIpAddress")
                if pip:
                    append_host(pip, hosts)
                    write_ndjson({"kind":"ec2_public","ip":pip,"region":region,"iid":i["InstanceId"]}, nd)
    for a in ec2.describe_addresses().get("Addresses", []):
        if a.get("PublicIp"):
            append_host(a["PublicIp"], hosts)
            write_ndjson({"kind":"eip","ip":a["PublicIp"],"region":region}, nd)

def acm(sess, region, nd):
    acm = sess.client("acm", config=cfg, region_name=region)
    try:
        for page in acm.get_paginator("list_certificates").paginate(CertificateStatuses=["ISSUED"]):
            for c in page["CertificateSummaryList"]:
                write_ndjson({"kind":"acm_cert","domain":c.get("DomainName"),"alt":c.get("SubjectAlternativeNameSummaries",[]),"region":region}, nd)
    except: pass

def main():
    master = boto3.Session()
    sessions = list_accounts(master)
    hosts = set()
    with open(OUT_NDJSON, "w") as nd:
        for s in sessions:
            # 全局服务
            route53(s, nd, hosts)
            cloudfront(s, nd, hosts)
            global_acc(s, nd, hosts)
            s3_web(s, nd, hosts)
            # 分区域服务
            for region in all_regions(s):
                elb_all(s, region, nd, hosts)
                api_gateway(s, region, nd, hosts)
                ec2_public(s, region, nd, hosts)
                acm(s, region, nd)

    # 输出 hosts 与 CSV 摘要
    with open(OUT_HOSTS, "w") as f:
        for h in sorted(hosts): f.write(h + "\n")

    # 简要 CSV（域名/IP & 类别）
    with open(OUT_CSV, "w", newline="") as c:
        w = csv.writer(c); w.writerow(["artifact","kind"])
        for line in open(OUT_NDJSON):
            obj = json.loads(line)
            if obj["kind"] in ("eip","ec2_public"):
                w.writerow([obj["ip"], obj["kind"]])
            elif "dns" in obj:
                w.writerow([obj["dns"], obj["kind"]])
            elif "domain" in obj:
                w.writerow([obj["domain"], obj["kind"]])

    print(f"[+] NDJSON: {OUT_NDJSON}\n[+] HOSTS: {OUT_HOSTS}\n[+] CSV: {OUT_CSV}")

if __name__ == "__main__":
    main()
