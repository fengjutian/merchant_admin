package controllers

import (
	model "merchant_back/internal/models"
	"merchant_back/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BusinessController 商家控制器
type BusinessController struct {
	businessService services.BusinessService
}

// NewBusinessController 创建商家控制器实例
func NewBusinessController(businessService services.BusinessService) *BusinessController {
	return &BusinessController{
		businessService: businessService,
	}
}

// GetBusinesses 获取商家列表
// @Summary 获取商家列表
// @Description 获取所有商家信息列表
// @Tags business
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/business [get]
func (bc *BusinessController) GetBusinesses(c *gin.Context) {
	businesses, err := bc.businessService.GetBusinesses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取商家列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取商家列表成功",
		"data":    businesses,
		"total":   len(businesses),
	})
}

// GetBusiness 获取单个商家
// @Summary 获取商家详情
// @Description 根据商家ID获取单个商家的详细信息
// @Tags business
// @Accept json
// @Produce json
// @Param id path int true "商家ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "无效的商家ID"
// @Failure 404 {object} map[string]interface{} "商家不存在"
// @Router /api/v1/business/{id} [get]
func (bc *BusinessController) GetBusiness(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商家ID",
		})
		return
	}

	business, err := bc.businessService.GetBusiness(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取商家成功",
		"data":    business,
	})
}

// CreateBusiness 创建商家
// @Summary 创建新商家
// @Description 创建一个新的商家账户
// @Tags business
// @Accept json
// @Produce json
// @Param business body model.Business true "商家信息"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "创建失败"
// @Router /api/v1/business [post]
func (bc *BusinessController) CreateBusiness(c *gin.Context) {
	var business model.Business
	if err := c.ShouldBindJSON(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	err := bc.businessService.CreateBusiness(&business)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建商家失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "商家创建成功",
		"data":    business,
	})
}

// UpdateBusiness 更新商家
// @Summary 更新商家信息
// @Description 根据商家ID更新商家的信息
// @Tags business
// @Accept json
// @Produce json
// @Param id path int true "商家ID"
// @Param business body model.Business true "商家更新信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "无效的商家ID或请求参数错误"
// @Failure 500 {object} map[string]interface{} "更新失败"
// @Router /api/v1/business/{id} [put]
func (bc *BusinessController) UpdateBusiness(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商家ID",
		})
		return
	}

	var business model.Business
	if err := c.ShouldBindJSON(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	err = bc.businessService.UpdateBusiness(id, &business)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新商家失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "商家更新成功",
		"data":    business,
	})
}

// UpdateBusinessStatus 更新商家状态
// @Summary 更新商家状态
// @Description 更新指定商家的状态（启用/禁用）
// @Tags business
// @Accept json
// @Produce json
// @Param id path int true "商家ID"
// @Param status body map[string]interface{} true "状态信息" example({"status": "active"})
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "无效的商家ID或请求参数错误"
// @Failure 500 {object} map[string]interface{} "更新失败"
// @Router /api/v1/business/{id}/status [put]
func (bc *BusinessController) UpdateBusinessStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商家ID",
		})
		return
	}

	var statusData map[string]interface{}
	if err := c.ShouldBindJSON(&statusData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	status, ok := statusData["status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少状态参数",
		})
		return
	}

	err = bc.businessService.UpdateBusinessStatus(id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新商家状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "商家状态更新成功",
		"data": map[string]interface{}{
			"id":     id,
			"status": status,
		},
	})
}

// DeleteBusiness 删除商家
// @Summary 删除商家
// @Description 根据商家ID删除商家账户
// @Tags business
// @Accept json
// @Produce json
// @Param id path int true "商家ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "无效的商家ID"
// @Failure 500 {object} map[string]interface{} "删除失败"
// @Router /api/v1/business/{id} [delete]
func (bc *BusinessController) DeleteBusiness(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的商家ID",
		})
		return
	}

	err = bc.businessService.DeleteBusiness(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除商家失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "商家删除成功",
	})
}

// GetBusinessByStatus 根据状态获取商家
// @Summary 根据状态获取商家列表
// @Description 根据商家状态获取商家列表
// @Tags business
// @Accept json
// @Produce json
// @Param status path string true "商家状态" Enums(active,inactive,suspended)
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "无效的状态参数"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/business/status/{status} [get]
func (bc *BusinessController) GetBusinessByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的状态参数",
		})
		return
	}

	businesses, err := bc.businessService.GetBusinessesByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取商家列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取商家列表成功",
		"data":    businesses,
		"total":   len(businesses),
		"status":  status,
	})
}

// SearchBusinesses 搜索商家
// @Summary 搜索商家
// @Description 根据关键词搜索商家
// @Tags business
// @Accept json
// @Produce json
// @Param q query string true "搜索关键词"
// @Success 200 {object} map[string]interface{} "搜索成功"
// @Failure 500 {object} map[string]interface{} "搜索失败"
// @Router /api/v1/business/search [get]
func (bc *BusinessController) SearchBusinesses(c *gin.Context) {
	query := c.Query("q")

	businesses, err := bc.businessService.SearchBusinesses(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索商家失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索完成",
		"data":    businesses,
		"total":   len(businesses),
		"query":   query,
	})
}

// GetBusinessesByType 根据类型获取商家
// @Summary 根据类型获取商家列表
// @Description 根据商家类型获取商家列表
// @Tags business
// @Accept json
// @Produce json
// @Param type query string true "商家类型"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "缺少类型参数"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/business/type [get]
func (bc *BusinessController) GetBusinessesByType(c *gin.Context) {
	businessType := c.Query("type")
	if businessType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少类型参数",
		})
		return
	}

	businesses, err := bc.businessService.GetBusinessesByType(businessType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取商家列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取商家列表成功",
		"data":    businesses,
		"total":   len(businesses),
		"type":    businessType,
	})
}

// GetBusinessesByRating 根据评分获取商家
// @Summary 根据评分获取商家列表
// @Description 获取评分高于指定值的商家列表
// @Tags business
// @Accept json
// @Produce json
// @Param minRating query number true "最低评分"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "缺少评分参数或参数无效"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/business/rating [get]
func (bc *BusinessController) GetBusinessesByRating(c *gin.Context) {
	minRatingStr := c.Query("minRating")
	if minRatingStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少最低评分参数",
		})
		return
	}

	minRating, err := strconv.ParseFloat(minRatingStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的评分参数",
		})
		return
	}

	businesses, err := bc.businessService.GetBusinessesByRating(minRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取商家列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "获取商家列表成功",
		"data":      businesses,
		"total":     len(businesses),
		"minRating": minRating,
	})
}

// GetBusinessesWithPagination 分页获取商家
// @Summary 分页获取商家列表
// @Description 分页获取商家列表，支持自定义页码和每页数量
// @Tags business
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/business/page [get]
func (bc *BusinessController) GetBusinessesWithPagination(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	businesses, total, err := bc.businessService.GetBusinessesWithPagination(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取商家列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"message":  "获取商家列表成功",
		"data":     businesses,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetNearbyBusinesses 获取附近商家
// @Summary 获取附近商家
// @Description 根据经纬度和半径获取附近的商家
// @Tags business
// @Accept json
// @Produce json
// @Param lat query number true "纬度"
// @Param lng query number true "经度"
// @Param radius query number false "搜索半径(公里)" default(5.0)
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "缺少经纬度参数或参数无效"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/business/nearby [get]
func (bc *BusinessController) GetNearbyBusinesses(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "5.0")

	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少经纬度参数",
		})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的纬度参数",
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的经度参数",
		})
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius <= 0 {
		radius = 5.0 // 默认5公里
	}

	businesses, err := bc.businessService.GetNearbyBusinesses(lat, lng, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取附近商家失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取附近商家成功",
		"data":    businesses,
		"total":   len(businesses),
		"lat":     lat,
		"lng":     lng,
		"radius":  radius,
	})
}

// GetBusinessCount 获取商家总数
// @Summary 获取商家总数
// @Description 获取所有商家的总数量
// @Tags business
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/business/count [get]
func (bc *BusinessController) GetBusinessCount(c *gin.Context) {
	count, err := bc.businessService.GetBusinessCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取商家总数失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取商家总数成功",
		"data":    count,
	})
}

// GetBusinessCountByType 根据类型获取商家数量
// @Summary 根据类型获取商家数量
// @Description 根据商家类型获取该类型的商家数量
// @Tags business
// @Accept json
// @Produce json
// @Param type query string true "商家类型"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "缺少类型参数"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/business/count/type [get]
func (bc *BusinessController) GetBusinessCountByType(c *gin.Context) {
	businessType := c.Query("type")
	if businessType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少类型参数",
		})
		return
	}

	count, err := bc.businessService.GetBusinessCountByType(businessType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取商家数量失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取商家数量成功",
		"data":    count,
		"type":    businessType,
	})
}
