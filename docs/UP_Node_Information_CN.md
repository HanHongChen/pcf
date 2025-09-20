# SMF向PCF發送PDU Session信息及訂閱機制

## 問題回答

### SMF是否會將UP node ID, port, UE addr, prefix, MTU傳給PCF？

根據對PCF代碼庫的分析，**目前的實現中SMF並不會直接將UP node ID和port信息傳給PCF**。但是，SMF會通過以下方式向PCF發送PDU Session相關的網路信息：

#### SMF向PCF發送的信息包括：

1. **UE地址信息**：
   - `Ipv4Address`: UE的IPv4地址
   - `Ipv6AddressPrefix`: UE的IPv6地址前綴
   - `IpDomain`: IP域信息

2. **PDU Session基本信息**：
   - `PduSessionId`: PDU Session ID
   - `PduSessionType`: PDU Session類型 (IPv4/IPv6/IPv4v6)
   - `Dnn`: 數據網路名稱
   - `SliceInfo`: 網路切片信息 (S-NSSAI)

3. **UE標識信息**：
   - `Supi`: 用戶永久標識符
   - `Gpsi`: 通用公共標識符
   - `Pei`: 永久設備標識符

4. **SMF信息**：
   - `SmfId`: SMF的標識符
   - `NotificationUri`: SMF的通知URI

#### 代碼中的相關實現：

```go
// 在 internal/sbi/processor/smpolicy.go 中
// SMF通過SmPolicyContextData結構發送信息給PCF
func (p *Processor) HandleCreateSmPolicyRequest(
    c *gin.Context,
    request models.SmPolicyContextData, // 包含上述所有信息
) {
    // PCF處理來自SMF的請求
}

// PCF創建綁定信息存儲這些網路參數
pcfBinding := models.PcfBinding{
    Supi:           request.Supi,
    Gpsi:           request.Gpsi,
    Ipv4Addr:       request.Ipv4Address,      // UE IPv4地址
    Ipv6Prefix:     request.Ipv6AddressPrefix, // UE IPv6前綴
    IpDomain:       request.IpDomain,         // IP域
    Dnn:            request.Dnn,              // 數據網路名稱
    Snssai:         request.SliceInfo,        // 切片信息
    PcfFqdn:        policyAuthorizationService.ApiPrefix,
    PcfIpEndPoints: *policyAuthorizationService.IpEndPoints,
}
```

### 關於UP node ID和port信息

**目前缺失的信息**：
- UP node ID (UPF標識符)
- UP node port (UPF端口)
- MTU信息

這些信息在當前的PCF實現中**並未被接收或處理**。如果需要這些信息，需要擴展`SmPolicyContextData`結構或創建新的通信機制。

## NF訂閱這些信息的方法

### 現有的訂閱機制

PCF提供了以下幾種訂閱和通知機制：

#### 1. SM Policy通知機制
```go
// SMF可以通過NotificationUri接收策略更新
type SmPolicyNotification struct {
    ResourceUri      string
    SmPolicyDecision *SmPolicyDecision
}

// PCF向SMF發送通知
go p.SendSMPolicyUpdateNotification(smPolicy.PolicyContext.NotificationUri, &notification)
```

#### 2. 影響數據訂閱機制  
```go
// 在 internal/sbi/consumer/udr_service.go 中
func (s *nudrService) CreateInfluenceDataSubscription(
    ue *pcf_context.UeContext, 
    request models.SmPolicyContextData
) (subscriptionID string, problemDetails *models.ProblemDetails, err error)
```

#### 3. 策略數據變更通知
```go
// 在 internal/sbi/processor/notifier.go 中
func (p *Processor) HandlePolicyDataChangeNotify(
    c *gin.Context,
    supi string,
    policyDataChangeNotification models.PolicyDataChangeNotification,
)
```

### 為NF實現訂閱UP node信息的建議方案

如果需要讓NF訂閱UP node信息，可以採用以下方法：

#### 方案1：擴展現有的SM Policy通知機制
```go
// 擴展SmPolicyDecision結構包含UP node信息
type SmPolicyDecision struct {
    // 現有字段...
    UpNodeInfo *UpNodeInformation `json:"upNodeInfo,omitempty"`
}

type UpNodeInformation struct {
    UpNodeId string `json:"upNodeId"`
    UpNodePort int32 `json:"upNodePort"`
    Mtu int32 `json:"mtu,omitempty"`
}
```

#### 方案2：創建專門的UP node信息訂閱API
```go
// 新增UP node信息訂閱端點
POST /npcf-smpolicycontrol/v1/up-node-subscriptions
{
    "supi": "imsi-123456789",
    "pduSessionId": 1,
    "notificationUri": "http://nf.example.com/notifications",
    "eventFilters": ["UP_NODE_CHANGE", "IP_ALLOCATION", "MTU_UPDATE"]
}
```

#### 方案3：利用現有的影響數據機制
```go
// 擴展TrafficInfluDataNotif包含UP node信息
type TrafficInfluDataNotif struct {
    // 現有字段...
    UpNodeInfo *UpNodeInformation `json:"upNodeInfo,omitempty"`
}
```

### 實現步驟

1. **擴展數據模型**：在`models.SmPolicyContextData`中添加UP node相關字段
2. **修改PCF處理邏輯**：在`HandleCreateSmPolicyRequest`中處理UP node信息
3. **實現訂閱API**：創建新的訂閱端點供NF使用
4. **添加通知機制**：當UP node信息變更時通知訂閱的NF
5. **更新存儲結構**：在`UeSmPolicyData`中存儲UP node信息

### 總結

- **目前**：SMF向PCF發送UE地址、前綴等網路信息，但不包括UP node ID和port
- **訂閱方式**：可以通過擴展現有的SM Policy通知機制或創建新的訂閱API來實現
- **建議**：采用方案2創建專門的UP node信息訂閱機制，提供更好的靈活性和擴展性