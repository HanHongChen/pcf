package sbi

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	local_models "github.com/free5gc/pcf/internal/models"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

// HTTPUpNodeSubscriptionsPost - Create UP Node Subscription
// POST /npcf-smpolicycontrol/v1/up-node-subscriptions
func (s *Server) HTTPUpNodeSubscriptionsPost(c *gin.Context) {
	var upNodeSubscription local_models.UpNodeSubscription

	// 解析請求體
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := util.GetProblemDetail("Bad request body", util.ERROR_REQUEST_PARAMETERS)
		logger.SmPolicyLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	// 反序列化JSON
	err = json.Unmarshal(requestBody, &upNodeSubscription)
	if err != nil {
		problemDetail := util.GetProblemDetail("Malformed request syntax", util.ERROR_REQUEST_PARAMETERS)
		logger.SmPolicyLog.Errorf("Unmarshal error: %+v", err)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	// 調用處理器
	s.Processor().HandleUpNodeSubscriptionCreate(c, upNodeSubscription)
}

// HTTPUpNodeSubscriptionsGet - Get UP Node Subscription
// GET /npcf-smpolicycontrol/v1/up-node-subscriptions/{subscriptionId}
func (s *Server) HTTPUpNodeSubscriptionsGet(c *gin.Context) {
	subscriptionId := c.Param("subscriptionId")
	
	if subscriptionId == "" {
		problemDetail := util.GetProblemDetail("Missing subscription ID", util.ERROR_REQUEST_PARAMETERS)
		logger.SmPolicyLog.Errorf("Missing subscription ID in request")
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	s.Processor().HandleUpNodeSubscriptionGet(c, subscriptionId)
}

// HTTPUpNodeSubscriptionsDelete - Delete UP Node Subscription  
// DELETE /npcf-smpolicycontrol/v1/up-node-subscriptions/{subscriptionId}
func (s *Server) HTTPUpNodeSubscriptionsDelete(c *gin.Context) {
	subscriptionId := c.Param("subscriptionId")
	
	if subscriptionId == "" {
		problemDetail := util.GetProblemDetail("Missing subscription ID", util.ERROR_REQUEST_PARAMETERS)
		logger.SmPolicyLog.Errorf("Missing subscription ID in request")
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	s.Processor().HandleUpNodeSubscriptionDelete(c, subscriptionId)
}

// AddUpNodeSubscriptionRoutes 添加UP節點訂閱相關路由
func (s *Server) AddUpNodeSubscriptionRoutes(group *gin.RouterGroup) {
	upNodeGroup := group.Group("/up-node-subscriptions")
	{
		upNodeGroup.POST("", s.HTTPUpNodeSubscriptionsPost)
		upNodeGroup.GET("/:subscriptionId", s.HTTPUpNodeSubscriptionsGet)
		upNodeGroup.DELETE("/:subscriptionId", s.HTTPUpNodeSubscriptionsDelete)
	}
}