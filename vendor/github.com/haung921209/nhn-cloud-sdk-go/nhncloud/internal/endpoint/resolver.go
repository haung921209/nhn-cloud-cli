package endpoint

import (
	"fmt"
	"strings"
)

type Service string

const (
	ServiceIAM           Service = "iam"
	ServiceCompute       Service = "compute"
	ServiceRDSMySQL      Service = "rds-mysql"
	ServiceRDSMariaDB    Service = "rds-mariadb"
	ServiceRDSPostgreSQL Service = "rds-postgresql"
	ServiceVPC           Service = "vpc"
	ServiceSecurityGroup Service = "security-group"
	ServiceFloatingIP    Service = "floating-ip"
	ServiceLoadBalancer  Service = "load-balancer"
	ServiceBlockStorage  Service = "block-storage"
	ServiceObjectStorage Service = "object-storage"
	ServiceNKS           Service = "nks"
	ServiceNCR           Service = "ncr"
	ServiceNCS           Service = "ncs"
)

func Resolve(service Service, region string) string {
	region = strings.ToLower(region)

	switch service {
	case ServiceIAM:
		return "https://core.api.nhncloudservice.com"
	case ServiceCompute:
		return fmt.Sprintf("https://%s-api-instance-infrastructure.nhncloudservice.com", region)
	case ServiceRDSMySQL:
		return fmt.Sprintf("https://%s-rds-mysql.api.nhncloudservice.com/v3.0", region)
	case ServiceRDSMariaDB:
		return fmt.Sprintf("https://%s-rds-mariadb.api.nhncloudservice.com/v3.0", region)
	case ServiceRDSPostgreSQL:
		// Note: NHN Cloud uses "rds-postgres" (not "rds-postgresql") in the endpoint
		return fmt.Sprintf("https://%s-rds-postgres.api.nhncloudservice.com/v1.0", region)
	case ServiceVPC:
		return fmt.Sprintf("https://%s-api-network-infrastructure.nhncloudservice.com", region)
	case ServiceSecurityGroup:
		return fmt.Sprintf("https://%s-api-network-infrastructure.nhncloudservice.com", region)
	case ServiceFloatingIP:
		return fmt.Sprintf("https://%s-api-network-infrastructure.nhncloudservice.com", region)
	case ServiceLoadBalancer:
		return fmt.Sprintf("https://%s-api-network-infrastructure.nhncloudservice.com", region)
	case ServiceBlockStorage:
		return fmt.Sprintf("https://%s-api-block-storage-infrastructure.nhncloudservice.com", region)
	case ServiceObjectStorage:
		return fmt.Sprintf("https://%s-api-object-storage.nhncloudservice.com", region)
	case ServiceNKS:
		return fmt.Sprintf("https://%s-api-kubernetes-infrastructure.nhncloudservice.com", region)
	case ServiceNCR:
		return "https://kr1-ncr.api.nhncloudservice.com"
	case ServiceNCS:
		return fmt.Sprintf("https://%s-ncs.api.nhncloudservice.com/ncs/v1.0", region)
	default:
		return ""
	}
}

func ResolveWithAppKey(service Service, region, appKey string) string {
	base := Resolve(service, region)
	// RDS services use appkey in header (X-Tc-App-Key), not in URL path
	// Only non-RDS services (like NCR, NCS) may need appkey in URL
	switch service {
	case ServiceRDSMySQL, ServiceRDSMariaDB, ServiceRDSPostgreSQL:
		return base
	default:
		if appKey != "" {
			return base + "/appkeys/" + appKey
		}
		return base
	}
}
