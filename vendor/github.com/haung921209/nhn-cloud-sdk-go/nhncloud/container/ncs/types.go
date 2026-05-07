package ncs

type ResponseHeader struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

type Workload struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Description       string            `json:"description,omitempty"`
	Type              string            `json:"type"`
	Replicas          int               `json:"replicas"`
	AvailableReplicas int               `json:"availableReplicas"`
	Status            string            `json:"status"`
	Containers        []Container       `json:"containers,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
}

type Container struct {
	Name         string                `json:"name"`
	Image        string                `json:"image"`
	Command      []string              `json:"command,omitempty"`
	Args         []string              `json:"args,omitempty"`
	Ports        []ContainerPort       `json:"ports,omitempty"`
	Env          []EnvVar              `json:"env,omitempty"`
	Resources    *ResourceRequirements `json:"resources,omitempty"`
	VolumeMounts []VolumeMount         `json:"volumeMounts,omitempty"`
}

type ContainerPort struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty"`
	Requests ResourceList `json:"requests,omitempty"`
}

type ResourceList struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	GPU    string `json:"nvidia.com/gpu,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type ListWorkloadsOutput struct {
	Header    *ResponseHeader `json:"header"`
	Workloads []Workload      `json:"workloads"`
}

type GetWorkloadOutput struct {
	Header *ResponseHeader `json:"header"`
	Workload
}

type CreateWorkloadInput struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace,omitempty"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type"`
	Replicas    int               `json:"replicas,omitempty"`
	Containers  []Container       `json:"containers"`
	Labels      map[string]string `json:"labels,omitempty"`
	Volumes     []Volume          `json:"volumes,omitempty"`
}

type Volume struct {
	Name                  string                 `json:"name"`
	EmptyDir              *EmptyDirVolumeSource  `json:"emptyDir,omitempty"`
	ConfigMap             *ConfigMapVolumeSource `json:"configMap,omitempty"`
	Secret                *SecretVolumeSource    `json:"secret,omitempty"`
	PersistentVolumeClaim *PVCVolumeSource       `json:"persistentVolumeClaim,omitempty"`
}

type EmptyDirVolumeSource struct {
	Medium    string `json:"medium,omitempty"`
	SizeLimit string `json:"sizeLimit,omitempty"`
}

type ConfigMapVolumeSource struct {
	Name string `json:"name"`
}

type SecretVolumeSource struct {
	SecretName string `json:"secretName"`
}

type PVCVolumeSource struct {
	ClaimName string `json:"claimName"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type CreateWorkloadOutput struct {
	Header *ResponseHeader `json:"header"`
	Workload
}

type UpdateWorkloadInput struct {
	Description string      `json:"description,omitempty"`
	Replicas    int         `json:"replicas,omitempty"`
	Containers  []Container `json:"containers,omitempty"`
}

type Template struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Version     string      `json:"version"`
	Type        string      `json:"type"`
	IsPublic    bool        `json:"isPublic"`
	Containers  []Container `json:"containers,omitempty"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

type ListTemplatesOutput struct {
	Header    *ResponseHeader `json:"header"`
	Templates []Template      `json:"templates"`
}

type GetTemplateOutput struct {
	Header *ResponseHeader `json:"header"`
	Template
}

type Service struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Type       string            `json:"type"`
	ClusterIP  string            `json:"clusterIP,omitempty"`
	ExternalIP string            `json:"externalIP,omitempty"`
	Ports      []ServicePort     `json:"ports,omitempty"`
	Selector   map[string]string `json:"selector,omitempty"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  string            `json:"updatedAt"`
}

type ServicePort struct {
	Name       string `json:"name,omitempty"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	NodePort   int    `json:"nodePort,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
}

type ListServicesOutput struct {
	Header   *ResponseHeader `json:"header"`
	Services []Service       `json:"services"`
}

type GetServiceOutput struct {
	Header *ResponseHeader `json:"header"`
	Service
}

type CreateServiceInput struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace,omitempty"`
	Type      string            `json:"type"`
	Ports     []ServicePort     `json:"ports"`
	Selector  map[string]string `json:"selector,omitempty"`
}

type CreateServiceOutput struct {
	Header *ResponseHeader `json:"header"`
	Service
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Stream    string `json:"stream"`
}

type GetWorkloadLogsOutput struct {
	Header *ResponseHeader `json:"header"`
	Logs   []LogEntry      `json:"logs"`
}

type HealthProbe struct {
	Type                string           `json:"type"`
	HTTPGet             *HTTPGetAction   `json:"httpGet,omitempty"`
	TCPSocket           *TCPSocketAction `json:"tcpSocket,omitempty"`
	Exec                *ExecAction      `json:"exec,omitempty"`
	InitialDelaySeconds int              `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int              `json:"periodSeconds,omitempty"`
	TimeoutSeconds      int              `json:"timeoutSeconds,omitempty"`
	SuccessThreshold    int              `json:"successThreshold,omitempty"`
	FailureThreshold    int              `json:"failureThreshold,omitempty"`
}

type HTTPGetAction struct {
	Path   string `json:"path"`
	Port   int    `json:"port"`
	Scheme string `json:"scheme,omitempty"`
}

type TCPSocketAction struct {
	Port int `json:"port"`
}

type ExecAction struct {
	Command []string `json:"command"`
}

type HealthCheckConfig struct {
	LivenessProbe  *HealthProbe `json:"livenessProbe,omitempty"`
	ReadinessProbe *HealthProbe `json:"readinessProbe,omitempty"`
}

type HealthCheckStatus struct {
	ContainerName string `json:"containerName"`
	Liveness      string `json:"liveness"`
	Readiness     string `json:"readiness"`
	LastCheck     string `json:"lastCheck"`
}

type GetHealthCheckStatusOutput struct {
	Header *ResponseHeader     `json:"header"`
	Status []HealthCheckStatus `json:"status"`
}

type UpdateResourceLimitsInput struct {
	Resources *ResourceRequirements `json:"resources"`
}

// Event types for workload events
type Event struct {
	EventID    string `json:"eventId"`
	EventType  string `json:"eventType"` // Normal, Warning
	Reason     string `json:"reason"`
	Message    string `json:"message"`
	Count      int    `json:"count"`
	FirstTime  string `json:"firstTimestamp"`
	LastTime   string `json:"lastTimestamp"`
	Source     string `json:"source,omitempty"`
	ObjectKind string `json:"objectKind,omitempty"`
	ObjectName string `json:"objectName,omitempty"`
}

type GetWorkloadEventsOutput struct {
	Header *ResponseHeader `json:"header"`
	Events []Event         `json:"events"`
}

// Volume types for persistent volume management
type PersistentVolume struct {
	VolumeID    string `json:"volumeId"`
	Name        string `json:"name"`
	Size        int    `json:"size"` // GB
	Status      string `json:"status"`
	VolumeType  string `json:"volumeType"`
	StorageType string `json:"storageType,omitempty"`
	CreatedAt   string `json:"createdAt"`
	AttachedTo  string `json:"attachedTo,omitempty"`
}

type ListVolumesOutput struct {
	Header  *ResponseHeader    `json:"header"`
	Volumes []PersistentVolume `json:"volumes"`
}

type VolumeAttachInput struct {
	VolumeID  string `json:"volumeId"`
	MountPath string `json:"mountPath"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type AttachVolumeOutput struct {
	Header *ResponseHeader `json:"header"`
}

type ExecInput struct {
	ContainerName string   `json:"containerName"`
	Command       []string `json:"command"`
	Stdin         bool     `json:"stdin,omitempty"`
	Stdout        bool     `json:"stdout,omitempty"`
	Stderr        bool     `json:"stderr,omitempty"`
	TTY           bool     `json:"tty,omitempty"`
}

type ExecOutput struct {
	Header     *ResponseHeader `json:"header"`
	ExecID     string          `json:"execId"`
	Output     string          `json:"output,omitempty"`
	ExitCode   int             `json:"exitCode"`
	Error      string          `json:"error,omitempty"`
	StartedAt  string          `json:"startedAt"`
	FinishedAt string          `json:"finishedAt,omitempty"`
}

type ContainerStatus struct {
	ContainerName string `json:"containerName"`
	State         string `json:"state"`
	Ready         bool   `json:"ready"`
	RestartCount  int    `json:"restartCount"`
	Image         string `json:"image"`
	ContainerID   string `json:"containerId,omitempty"`
	StartedAt     string `json:"startedAt,omitempty"`
	FinishedAt    string `json:"finishedAt,omitempty"`
	ExitCode      int    `json:"exitCode,omitempty"`
	Reason        string `json:"reason,omitempty"`
	Message       string `json:"message,omitempty"`
}

type GetContainerStatusOutput struct {
	Header     *ResponseHeader   `json:"header"`
	Containers []ContainerStatus `json:"containers"`
}

type AutoScalingPolicy struct {
	MinReplicas                       int `json:"minReplicas"`
	MaxReplicas                       int `json:"maxReplicas"`
	TargetCPUUtilizationPercentage    int `json:"targetCpuUtilizationPercentage,omitempty"`
	TargetMemoryUtilizationPercentage int `json:"targetMemoryUtilizationPercentage,omitempty"`
	ScaleUpStabilizationSeconds       int `json:"scaleUpStabilizationSeconds,omitempty"`
	ScaleDownStabilizationSeconds     int `json:"scaleDownStabilizationSeconds,omitempty"`
}

type ConfigureAutoScalingInput struct {
	Enabled bool               `json:"enabled"`
	Policy  *AutoScalingPolicy `json:"policy,omitempty"`
}

type AutoScalingStatus struct {
	Enabled         bool               `json:"enabled"`
	CurrentReplicas int                `json:"currentReplicas"`
	DesiredReplicas int                `json:"desiredReplicas"`
	Policy          *AutoScalingPolicy `json:"policy,omitempty"`
	LastScaleTime   string             `json:"lastScaleTime,omitempty"`
	Conditions      []struct {
		Type    string `json:"type"`
		Status  string `json:"status"`
		Reason  string `json:"reason,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"conditions,omitempty"`
}

type GetAutoScalingStatusOutput struct {
	Header *ResponseHeader   `json:"header"`
	Status AutoScalingStatus `json:"autoScaling"`
}

type NetworkPolicyRule struct {
	Direction string   `json:"direction"` // Ingress, Egress
	Protocol  string   `json:"protocol"`  // TCP, UDP, ICMP
	Port      int      `json:"port,omitempty"`
	FromCIDR  []string `json:"fromCidr,omitempty"`
	ToCIDR    []string `json:"toCidr,omitempty"`
	Action    string   `json:"action"` // Allow, Deny
}

type NetworkPolicy struct {
	PolicyID    string              `json:"policyId"`
	PolicyName  string              `json:"policyName"`
	Description string              `json:"description,omitempty"`
	Rules       []NetworkPolicyRule `json:"rules"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

type ListNetworkPoliciesOutput struct {
	Header   *ResponseHeader `json:"header"`
	Policies []NetworkPolicy `json:"policies"`
}

type GetNetworkPolicyOutput struct {
	Header *ResponseHeader `json:"header"`
	Policy NetworkPolicy   `json:"policy"`
}

type CreateNetworkPolicyInput struct {
	PolicyName  string              `json:"policyName"`
	Description string              `json:"description,omitempty"`
	Rules       []NetworkPolicyRule `json:"rules"`
}

type UpdateNetworkPolicyInput struct {
	PolicyName  string              `json:"policyName,omitempty"`
	Description string              `json:"description,omitempty"`
	Rules       []NetworkPolicyRule `json:"rules,omitempty"`
}
