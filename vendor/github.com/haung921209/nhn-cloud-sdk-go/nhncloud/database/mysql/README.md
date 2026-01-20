# MySQL API Guide

Complete guide for NHN Cloud RDS for MySQL v3.0 API.

## Table of Contents

- [Quick Start](#quick-start)
- [Instance Management](#instance-management)
- [Security Groups](#security-groups)
- [Parameter Groups](#parameter-groups)
- [Backups](#backups)
- [High Availability](#high-availability)
- [Read Replicas](#read-replicas)
- [DB Users & Schemas](#db-users--schemas)
- [Network Configuration](#network-configuration)
- [Monitoring](#monitoring)
- [Complete API Reference](#complete-api-reference)

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/haung921209/nhn-cloud-sdk-go/v2/nhncloud/database/mysql"
)

func main() {
    // Initialize client
    cfg := mysql.Config{
        Region:    "kr1",
        AppKey:    "your-mysql-app-key",
        AccessKey: "your-access-key",
        SecretKey: "your-secret-key",
    }
    
    client, err := mysql.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    // List instances
    ctx := context.Background()
    instances, err := client.ListInstances(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, inst := range instances.DBInstances {
        fmt.Printf("%s: %s (%s)\n", 
            inst.DBInstanceName, 
            inst.DBInstanceStatus,
            inst.DBVersion)
    }
}
```

---

## Instance Management

### List Instances

```go
instances, err := client.ListInstances(ctx)
if err != nil {
    log.Fatal(err)
}

for _, inst := range instances.DBInstances {
    fmt.Printf("ID: %s\n", inst.DBInstanceID)
    fmt.Printf("Name: %s\n", inst.DBInstanceName)
    fmt.Printf("Status: %s\n", inst.DBInstanceStatus)
    fmt.Printf("Version: %s\n", inst.DBVersion)
    fmt.Printf("Flavor: %s\n", inst.DBFlavorName)
    fmt.Printf("Storage: %s %dGB\n", inst.Storage.StorageType, inst.Storage.StorageSize)
}
```

### Get Instance Details

```go
instance, err := client.GetInstance(ctx, "instance-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Instance: %s\n", instance.DBInstance.DBInstanceName)
fmt.Printf("Network: %+v\n", instance.DBInstance.Network)
fmt.Printf("Backup: %+v\n", instance.DBInstance.Backup)
```

### Create Instance

```go
req := &mysql.CreateInstanceRequest{
    DBInstanceName: "my-production-db",
    DBFlavorID:     "m2.c4m8",
    DBVersion:      "MYSQL_V8043",
    DBUserName:     "admin",
    DBPassword:     "SecurePass123",
    DBPort:         ptr.Int(3306),
    ParameterGroupID: "default-mysql-8.0",
    Network: mysql.CreateInstanceNetworkConfig{
        SubnetID:         "subnet-abc123",
        AvailabilityZone: "kr-pub-a",
        UsePublicAccess:  ptr.Bool(false),
    },
    Storage: mysql.CreateInstanceStorageConfig{
        StorageType: "General SSD",
        StorageSize: 100,
    },
    Backup: mysql.CreateInstanceBackupConfig{
        BackupPeriod: 7,
        BackupSchedules: []mysql.CreateInstanceBackupSchedule{
            {
                BackupWndBgnTime:  "02:00:00",
                BackupWndDuration: "ONE_HOUR",
            },
        },
    },
    UseHighAvailability:   ptr.Bool(true),
    ReplicationMode:       "SEMISYNC",
    UseDeletionProtection: ptr.Bool(true),
}

result, err := client.CreateInstance(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Creating instance... Job ID: %s\n", result.JobID)
```

### Lifecycle Operations

```go
// Start instance
result, err := client.StartInstance(ctx, "instance-id")

// Stop instance
result, err := client.StopInstance(ctx, "instance-id")

// Restart instance
req := &mysql.RestartInstanceRequest{
    ExecuteBackup:     ptr.Bool(true),
    UseOnlineFailover: ptr.Bool(true),
}
result, err := client.RestartInstance(ctx, "instance-id", req)

// Delete instance
result, err := client.DeleteInstance(ctx, "instance-id")
```

---

## Security Groups

### List Security Groups

```go
groups, err := client.ListSecurityGroups(ctx)
if err != nil {
    log.Fatal(err)
}

for _, group := range groups.DBSecurityGroups {
    fmt.Printf("Group: %s\n", group.DBSecurityGroupName)
    for _, rule := range group.Rules {
        fmt.Printf("  Rule: %s from %s\n", rule.Port.PortType, rule.CIDR)
    }
}
```

### Create Security Group with Rules

```go
req := &mysql.CreateSecurityGroupRequest{
    DBSecurityGroupName: "web-server-access",
    Description:         "Allow MySQL access from web servers",
    Rules: []mysql.SecurityRule{
        {
            Description: "Web server subnet",
            Direction:   "INGRESS",
            EtherType:   "IPV4",
            Port: mysql.RulePort{
                PortType: "DB_PORT",
            },
            CIDR: "10.0.1.0/24",
        },
    },
}

result, err := client.CreateSecurityGroup(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created security group: %s\n", result.DBSecurityGroupID)
```

### Add Security Rule

```go
req := &mysql.CreateSecurityRuleRequest{
    Description: "Office network",
    Direction:   "INGRESS",
    EtherType:   "IPV4",
    Port: mysql.RulePort{
        PortType: "DB_PORT",
    },
    CIDR: "203.0.113.0/24",
}

result, err := client.CreateSecurityRule(ctx, "security-group-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Added rule: %s\n", result.RuleID)
```

---

## Parameter Groups

### List Available Parameters

```go
group, err := client.GetParameterGroup(ctx, "param-group-id")
if err != nil {
    log.Fatal(err)
}

for _, param := range group.ParameterGroup.Parameters {
    fmt.Printf("Parameter: %s\n", param.ParameterName)
    fmt.Printf("  Value: %s\n", param.Value)
    fmt.Printf("  Default: %s\n", param.DefaultValue)
    fmt.Printf("  Modifiable: %v\n", param.IsModifiable)
}
```

### Create Custom Parameter Group

```go
req := &mysql.CreateParameterGroupRequest{
    ParameterGroupName: "my-optimized-params",
    Description:        "Custom parameters for high-traffic DB",
    DBVersion:          "MYSQL_V8043",
}

result, err := client.CreateParameterGroup(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created parameter group: %s\n", result.ParameterGroupID)
```

### Modify Parameters

```go
req := &mysql.ModifyParametersRequest{
    ModifiedParameters: []mysql.ModifiedParameter{
        {
            ParameterID: "max_connections",
            Value:       "500",
        },
        {
            ParameterID: "innodb_buffer_pool_size",
            Value:       "8589934592", // 8GB
        },
    },
}

_, err := client.ModifyParameters(ctx, "param-group-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Parameters updated successfully")
```

---

## Backups

### List Backups

```go
backups, err := client.ListBackups(ctx, "instance-id")
if err != nil {
    log.Fatal(err)
}

for _, backup := range backups.Backups {
    fmt.Printf("Backup: %s\n", backup.BackupName)
    fmt.Printf("  Status: %s\n", backup.BackupStatus)
    fmt.Printf("  Size: %d bytes\n", backup.BackupSize)
    fmt.Printf("  Created: %s\n", backup.CreatedAt)
}
```

### Create Manual Backup

```go
req := &mysql.CreateBackupRequest{
    BackupName: "before-upgrade-backup",
}

result, err := client.CreateBackup(ctx, "instance-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Backup started. Job ID: %s\n", result.JobID)
```

### Restore from Backup

```go
req := &mysql.RestoreBackupRequest{
    DBInstanceName: "restored-database",
}

result, err := client.RestoreBackup(ctx, "backup-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Restore started. Job ID: %s\n", result.JobID)
```

### Export to Object Storage

```go
req := &mysql.ExportBackupRequest{
    TenantID:        "your-tenant-id",
    Username:        "your-username",
    Password:        "your-password",
    TargetContainer: "backups",
    ObjectPath:      "mysql/2024/backup-20240114.sql",
}

result, err := client.ExportBackup(ctx, "backup-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Export started. Job ID: %s\n", result.JobID)
```

---

## High Availability

### Enable HA

```go
req := &mysql.EnableHARequest{
    UseHighAvailability: true,
    PingInterval:        ptr.Int(3),
    ReplicationMode:     "SEMISYNC",
}

result, err := client.EnableHA(ctx, "instance-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("HA enabled. Job ID: %s\n", result.JobID)
```

### HA Operations

```go
// Pause HA (for maintenance)
result, err := client.PauseHA(ctx, "instance-id")

// Resume HA
result, err := client.ResumeHA(ctx, "instance-id")

// Repair HA (if broken)
result, err := client.RepairHA(ctx, "instance-id")

// Split HA (separate master and candidate)
result, err := client.SplitHA(ctx, "instance-id")

// Disable HA
result, err := client.DisableHA(ctx, "instance-id")
```

---

## Read Replicas

### Create Read Replica

```go
req := &mysql.CreateReplicaRequest{
    DBInstanceName: "my-read-replica-1",
}

result, err := client.CreateReplica(ctx, "master-instance-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Creating replica... Job ID: %s\n", result.JobID)
```

### Promote Replica to Master

```go
result, err := client.PromoteReplica(ctx, "replica-instance-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Promoting replica... Job ID: %s\n", result.JobID)
```

---

## DB Users & Schemas

### Create DB User

```go
req := &mysql.CreateDBUserRequest{
    DBUserName:           "app_user",
    DBPassword:           "SecurePass123",
    Host:                 "%",
    AuthorityType:        "CRUD",
    AuthenticationPlugin: "CACHING_SHA2",
    TLSOption:            "NONE",
}

result, err := client.CreateDBUser(ctx, "instance-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User created. Job ID: %s\n", result.JobID)
```

### Create Schema

```go
req := &mysql.CreateSchemaRequest{
    DBSchemaName: "my_application_db",
}

result, err := client.CreateSchema(ctx, "instance-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Schema created. Job ID: %s\n", result.JobID)
```

---

## Network Configuration

### Get Network Info

```go
networkInfo, err := client.GetNetworkInfo(ctx, "instance-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Availability Zone: %s\n", networkInfo.NetworkInfo.AvailabilityZone)
for _, endpoint := range networkInfo.NetworkInfo.EndPoints {
    fmt.Printf("  %s: %s (%s)\n", 
        endpoint.EndPointType, 
        endpoint.Domain, 
        endpoint.IPAddress)
}
```

### Enable Public Access

```go
req := &mysql.ModifyNetworkInfoRequest{
    UsePublicAccess: true,
}

result, err := client.ModifyNetworkInfo(ctx, "instance-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Public access enabled. Job ID: %s\n", result.JobID)
```

### Expand Storage

```go
req := &mysql.ModifyStorageInfoRequest{
    StorageSize:       200, // GB
    UseOnlineFailover: ptr.Bool(true),
}

result, err := client.ModifyStorageInfo(ctx, "instance-id", req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Storage expansion started. Job ID: %s\n", result.JobID)
```

---

## Monitoring

### List Log Files

```go
logs, err := client.ListLogFiles(ctx, "instance-id")
if err != nil {
    log.Fatal(err)
}

for _, log := range logs.LogFiles {
    fmt.Printf("Log: %s (%d bytes)\n", log.LogFileName, log.LogFileSize)
}
```

### Get Metrics

```go
metrics, err := client.ListMetrics(ctx)
if err != nil {
    log.Fatal(err)
}

for _, metric := range metrics.Metrics {
    fmt.Printf("Metric: %s (%s)\n", metric.MetricName, metric.Unit)
}
```

### Get Metric Statistics

```go
stats, err := client.GetMetricStatistics(
    ctx,
    "instance-id",
    "2024-01-14T00:00:00+09:00", // from
    "2024-01-14T23:59:59+09:00", // to
    300, // interval in seconds (5 minutes)
)
if err != nil {
    log.Fatal(err)
}

for _, stat := range stats.MetricStatistics {
    fmt.Printf("Metric: %s\n", stat.MetricName)
    for _, value := range stat.Values {
        fmt.Printf("  %s: %.2f\n", value.Timestamp, value.Value)
    }
}
```

---

## Complete API Reference

### Instance Management (9 APIs)
- `ListInstances(ctx)` - List all instances
- `GetInstance(ctx, instanceID)` - Get instance details
- `CreateInstance(ctx, req)` - Create new instance
- `ModifyInstance(ctx, instanceID, req)` - Modify instance
- `DeleteInstance(ctx, instanceID)` - Delete instance
- `StartInstance(ctx, instanceID)` - Start stopped instance
- `StopInstance(ctx, instanceID)` - Stop running instance
- `RestartInstance(ctx, instanceID, req)` - Restart instance
- `ForceRestartInstance(ctx, instanceID)` - Force restart

### Security Groups (8 APIs)
- `ListSecurityGroups(ctx)`
- `GetSecurityGroup(ctx, groupID)`
- `CreateSecurityGroup(ctx, req)`
- `UpdateSecurityGroup(ctx, groupID, req)`
- `DeleteSecurityGroup(ctx, groupID)`
- `CreateSecurityRule(ctx, groupID, req)`
- `UpdateSecurityRule(ctx, groupID, ruleID, req)`
- `DeleteSecurityRule(ctx, groupID, ruleID)`

### Parameter Groups (8 APIs)
- `ListParameterGroups(ctx)`
- `GetParameterGroup(ctx, groupID)`
- `CreateParameterGroup(ctx, req)`
- `CopyParameterGroup(ctx, groupID, req)`
- `UpdateParameterGroup(ctx, groupID, req)`
- `ModifyParameters(ctx, groupID, req)`
- `ResetParameterGroup(ctx, groupID)`
- `DeleteParameterGroup(ctx, groupID)`

### Backups (6 APIs)
- `ListBackups(ctx, instanceID)`
- `CreateBackup(ctx, instanceID, req)`
- `BackupToObjectStorage(ctx, instanceID, req)`
- `RestoreBackup(ctx, backupID, req)`
- `ExportBackup(ctx, backupID, req)`
- `DeleteBackup(ctx, backupID)`

### High Availability (5 APIs)
- `EnableHA(ctx, instanceID, req)`
- `DisableHA(ctx, instanceID)`
- `PauseHA(ctx, instanceID)`
- `ResumeHA(ctx, instanceID)`
- `RepairHA(ctx, instanceID)`
- `SplitHA(ctx, instanceID)`

### Read Replicas (2 APIs)
- `CreateReplica(ctx, instanceID, req)`
- `PromoteReplica(ctx, instanceID)`

### DB Users & Schemas (7 APIs)
- `ListDBUsers(ctx, instanceID)`
- `CreateDBUser(ctx, instanceID, req)`
- `UpdateDBUser(ctx, instanceID, userID, req)`
- `DeleteDBUser(ctx, instanceID, userID)`
- `ListSchemas(ctx, instanceID)`
- `CreateSchema(ctx, instanceID, req)`
- `DeleteSchema(ctx, instanceID, schemaID)`

### Network (4 APIs)
- `GetNetworkInfo(ctx, instanceID)`
- `ModifyNetworkInfo(ctx, instanceID, req)`
- `ModifyStorageInfo(ctx, instanceID, req)`
- `ModifyDeletionProtection(ctx, instanceID, req)`

### Notifications (5 APIs)
- `ListNotificationGroups(ctx)`
- `GetNotificationGroup(ctx, groupID)`
- `CreateNotificationGroup(ctx, req)`
- `UpdateNotificationGroup(ctx, groupID, req)`
- `DeleteNotificationGroup(ctx, groupID)`

### Monitoring (3 APIs)
- `ListLogFiles(ctx, instanceID)`
- `ListMetrics(ctx)`
- `GetMetricStatistics(ctx, instanceID, from, to, interval)`

### Reference Data (4 APIs)
- `ListFlavors(ctx)` - List available instance types
- `ListVersions(ctx)` - List MySQL versions
- `ListStorageTypes(ctx)` - List storage types
- `ListSubnets(ctx)` - List available subnets

---

## Helper Package

For pointer conversions:

```go
import "github.com/haung921209/nhn-cloud-sdk-go/v2/internal/ptr"

// Use ptr.Int, ptr.Bool, ptr.String for optional fields
req := &mysql.CreateInstanceRequest{
    DBPort: ptr.Int(3306),
    UseHighAvailability: ptr.Bool(true),
}
```

---

## Official Documentation

- [MySQL API Guide (Korean)](https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/)
- [MySQL Console Guide](https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/console-guide/)

---

**Total APIs**: 64  
**Version**: v3.0  
**Status**: ✅ Complete
