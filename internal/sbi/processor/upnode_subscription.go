package processor

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
	pcf_context "github.com/free5gc/pcf/internal/context"
	local_models "github.com/free5gc/pcf/internal/models"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

// HandleUpNodeSubscriptionCreate 處理UP節點信息訂閱創建請求
// POST /npcf-smpolicycontrol/v1/up-node-subscriptions
func (p *Processor) HandleUpNodeSubscriptionCreate(
	c *gin.Context,
	request local_models.UpNodeSubscription,
) {
	logger.SmPolicyLog.Infof("Handle UP Node Subscription Create for SUPI: %s, PDU Session: %d", 
		request.Supi, request.PduSessionId)

	// 驗證必要參數
	if request.Supi == "" || request.NotificationUri == "" {
		problemDetail := util.GetProblemDetail("Missing mandatory parameters", util.ERROR_INITIAL_PARAMETERS)
		logger.SmPolicyLog.Warnf("Missing mandatory parameters in UP node subscription request")
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	// 檢查UE是否存在
	pcfSelf := p.Context()
	var ue *pcf_context.UeContext
	if val, exist := pcfSelf.UePool.Load(request.Supi); exist {
		ue = val.(*pcf_context.UeContext)
	}

	if ue == nil {
		problemDetail := util.GetProblemDetail("UE not found in PCF", util.USER_UNKNOWN)
		logger.SmPolicyLog.Warnf("UE[%s] not found in PCF", request.Supi)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	// 檢查PDU Session是否存在
	smPolicyID := fmt.Sprintf("%s-%d", request.Supi, request.PduSessionId)
	smPolicy := ue.SmPolicyData[smPolicyID]
	if smPolicy == nil {
		problemDetail := util.GetProblemDetail("PDU Session not found", util.CONTEXT_NOT_FOUND)
		logger.SmPolicyLog.Warnf("PDU Session[%d] not found for UE[%s]", request.PduSessionId, request.Supi)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	// 生成訂閱ID
	subscriptionId := uuid.New().String()
	
	// 存儲訂閱信息
	if smPolicy.UpNodeSubscriptions == nil {
		smPolicy.UpNodeSubscriptions = make(map[string]interface{})
	}
	smPolicy.UpNodeSubscriptions[subscriptionId] = &request

	// 獲取當前UP節點信息（模擬從SMF或UPF獲取）
	currentUpNodeInfo := p.getCurrentUpNodeInfo(smPolicy)

	// 創建響應
	subscriptionUri := fmt.Sprintf("/npcf-smpolicycontrol/v1/up-node-subscriptions/%s", subscriptionId)
	response := local_models.UpNodeSubscriptionResponse{
		SubscriptionId:    subscriptionId,
		SubscriptionUri:   subscriptionUri,
		CurrentUpNodeInfo: currentUpNodeInfo,
	}

	logger.SmPolicyLog.Infof("UP Node subscription created with ID: %s", subscriptionId)
	
	// 設置Location header
	c.Header("Location", util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, subscriptionId))
	c.JSON(http.StatusCreated, response)
}

// HandleUpNodeSubscriptionGet 獲取UP節點訂閱信息
// GET /npcf-smpolicycontrol/v1/up-node-subscriptions/{subscriptionId}
func (p *Processor) HandleUpNodeSubscriptionGet(
	c *gin.Context,
	subscriptionId string,
) {
	logger.SmPolicyLog.Infof("Handle UP Node Subscription Get for ID: %s", subscriptionId)

	// 查找訂閱信息
	subscription, smPolicy := p.findUpNodeSubscription(subscriptionId)
	if subscription == nil {
		problemDetail := util.GetProblemDetail("Subscription not found", util.CONTEXT_NOT_FOUND)
		logger.SmPolicyLog.Warnf("UP Node subscription[%s] not found", subscriptionId)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	// 獲取當前UP節點信息
	currentUpNodeInfo := p.getCurrentUpNodeInfo(smPolicy)

	response := local_models.UpNodeSubscriptionResponse{
		SubscriptionId:    subscriptionId,
		SubscriptionUri:   fmt.Sprintf("/npcf-smpolicycontrol/v1/up-node-subscriptions/%s", subscriptionId),
		CurrentUpNodeInfo: currentUpNodeInfo,
	}

	c.JSON(http.StatusOK, response)
}

// HandleUpNodeSubscriptionDelete 刪除UP節點訂閱
// DELETE /npcf-smpolicycontrol/v1/up-node-subscriptions/{subscriptionId}
func (p *Processor) HandleUpNodeSubscriptionDelete(
	c *gin.Context,
	subscriptionId string,
) {
	logger.SmPolicyLog.Infof("Handle UP Node Subscription Delete for ID: %s", subscriptionId)

	// 查找並刪除訂閱信息
	subscription, smPolicy := p.findUpNodeSubscription(subscriptionId)
	if subscription == nil {
		problemDetail := util.GetProblemDetail("Subscription not found", util.CONTEXT_NOT_FOUND)
		logger.SmPolicyLog.Warnf("UP Node subscription[%s] not found", subscriptionId)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	// 刪除訂閱
	delete(smPolicy.UpNodeSubscriptions, subscriptionId)
	
	logger.SmPolicyLog.Infof("UP Node subscription[%s] deleted successfully", subscriptionId)
	c.JSON(http.StatusNoContent, nil)
}

// SendUpNodeNotification 發送UP節點變更通知給訂閱者
func (p *Processor) SendUpNodeNotification(
	smPolicy *pcf_context.UeSmPolicyData,
	eventType local_models.UpNodeEventType,
	upNodeInfo *local_models.UpNodeInformation,
) {
	if smPolicy.UpNodeSubscriptions == nil {
		return
	}

	for subscriptionId, subscriptionInterface := range smPolicy.UpNodeSubscriptions {
		subscription, ok := subscriptionInterface.(*local_models.UpNodeSubscription)
		if !ok {
			continue
		}
		
		// 檢查事件過濾器
		if !p.shouldNotifySubscriber(subscription, eventType) {
			continue
		}

		// 創建通知
		notification := local_models.UpNodeNotification{
			ResourceUri: fmt.Sprintf("/npcf-smpolicycontrol/v1/up-node-subscriptions/%s", subscriptionId),
			EventType:   eventType,
			UpNodeInfo:  upNodeInfo,
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		// 異步發送通知
		go p.sendUpNodeNotificationToSubscriber(subscription.NotificationUri, &notification)
	}
}

// getCurrentUpNodeInfo 獲取當前UP節點信息（模擬實現）
func (p *Processor) getCurrentUpNodeInfo(smPolicy *pcf_context.UeSmPolicyData) *local_models.UpNodeInformation {
	// 這裡應該從實際的UPF或SMF獲取信息
	// 當前提供模擬數據作為示例
	
	policyContext := smPolicy.PolicyContext
	if policyContext == nil {
		return nil
	}

	upNodeInfo := &local_models.UpNodeInformation{
		UpNodeId:       "upf-" + strconv.Itoa(int(policyContext.PduSessionId)), // 模擬UP node ID
		UpNodePort:     8805, // 標準GTP-U端口
		Mtu:            1500, // 默認MTU
		UpfIpv4Address: "192.168.1.100", // 模擬UPF IP
	}

	// 添加N3介面信息
	if policyContext.Ipv4Address != "" || policyContext.Ipv6AddressPrefix != "" {
		upNodeInfo.N3InterfaceInfo = &local_models.N3InterfaceInformation{
			GtpuEndpoint: &local_models.GtpuEndpoint{
				IpAddress: "192.168.1.100",
				Port:      8805,
				Teid:      fmt.Sprintf("teid-%d", policyContext.PduSessionId),
			},
			TunnelInfo: &local_models.TunnelInfo{
				TunnelId:   fmt.Sprintf("tunnel-%d", policyContext.PduSessionId),
				TunnelType: "GTP",
			},
		}
	}

	// 添加N6介面信息
	upNodeInfo.N6InterfaceInfo = &local_models.N6InterfaceInformation{
		DnAccessInfo: &local_models.DnAccessInfo{
			Gateway:    "10.0.0.1",
			DnsServers: []string{"8.8.8.8", "8.8.4.4"},
		},
	}

	return upNodeInfo
}

// findUpNodeSubscription 查找UP節點訂閱信息
func (p *Processor) findUpNodeSubscription(subscriptionId string) (*local_models.UpNodeSubscription, *pcf_context.UeSmPolicyData) {
	pcfSelf := p.Context()
	
	// 遍歷所有UE尋找訂閱
	var foundSubscription *local_models.UpNodeSubscription
	var foundSmPolicy *pcf_context.UeSmPolicyData
	
	pcfSelf.UePool.Range(func(key, value interface{}) bool {
		ue := value.(*pcf_context.UeContext)
		for _, smPolicy := range ue.SmPolicyData {
			if smPolicy.UpNodeSubscriptions != nil {
				if subscriptionInterface, exists := smPolicy.UpNodeSubscriptions[subscriptionId]; exists {
					if subscription, ok := subscriptionInterface.(*local_models.UpNodeSubscription); ok {
						foundSubscription = subscription
						foundSmPolicy = smPolicy
						return false // 找到了，停止遍歷
					}
				}
			}
		}
		return true // 繼續遍歷
	})
	
	return foundSubscription, foundSmPolicy
}

// shouldNotifySubscriber 檢查是否應該通知訂閱者
func (p *Processor) shouldNotifySubscriber(subscription *local_models.UpNodeSubscription, eventType local_models.UpNodeEventType) bool {
	// 如果沒有設置事件過濾器，則通知所有事件
	if len(subscription.EventFilters) == 0 {
		return true
	}

	// 檢查事件類型是否在過濾器中
	for _, filter := range subscription.EventFilters {
		if filter == eventType {
			return true
		}
	}

	return false
}

// sendUpNodeNotificationToSubscriber 向訂閱者發送UP節點通知
func (p *Processor) sendUpNodeNotificationToSubscriber(notificationUri string, notification *local_models.UpNodeNotification) {
	// 這裡應該實現實際的HTTP POST請求發送通知
	// 當前僅記錄日誌作為示例
	logger.SmPolicyLog.Infof("Sending UP Node notification to %s: EventType=%s, UpNodeId=%s", 
		notificationUri, notification.EventType, notification.UpNodeInfo.UpNodeId)
	
	// 實際實現應該類似於：
	// client := &http.Client{}
	// jsonData, _ := json.Marshal(notification)
	// resp, err := client.Post(notificationUri, "application/json", bytes.NewBuffer(jsonData))
	// ... 處理響應和錯誤
}