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

// DeleteBusiness 删除商家
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

// SearchBusinesses 搜索商家
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
