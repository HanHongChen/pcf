# UP節點信息訂閱使用範例

## 概述

本文檔提供了如何使用PCF的UP節點信息訂閱功能的詳細範例，包括API調用方法和實際使用場景。

## API端點

### 創建UP節點訂閱

**請求**：
```http
POST /npcf-smpolicycontrol/v1/up-node-subscriptions
Content-Type: application/json

{
    "subscriberId": "nf-consumer-001",
    "supi": "imsi-123456789012345",
    "pduSessionId": 1,
    "notificationUri": "http://consumer-nf.example.com/notifications/up-node",
    "eventFilters": [
        "UP_NODE_CHANGE",
        "IP_ALLOCATION", 
        "MTU_UPDATE",
        "N3_CHANGE"
    ],
    "expiryTime": "2024-12-31T23:59:59Z"
}
```

**響應**：
```http
HTTP/1.1 201 Created
Location: /npcf-smpolicycontrol/v1/up-node-subscriptions/550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json

{
    "subscriptionId": "550e8400-e29b-41d4-a716-446655440000",
    "subscriptionUri": "/npcf-smpolicycontrol/v1/up-node-subscriptions/550e8400-e29b-41d4-a716-446655440000",
    "currentUpNodeInfo": {
        "upNodeId": "upf-1",
        "upNodePort": 8805,
        "mtu": 1500,
        "upfIpv4Address": "192.168.1.100",
        "n3InterfaceInfo": {
            "gtpuEndpoint": {
                "ipAddress": "192.168.1.100",
                "port": 8805,
                "teid": "teid-1"
            },
            "tunnelInfo": {
                "tunnelId": "tunnel-1",
                "tunnelType": "GTP"
            }
        },
        "n6InterfaceInfo": {
            "dnAccessInfo": {
                "gateway": "10.0.0.1",
                "dnsServers": ["8.8.8.8", "8.8.4.4"]
            }
        }
    }
}
```

### 查詢UP節點訂閱

**請求**：
```http
GET /npcf-smpolicycontrol/v1/up-node-subscriptions/550e8400-e29b-41d4-a716-446655440000
```

**響應**：
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "subscriptionId": "550e8400-e29b-41d4-a716-446655440000",
    "subscriptionUri": "/npcf-smpolicycontrol/v1/up-node-subscriptions/550e8400-e29b-41d4-a716-446655440000",
    "currentUpNodeInfo": {
        "upNodeId": "upf-1",
        "upNodePort": 8805,
        "mtu": 1500,
        "upfIpv4Address": "192.168.1.100",
        "n3InterfaceInfo": {
            "gtpuEndpoint": {
                "ipAddress": "192.168.1.100", 
                "port": 8805,
                "teid": "teid-1"
            }
        }
    }
}
```

### 刪除UP節點訂閱

**請求**：
```http
DELETE /npcf-smpolicycontrol/v1/up-node-subscriptions/550e8400-e29b-41d4-a716-446655440000
```

**響應**：
```http
HTTP/1.1 204 No Content
```

## 通知範例

當UP節點信息發生變更時，PCF會向訂閱者發送通知：

```http
POST http://consumer-nf.example.com/notifications/up-node
Content-Type: application/json

{
    "resourceUri": "/npcf-smpolicycontrol/v1/up-node-subscriptions/550e8400-e29b-41d4-a716-446655440000",
    "eventType": "UP_NODE_CHANGE",
    "upNodeInfo": {
        "upNodeId": "upf-2",
        "upNodePort": 8805,
        "mtu": 1400,
        "upfIpv4Address": "192.168.1.101",
        "n3InterfaceInfo": {
            "gtpuEndpoint": {
                "ipAddress": "192.168.1.101",
                "port": 8805,
                "teid": "teid-1"
            }
        }
    },
    "timestamp": "2024-01-15T10:30:00Z",
    "additionalInfo": {
        "reason": "UPF_RESELECTION",
        "previousUpNodeId": "upf-1"
    }
}
```

## 使用場景

### 1. 網路功能監控

```go
// NF監控服務訂閱UP節點變更
subscription := UpNodeSubscription{
    SubscriberId:    "monitoring-service",
    Supi:           "imsi-123456789012345",
    PduSessionId:   1,
    NotificationUri: "http://monitoring.example.com/up-node-events",
    EventFilters:   []UpNodeEventType{
        UpNodeEventType_UP_NODE_CHANGE,
        UpNodeEventType_N3_CHANGE,
        UpNodeEventType_MTU_UPDATE,
    },
}
```

### 2. 流量優化服務

```go
// 流量優化服務訂閱IP分配和MTU變更
subscription := UpNodeSubscription{
    SubscriberId:    "traffic-optimizer",
    Supi:           "imsi-123456789012345", 
    PduSessionId:   1,
    NotificationUri: "http://optimizer.example.com/ip-mtu-events",
    EventFilters:   []UpNodeEventType{
        UpNodeEventType_IP_ALLOCATION,
        UpNodeEventType_MTU_UPDATE,
    },
}
```

### 3. 故障檢測與恢復

```go
// 故障管理服務訂閱所有UP節點相關事件
subscription := UpNodeSubscription{
    SubscriberId:    "fault-management",
    Supi:           "imsi-123456789012345",
    PduSessionId:   1,
    NotificationUri: "http://fault-mgmt.example.com/up-node-events",
    // 不設置eventFilters以接收所有事件
}
```

## 事件類型說明

| 事件類型 | 描述 | 觸發條件 |
|---------|------|----------|
| `UP_NODE_CHANGE` | UP節點變更 | UPF重選、故障切換 |
| `IP_ALLOCATION` | IP地址分配 | UE獲得新IP地址 |
| `MTU_UPDATE` | MTU更新 | 網路MTU配置變更 |
| `N3_CHANGE` | N3介面變更 | GTP隧道參數變更 |
| `N6_CHANGE` | N6介面變更 | 數據網路接入變更 |
| `TUNNEL_CHANGE` | 隧道變更 | GTP隧道重建或修改 |

## 錯誤處理

### 常見錯誤碼

- `400 Bad Request`: 請求參數錯誤
- `404 Not Found`: UE或PDU Session不存在
- `409 Conflict`: 訂閱已存在
- `500 Internal Server Error`: 內部處理錯誤

### 錯誤響應範例

```json
{
    "type": "https://example.com/pcf/errors/user-unknown",
    "title": "User Unknown",
    "status": 404,
    "detail": "UE not found in PCF",
    "instance": "/npcf-smpolicycontrol/v1/up-node-subscriptions",
    "cause": "USER_UNKNOWN"
}
```

## 實現注意事項

1. **訂閱生命週期管理**：需要實現訂閱到期自動清理機制
2. **通知可靠性**：實現重試機制確保通知送達
3. **事件去重**：避免重複事件通知
4. **安全性**：驗證訂閱者身份和權限
5. **性能優化**：對於大量訂閱的場景需要優化查找和通知性能

## 集成指南

### 客戶端實現

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func createUpNodeSubscription(pcfUrl string, subscription UpNodeSubscription) (*UpNodeSubscriptionResponse, error) {
    jsonData, err := json.Marshal(subscription)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(
        pcfUrl+"/npcf-smpolicycontrol/v1/up-node-subscriptions",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("subscription failed with status: %d", resp.StatusCode)
    }

    var response UpNodeSubscriptionResponse
    err = json.NewDecoder(resp.Body).Decode(&response)
    return &response, err
}
```

### 通知處理

```go
func handleUpNodeNotification(w http.ResponseWriter, r *http.Request) {
    var notification UpNodeNotification
    err := json.NewDecoder(r.Body).Decode(&notification)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // 處理UP節點變更通知
    switch notification.EventType {
    case UpNodeEventType_UP_NODE_CHANGE:
        // 處理UP節點變更
        handleUpNodeChange(notification.UpNodeInfo)
    case UpNodeEventType_IP_ALLOCATION:
        // 處理IP分配
        handleIpAllocation(notification.UpNodeInfo)
    case UpNodeEventType_MTU_UPDATE:
        // 處理MTU更新
        handleMtuUpdate(notification.UpNodeInfo)
    }

    w.WriteHeader(http.StatusNoContent)
}
```