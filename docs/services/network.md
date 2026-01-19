# Network Service Guide

This guide covers **VPC**, **Subnets**, and **Floating IPs**.

## 1. VPC & Subnets

### List VPCs
```bash
nhncloud network describe-vpcs
```

### List Subnets
Find the UUIDs required for creating Instances and RDS databases.
```bash
nhncloud network describe-subnets
```

## 2. Floating IPs (Public Access)

### Create Floating IP
Allocate a Public IP address.
```bash
nhncloud network create-floating-ip --floating-ip <ip-address> # Or auto-allocate
```

### List Floating IPs
```bash
nhncloud network describe-floating-ips
```
