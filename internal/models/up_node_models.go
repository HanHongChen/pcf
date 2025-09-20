package models

// UpNodeInformation 定義UP節點信息結構
type UpNodeInformation struct {
	// UP node ID (UPF標識符)
	UpNodeId string `json:"upNodeId" yaml:"upNodeId" bson:"upNodeId"`
	
	// UP node端口信息
	UpNodePort int32 `json:"upNodePort,omitempty" yaml:"upNodePort,omitempty" bson:"upNodePort,omitempty"`
	
	// MTU信息
	Mtu int32 `json:"mtu,omitempty" yaml:"mtu,omitempty" bson:"mtu,omitempty"`
	
	// UPF IP地址
	UpfIpv4Address string `json:"upfIpv4Address,omitempty" yaml:"upfIpv4Address,omitempty" bson:"upfIpv4Address,omitempty"`
	UpfIpv6Address string `json:"upfIpv6Address,omitempty" yaml:"upfIpv6Address,omitempty" bson:"upfIpv6Address,omitempty"`
	
	// N3介面資訊 (AN-UPF)
	N3InterfaceInfo *N3InterfaceInformation `json:"n3InterfaceInfo,omitempty" yaml:"n3InterfaceInfo,omitempty" bson:"n3InterfaceInfo,omitempty"`
	
	// N6介面資訊 (UPF-DN)
	N6InterfaceInfo *N6InterfaceInformation `json:"n6InterfaceInfo,omitempty" yaml:"n6InterfaceInfo,omitempty" bson:"n6InterfaceInfo,omitempty"`
}

// N3InterfaceInformation N3介面詳細資訊 (AN-UPF介面)
type N3InterfaceInformation struct {
	// GTP-U端點資訊
	GtpuEndpoint *GtpuEndpoint `json:"gtpuEndpoint,omitempty" yaml:"gtpuEndpoint,omitempty" bson:"gtpuEndpoint,omitempty"`
	
	// 隧道資訊
	TunnelInfo *TunnelInfo `json:"tunnelInfo,omitempty" yaml:"tunnelInfo,omitempty" bson:"tunnelInfo,omitempty"`
}

// N6InterfaceInformation N6介面詳細資訊 (UPF-DN介面)
type N6InterfaceInformation struct {
	// 數據網路接入資訊
	DnAccessInfo *DnAccessInfo `json:"dnAccessInfo,omitempty" yaml:"dnAccessInfo,omitempty" bson:"dnAccessInfo,omitempty"`
}

// GtpuEndpoint GTP-U端點資訊
type GtpuEndpoint struct {
	IpAddress string `json:"ipAddress" yaml:"ipAddress" bson:"ipAddress"`
	Port      int32  `json:"port" yaml:"port" bson:"port"`
	Teid      string `json:"teid,omitempty" yaml:"teid,omitempty" bson:"teid,omitempty"`
}

// TunnelInfo 隧道資訊
type TunnelInfo struct {
	TunnelId   string `json:"tunnelId,omitempty" yaml:"tunnelId,omitempty" bson:"tunnelId,omitempty"`
	TunnelType string `json:"tunnelType,omitempty" yaml:"tunnelType,omitempty" bson:"tunnelType,omitempty"`
}

// DnAccessInfo 數據網路接入資訊
type DnAccessInfo struct {
	Gateway    string `json:"gateway,omitempty" yaml:"gateway,omitempty" bson:"gateway,omitempty"`
	DnsServers []string `json:"dnsServers,omitempty" yaml:"dnsServers,omitempty" bson:"dnsServers,omitempty"`
}

// UpNodeSubscription UP節點信息訂閱請求
type UpNodeSubscription struct {
	// 訂閱者資訊
	SubscriberId string `json:"subscriberId" yaml:"subscriberId" bson:"subscriberId"`
	
	// 目標UE
	Supi string `json:"supi" yaml:"supi" bson:"supi"`
	
	// PDU Session ID
	PduSessionId int32 `json:"pduSessionId" yaml:"pduSessionId" bson:"pduSessionId"`
	
	// 通知URI
	NotificationUri string `json:"notificationUri" yaml:"notificationUri" bson:"notificationUri"`
	
	// 事件過濾器
	EventFilters []UpNodeEventType `json:"eventFilters,omitempty" yaml:"eventFilters,omitempty" bson:"eventFilters,omitempty"`
	
	// 訂閱有效期
	ExpiryTime *string `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty" bson:"expiryTime,omitempty"`
}

// UpNodeEventType UP節點事件類型
type UpNodeEventType string

const (
	UpNodeEventType_UP_NODE_CHANGE  UpNodeEventType = "UP_NODE_CHANGE"
	UpNodeEventType_IP_ALLOCATION   UpNodeEventType = "IP_ALLOCATION"
	UpNodeEventType_MTU_UPDATE      UpNodeEventType = "MTU_UPDATE"
	UpNodeEventType_N3_CHANGE       UpNodeEventType = "N3_CHANGE"
	UpNodeEventType_N6_CHANGE       UpNodeEventType = "N6_CHANGE"
	UpNodeEventType_TUNNEL_CHANGE   UpNodeEventType = "TUNNEL_CHANGE"
)

// UpNodeNotification UP節點變更通知
type UpNodeNotification struct {
	// 資源URI
	ResourceUri string `json:"resourceUri" yaml:"resourceUri" bson:"resourceUri"`
	
	// 事件類型
	EventType UpNodeEventType `json:"eventType" yaml:"eventType" bson:"eventType"`
	
	// UP節點信息
	UpNodeInfo *UpNodeInformation `json:"upNodeInfo,omitempty" yaml:"upNodeInfo,omitempty" bson:"upNodeInfo,omitempty"`
	
	// 時間戳
	Timestamp string `json:"timestamp" yaml:"timestamp" bson:"timestamp"`
	
	// 額外資訊
	AdditionalInfo map[string]interface{} `json:"additionalInfo,omitempty" yaml:"additionalInfo,omitempty" bson:"additionalInfo,omitempty"`
}

// UpNodeSubscriptionResponse UP節點訂閱響應
type UpNodeSubscriptionResponse struct {
	// 訂閱ID
	SubscriptionId string `json:"subscriptionId" yaml:"subscriptionId" bson:"subscriptionId"`
	
	// 訂閱URI
	SubscriptionUri string `json:"subscriptionUri" yaml:"subscriptionUri" bson:"subscriptionUri"`
	
	// 當前UP節點信息
	CurrentUpNodeInfo *UpNodeInformation `json:"currentUpNodeInfo,omitempty" yaml:"currentUpNodeInfo,omitempty" bson:"currentUpNodeInfo,omitempty"`
}