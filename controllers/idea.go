package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/constant"
	"github.com/anhhuy1010/DATN-cms-ideas/helpers/respond"
	"github.com/anhhuy1010/DATN-cms-ideas/helpers/util"
	"github.com/anhhuy1010/DATN-cms-ideas/models"
	request "github.com/anhhuy1010/DATN-cms-ideas/request/ideas"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IdeasController struct {
}

func (ideaCtl IdeasController) List(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var req request.GetListRequest
	err := c.ShouldBindWith(&req, binding.Query)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	cond := bson.M{"is_active": 1} // Chỉ lấy những ý tưởng đã được duyệt
	cond["is_delete"] = 0          // Chỉ lấy những ý tưởng chưa bị xóa
	cond["status"] = 0             // Chỉ lấy những ý tưởng chưa được mua

	if req.Industry != nil && *req.Industry != "" {
		cond["industry"] = req.Industry
	}
	if req.PriceTier != nil && *req.PriceTier != "" {
		var priceCond bson.M

		switch *req.PriceTier {
		case "tier1":
			// Dưới 1 triệu
			priceCond = bson.M{"$lte": 1000000}
		case "tier2":
			// Từ 1 triệu đến dưới 3 triệu
			priceCond = bson.M{
				"$gte": 1000000,
				"$lte": 3000000,
			}
		case "tier3":
			// Từ 3 triệu đến dưới 5 triệu
			priceCond = bson.M{
				"$gte": 3000000,
				"$lte": 5000000,
			}
		case "tier4":
			// Trên 5 triệu
			priceCond = bson.M{"$gte": 5000000}
		}

		if priceCond != nil {
			cond["price"] = priceCond
		}
	}
	if req.Ideasname != nil {
		cond["ideasname"] = bson.M{
			"$regex":   *req.Ideasname,
			"$options": "i", // Tìm kiếm không phân biệt chữ hoa thường
		}
	}

	var optionsQuery models.ModelOption
	var page, limit int

	if req.TopViewOnly {
		optionsQuery, page, limit = models.GetPagingOption(1, 5, "-view")
		optionsQuery.SortBy = "view"
		optionsQuery.SortDir = -1
		page = 1
		limit = 5
	} else {
		optionsQuery, page, limit = models.GetPagingOption(req.Page, 3, req.Sort)
		if req.Sort == "" {
			optionsQuery.SortBy = "post_day"
			optionsQuery.SortDir = -1
		}
	}
	var respData []request.ListResponse
	ideas, err := ideaModel.Pagination(c, cond, optionsQuery)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	respData = make([]request.ListResponse, 0, len(ideas))
	for _, idea := range ideas {
		image := ""
		if len(idea.Image) > 0 {
			image = idea.Image[0] // lấy ảnh đầu tiên
		}

		res := request.ListResponse{
			Uuid:          idea.Uuid,
			CustomerName:  idea.CustomerName,
			IdeasName:     idea.IdeasName,
			Industry:      idea.Industry,
			Image:         image,
			ContentDetail: idea.ContentDetail,
			PostDay:       idea.PostDay,
			Price:         idea.Price,
			View:          idea.View,
			IsActive:      idea.IsActive,
			IsDelete:      idea.IsDelete,
		}
		respData = append(respData, res)
	}
	total, err := ideaModel.Count(c, cond)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	pages := int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}

// ///////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) Detail(c *gin.Context) {

	ideaModel := new(models.Ideas)
	var reqUri request.GetDetailUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	idea, err := ideaModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("uuid no found!"))
		return
	}
	update := bson.M{
		"$inc": bson.M{"view": 1},
	}
	_, _ = ideaModel.UpdateByCondition(condition, update)

	response := request.GetDetailResponse{
		IdeasName:       idea.IdeasName,
		Industry:        idea.Industry,
		Procedure:       idea.Procedure,
		ContentDetail:   idea.ContentDetail,
		Value_Benefits:  idea.Value_Benefits,
		Is_Intellect:    idea.Is_Intellect,
		Price:           idea.Price,
		CustomerName:    idea.CustomerName,
		CustomerEmail:   idea.CustomerEmail,
		Image:           idea.Image,
		View:            idea.View + 1,
		CustomerUuid:    idea.CustomerUuid,
		Image_Intellect: idea.Image_Intellect,
		PostDay:         idea.PostDay,
		IsActive:        idea.IsActive,
		IsDelete:        idea.IsDelete,
		Status:          idea.Status,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

// ///////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) SearchIdeasByTitle(c *gin.Context) {
	title := c.Query("title")
	// Tìm kiếm không phân biệt chữ hoa thường
	filter := bson.M{
		"ideasname": bson.M{
			"$regex":   title,
			"$options": "i",
		},
	}

	ideaModel := models.Ideas{}
	results, err := ideaModel.Find(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search ideas", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// ////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) GetMyIdeas(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var req request.GetListRequest
	err := c.ShouldBindWith(&req, binding.Query)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	customerUuid, _ := c.Get("customer_uuid")
	if customerUuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customeruuid is required"})
		return
	}

	cond := bson.M{
		"customeruuid": customerUuid, // lọc theo người dùng hiện tại
	}

	optionsQuery, page, limit := models.GetPagingOption(req.Page, req.Limit, req.Sort)
	if req.Sort == "" {
		optionsQuery.SortBy = "post_day"
		optionsQuery.SortDir = -1
	}
	respData := make([]request.ListResponse, 0)

	ideas, err := ideaModel.Pagination(c, cond, optionsQuery)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ideas"})
		return
	}

	for _, idea := range ideas {
		image := ""
		if len(idea.Image) > 0 {
			image = idea.Image[0] // lấy ảnh đầu tiên
		}

		res := request.ListResponse{
			Uuid:          idea.Uuid,
			CustomerName:  idea.CustomerName,
			IdeasName:     idea.IdeasName,
			Industry:      idea.Industry,
			Image:         image,
			ContentDetail: idea.ContentDetail,
			PostDay:       idea.PostDay,
			Price:         idea.Price,
			View:          idea.View,
			IsActive:      idea.IsActive,
			IsDelete:      idea.IsDelete,
		}
		respData = append(respData, res)
	}

	total, err := ideaModel.Count(c, cond)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count ideas"})
		return
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}

// ////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) GetMyIdeaDetail(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var reqUri request.GetDetailUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	customerUuid, exists := c.Get("customer_uuid")
	if !exists || customerUuid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Tìm ý tưởng theo uuid và customer_uuid (đảm bảo là của người dùng hiện tại)
	cond := bson.M{
		"customeruuid": customerUuid,
		"uuid":         reqUri.Uuid,
	}
	idea, err := ideaModel.FindOne(cond)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Idea not found"})
		return
	}

	// Trả về dữ liệu chi tiết
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"uuid":            idea.Uuid,
			"ideasname":       idea.IdeasName,
			"industry":        idea.Industry,
			"content_detail":  idea.ContentDetail,
			"value_benefits":  idea.Value_Benefits,
			"price":           idea.Price,
			"is_intellect":    idea.Is_Intellect,
			"is_active":       idea.IsActive,
			"post_day":        idea.PostDay,
			"image":           idea.Image,
			"view":            idea.View,
			"customer_name":   idea.CustomerName,
			"customer_email":  idea.CustomerEmail,
			"image_intellect": idea.Image_Intellect,
		},
	})
}

// ///////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) PostIdea(c *gin.Context) {
	customerUUID, _ := c.Get("customer_uuid")
	customerName, _ := c.Get("customer_name")
	customerEmail, _ := c.Get("customer_email")

	// Lấy dữ liệu json ý tưởng từ form field "data"
	var input struct {
		IdeasName       string   `json:"ideasname" binding:"required"`
		Industry        string   `json:"industry" binding:"required"`
		ContentDetail   string   `json:"content_detail"`
		Value_Benefits  string   `json:"value_benefits"`
		Price           int      `json:"price"`
		Is_Intellect    int      `json:"is_intellect"`
		Image           []string `json:"image"`
		Image_Intellect string   `json:"image_intellect"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idea := models.Ideas{
		Uuid:            uuid.New().String(),
		CustomerUuid:    customerUUID.(string),
		IdeasName:       strings.TrimSpace(input.IdeasName),
		Industry:        input.Industry,
		Procedure:       "Ideas",
		ContentDetail:   input.ContentDetail,
		Value_Benefits:  input.Value_Benefits,
		View:            0,
		Is_Intellect:    input.Is_Intellect,
		Price:           input.Price,
		IsActive:        0,
		IsDelete:        0,
		PostDay:         time.Now(),
		CustomerName:    customerName.(string),
		CustomerEmail:   customerEmail.(string),
		Image:           input.Image,
		Image_Intellect: input.Image_Intellect,
		Status:          0,
	}

	if err := idea.Insert(context.TODO()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Đăng ý tưởng thất bại", "detail": err.Error()})
		return
	}

	notificationPayload := map[string]interface{}{
		"receiver_uuid": idea.CustomerUuid,
		"name":          idea.IdeasName,
		"order_uuid":    idea.Uuid,
		"type":          "idea_post",
	}

	jsonData, _ := json.Marshal(notificationPayload)
	notificationServiceURL := "https://cms-notification-app-709c465dc8fb.herokuapp.com/v1/notification" // Cập nhật theo host thật

	resp, err := http.Post(notificationServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Không thể gửi thông báo: %v", err)
		// Không return lỗi – chỉ log
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng ý tưởng thành công, chờ duyệt",
		"data": gin.H{
			"uuid":          idea.Uuid,
			"ideasname":     idea.IdeasName,
			"industry":      idea.Industry,
			"price":         idea.Price,
			"customer_name": idea.CustomerName,
			"post_day":      idea.PostDay,
		},
	})
}

// /////////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) DeleteMyIdea(c *gin.Context) {
	userModel := new(models.Ideas)
	favModel := new(models.Favorite)
	var reqUri request.DeleteUri
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	customerUuid, exists := c.Get("customer_uuid")
	if !exists || customerUuid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	condition := bson.M{
		"uuid":         reqUri.Uuid,
		"customeruuid": customerUuid,
	}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	user.IsDelete = constant.DELETE

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	favCond := bson.M{
		"post_uuid": reqUri.Uuid,
	}
	fav, err := favModel.FindOne(favCond)
	if err == nil && fav != nil {
		fav.IsDelete = constant.DELETE
		_, _ = fav.Update()
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "Delete successfully"))
}

// //////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) UpdateMyIdea(c *gin.Context) {
	userModel := new(models.Ideas)
	var reqUri request.UpdateUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	customerUuid, exists := c.Get("customer_uuid")
	if !exists || customerUuid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req request.UpdateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{
		"uuid":         reqUri.Uuid,
		"customeruuid": customerUuid,
	}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	if req.IdeasName != "" {
		user.IdeasName = req.IdeasName
	}
	if req.ContentDetail != "" {
		user.ContentDetail = req.ContentDetail
	}
	if req.Image_Intellect != nil {
		user.Image_Intellect = *req.Image_Intellect

		if *req.Image_Intellect == "" {
			user.Is_Intellect = 0
		} else {
			user.Is_Intellect = 1
		}
	}
	if req.Industry != "" {
		user.Industry = req.Industry
	}
	if req.Value_Benefits != "" {
		user.Value_Benefits = req.Value_Benefits
	}
	if len(req.Image) > 0 {
		user.Image = req.Image
	}
	if req.Price != 0 {
		user.Price = int(req.Price)
	}
	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// /////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) AcceptStatus(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var reqUri request.UpdateStatusUri
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	var req request.AcceptStatusRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	if *req.IsActive < 0 || *req.IsActive > 1 {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Stauts just can be set in range [0..1]"))
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := ideaModel.FindOne(condition)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Idea no found!"))
		return
	}

	user.IsActive = *req.IsActive

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	notificationPayload := map[string]interface{}{
		"receiver_uuid": user.CustomerUuid,
		"name":          user.IdeasName,
		"order_uuid":    user.Uuid,
		"type":          "accept_idea",
	}

	jsonData, _ := json.Marshal(notificationPayload)
	notificationServiceURL := "https://cms-notification-app-709c465dc8fb.herokuapp.com/v1/notification" // Cập nhật theo host thật

	resp, err := http.Post(notificationServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Không thể gửi thông báo: %v", err)
		// Không return lỗi – chỉ log
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

////////////////////////////////////////////////////////////////////////////////

func (ideaCtl IdeasController) RejectStatus(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var reqUri request.UpdateStatusUris
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	var req request.RejectStatusRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	if *req.IsDelete < 0 || *req.IsDelete > 1 {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Stauts just can be set in range [0..1]"))
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := ideaModel.FindOneAdmin(condition)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Idea no found!"))
		return
	}

	user.IsDelete = *req.IsDelete

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}

	notificationPayload := map[string]interface{}{
		"receiver_uuid": user.CustomerUuid,
		"name":          user.IdeasName,
		"order_uuid":    user.Uuid,
		"type":          "reject_idea", // có thể là "idea_post", "system", v.v.
	}

	jsonData, _ := json.Marshal(notificationPayload)
	notificationServiceURL := "https://cms-notification-app-709c465dc8fb.herokuapp.com/v1/notification" // Cập nhật theo host thật

	resp, err := http.Post(notificationServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Không thể gửi thông báo: %v", err)
		// Không return lỗi – chỉ log
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// /////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) DetailWatting(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var reqUri request.GetDetailUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	idea, err := ideaModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("uuid no found!"))
		return
	}
	response := request.GetDetailResponse{
		IdeasName:       idea.IdeasName,
		Industry:        idea.Industry,
		Procedure:       idea.Procedure,
		ContentDetail:   idea.ContentDetail,
		Value_Benefits:  idea.Value_Benefits,
		Is_Intellect:    idea.Is_Intellect,
		Price:           idea.Price,
		CustomerName:    idea.CustomerName,
		CustomerEmail:   idea.CustomerEmail,
		Image:           idea.Image,
		View:            0, // Không tăng view khi xem chi tiết ý tưởng đang chờ duyệt
		Image_Intellect: idea.Image_Intellect,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

// ////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) AddFavorite(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var req request.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Check model: %v", err) // Thêm dòng này
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	customerUuid, exists := c.Get("customer_uuid")
	if !exists || customerUuid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	customerUuidStr, ok := customerUuid.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid must be string"})
		return
	}

	// Kiểm tra post_uuid có tồn tại trong cơ sở dữ liệu không
	cond := bson.M{
		"uuid": req.PostUuid,
	}

	idea, err := ideaModel.FindOne(cond)
	if err != nil {
		log.Printf("Tìm ý tưởng: %v", err) // Thêm dòng này
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy ý tưởng"})
		return
	}

	favModel := new(models.Favorite)
	existCond := bson.M{
		"customer_uuid": customerUuidStr,
		"post_uuid":     req.PostUuid,
		"is_delete":     0,
	}
	existingFavorite, _ := favModel.FindOne(existCond)
	if existingFavorite != nil && existingFavorite.Uuid != "" {
		log.Printf("Kiểm tra ý tưởng: %v", existingFavorite.Uuid)
		c.JSON(http.StatusConflict, gin.H{"error": "ý tưởng đã có trong danh sách yêu thích"})
		return
	}
	fav := models.Favorite{
		Uuid:            util.GenerateUUID(),
		CustomerUuid:    customerUuidStr,
		PostUuid:        req.PostUuid,
		CreatedAt:       util.NowVN(),
		IsDelete:        0,
		IdeasName:       idea.IdeasName,
		Industry:        idea.Industry,
		Procedure:       idea.Procedure,
		ContentDetail:   idea.ContentDetail,
		Value_Benefits:  idea.Value_Benefits,
		View:            idea.View,
		Price:           idea.Price,
		IsActive:        idea.IsActive,
		PostDay:         idea.PostDay,
		CustomerName:    idea.CustomerName,
		Image:           idea.Image,
		Image_Intellect: idea.Image_Intellect,
		CustomerEmail:   idea.CustomerEmail,
	}

	_, err = fav.Insert()
	if err != nil {
		log.Printf("Insert Favorite Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể thêm vào danh sách yêu thích"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã thêm vào danh sách yêu thích"})
}

// /////////////////////////////////////////////////////////
func (ideaCtl IdeasController) ListFavorite(c *gin.Context) {
	var req request.GetListRequestFavorite
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	customerUuid, exists := c.Get("customer_uuid")
	if !exists || customerUuid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	customerUuidStr, ok := customerUuid.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid must be string"})
		return
	}

	// Bước 1: Lấy danh sách favorite của customer
	cond := bson.M{
		"customer_uuid": customerUuidStr,
		"is_delete":     0,
	}
	optionsQuery, page, limit := models.GetPagingOption(req.Page, req.Limit, req.Sort)

	favModel := new(models.Favorite)
	favs, err := favModel.Pagination(c, cond, optionsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể truy vấn danh sách yêu thích"})
		return
	}

	// Bước 2: Lấy list PostUuid trong favorite
	postUuids := make([]string, 0, len(favs))
	postUuidToFavorite := make(map[string]*models.Favorite)
	for _, fav := range favs {
		postUuids = append(postUuids, fav.PostUuid)
		postUuidToFavorite[fav.PostUuid] = fav
	}

	if len(postUuids) == 0 {
		// Không có favorite nào
		c.JSON(http.StatusOK, respond.SuccessPagination([]request.ListResponseFavorite{}, page, limit, 0, 0))
		return
	}

	// Bước 3: Query ideas có status == 0
	ideasModel := new(models.Ideas)
	ideaCond := bson.M{
		"uuid":   bson.M{"$in": postUuids},
		"status": 0,
	}
	ideas, err := ideasModel.Find(ideaCond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể truy vấn ideas"})
		return
	}

	// Bước 4: Build response chỉ với favorite có idea hợp lệ
	respData := make([]request.ListResponseFavorite, 0, len(ideas))
	for _, idea := range ideas {
		fav := postUuidToFavorite[idea.Uuid]
		image := ""
		if len(fav.Image) > 0 {
			image = fav.Image[0]
		}

		respData = append(respData, request.ListResponseFavorite{
			Uuid:          fav.Uuid,
			PostUuid:      fav.PostUuid,
			CreatedAt:     fav.CreatedAt,
			CustomerUuid:  fav.CustomerUuid,
			IdeasName:     fav.IdeasName,
			Industry:      fav.Industry,
			Procedure:     fav.Procedure,
			ContentDetail: fav.ContentDetail,
			View:          fav.View,
			Price:         fav.Price,
			PostDay:       fav.PostDay,
			CustomerName:  fav.CustomerName,
			Image:         image,
		})
	}

	// Tổng số bản ghi chính là số lượng ideas hợp lệ
	total := int64(len(respData))
	pages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}

func (ideaCtl IdeasController) DeleteFavorite(c *gin.Context) {
	userModel := new(models.Favorite)
	var reqUri request.DeleteUriFav

	// Bind chỉ lấy post_uuid từ URI thôi
	if err := c.ShouldBindUri(&reqUri); err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	// Lấy customer_uuid từ context (do middleware set)
	customerUuid, exists := c.Get("customer_uuid")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid is missing"})
		return
	}

	customerUuidStr, ok := customerUuid.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid must be string"})
		return
	}

	// Log ra cho bạn kiểm tra
	log.Printf("[DEBUG]", "post_uuid: %s, customer_uuid: %s\n", reqUri.PostUuid, customerUuidStr)

	condition := bson.M{
		"post_uuid":     reqUri.PostUuid,
		"customer_uuid": customerUuidStr,
	}

	fav, err := userModel.FindOne(condition)
	log.Printf("[DEBUG] Query condition: %+v\n", condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Favorite not found!"))
		return
	}

	fav.IsDelete = constant.DELETE

	_, err = fav.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}

	c.JSON(http.StatusOK, respond.Success(fav.PostUuid, "Delete successfully"))
}

// /////////////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) ListID(c *gin.Context) {
	ideaModel := new(models.Favorite)
	// Điều kiện lọc
	customerUuid, exists := c.Get("customer_uuid")
	if !exists || customerUuid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	customerUuidStr, ok := customerUuid.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid must be string"})
		return
	}
	cond := bson.M{"customer_uuid": customerUuidStr}
	// Truy vấn toàn bộ không giới hạn
	ideas, err := ideaModel.Find(cond)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	// Lấy UUID
	respData := make([]request.ListResponseUuid, 0)
	for _, idea := range ideas {
		respData = append(respData, request.ListResponseUuid{
			PostUuid: idea.PostUuid,
		})
	}
	// Không cần phân trang, chỉ trả về thành công
	c.JSON(http.StatusOK, respond.Success(respData, "Successfully"))
}

// ////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) Price(c *gin.Context) {

	ideaModel := new(models.Ideas)
	var reqUri request.GetPricelUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	idea, err := ideaModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("uuid no found!"))
		return
	}
	response := request.GetPriceResponse{
		Price: idea.Price,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}
func (ideaCtl IdeasController) UpdateStatus(c *gin.Context) {
	ideaUuid := c.Param("uuid")

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status là bắt buộc"})
		return
	}

	ideaModel := new(models.Ideas)
	condition := bson.M{"uuid": ideaUuid}
	update := bson.M{"$set": bson.M{"status": req.Status}}

	_, err := ideaModel.UpdateByCondition(condition, update)
	if err != nil {
		log.Println("[IdeasService] Lỗi cập nhật status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật status idea"})
		return
	}
	idea, err := ideaModel.FindOne(condition)
	if err != nil {
		log.Println("[IdeasService] Không tìm thấy ý tưởng sau khi cập nhật")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tìm thấy ý tưởng sau khi cập nhật"})
		return
	}
	notificationPayload := map[string]interface{}{
		"receiver_uuid": idea.CustomerUuid,
		"name":          idea.IdeasName,
		"order_uuid":    idea.Uuid,
		"type":          "buy_idea", // có thể là "idea_post", "system", v.v.
	}

	jsonData, _ := json.Marshal(notificationPayload)
	notificationServiceURL := "https://cms-notification-app-709c465dc8fb.herokuapp.com/v1/notification" // Cập nhật theo host thật

	resp, err := http.Post(notificationServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Không thể gửi thông báo: %v", err)
		// Không return lỗi – chỉ log
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật status idea thành công"})
}
func fetchPurchasedIdeaUUIDs(customerUUID string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("https://cms-payment-app-358737dd8384.herokuapp.com/v1/payment/success/ideas/%s", customerUUID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service trả về mã lỗi: %d", resp.StatusCode)
	}

	var result struct {
		ProduceUUIDs []string `json:"produce_uuids"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.ProduceUUIDs, nil
}

// ///////////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) GetPurchasedIdeas(c *gin.Context) {
	var req request.GetListRequest
	err := c.ShouldBindWith(&req, binding.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	// Kiểm tra customer_uuid có được truyền không
	customerUuidVal, exists := c.Get("customer_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin người dùng trong token"})
		return
	}
	customerUuid := customerUuidVal.(string)

	// Gọi Payment Service để lấy danh sách uuid đã mua
	uuids, err := fetchPurchasedIdeaUUIDs(customerUuid)
	if err != nil {
		log.Printf("Lỗi gọi payment service cho customer %s: %v", customerUuid, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy dữ liệu payment"})
		return
	}

	// Nếu danh sách mua trống thì trả về rỗng luôn
	if len(uuids) == 0 {
		c.JSON(http.StatusOK, respond.SuccessPagination([]request.ListBuyResponse{}, 1, 10, 0, 0))
		return
	}

	// Điều kiện truy vấn
	cond := bson.M{
		"uuid": bson.M{"$in": uuids},
	}

	// Phân trang và sắp xếp
	optionsQuery, page, limit := models.GetPagingOption(req.Page, req.Limit, req.Sort)
	if req.Sort == "" {
		optionsQuery.SortBy = "post_day"
		optionsQuery.SortDir = -1
	}

	// Lấy danh sách ideas từ DB
	ideaModel := new(models.Ideas)

	// Chuẩn hóa dữ liệu trả về
	respData := make([]request.ListBuyResponse, 0)

	ideas, err := ideaModel.Pagination(c, cond, optionsQuery)
	if err != nil {
		log.Printf("Lỗi truy vấn Pagination ideas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách ý tưởng"})
		return
	}

	for _, idea := range ideas {
		image := ""
		if len(idea.Image) > 0 {
			image = idea.Image[0]
		}
		res := request.ListBuyResponse{
			Uuid:          idea.Uuid,
			CustomerName:  idea.CustomerName,
			IdeasName:     idea.IdeasName,
			Industry:      idea.Industry,
			Image:         image,
			ContentDetail: idea.ContentDetail,
			PostDay:       idea.PostDay,
		}
		respData = append(respData, res)
	}

	// Tổng số và số trang
	total, err := ideaModel.Count(c, cond)
	if err != nil {
		log.Printf("Lỗi Count ideas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đếm tổng số ý tưởng"})
		return
	}
	pages := int(math.Ceil(float64(total) / float64(limit)))

	// Trả về
	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) UpdateStatusAdmin(c *gin.Context) {
	userModel := new(models.Ideas)
	var reqUri request.UpdateStatusUriAdmin
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	var req request.UpdateStatusRequestAdmin
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOneAdmin(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotFound, respond.ErrorCommon("Ideas not found!"))
		return
	}
	if req.IsDelete != nil {
		user.IsDelete = *req.IsDelete
	}
	// Cập nhật DB
	if _, err := user.Update(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, respond.UpdatedFail())
		return
	}
	notificationPayload := map[string]interface{}{
		"receiver_uuid": user.CustomerUuid,
		"order_uuid":    user.Uuid,
		"name":          user.IdeasName,
		"type":          "delete_idea", // có thể là "idea_post", "system", v.v.
	}

	jsonData, _ := json.Marshal(notificationPayload)
	notificationServiceURL := "https://cms-notification-app-709c465dc8fb.herokuapp.com/v1/notification" // Cập nhật theo host thật

	resp, err := http.Post(notificationServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Không thể gửi thông báo: %v", err)
		// Không return lỗi – chỉ log
	}

	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) ListForAdmin(c *gin.Context) {
	userModel := new(models.Ideas)
	var req request.GetListWattingRequest
	err := c.ShouldBindWith(&req, binding.Query)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	cond := bson.M{}
	if req.Ideasname != nil {
		cond["ideasname"] = bson.M{
			"$regex":   *req.Ideasname,
			"$options": "i",
		}
	}
	if req.Industry != nil && *req.Industry != "" {
		cond["industry"] = *req.Industry
	}

	if c.Query("is_active") != "" && req.IsActive != nil {
		cond["is_active"] = *req.IsActive
	}

	if c.Query("is_delete") != "" && req.IsDelete != nil {
		cond["is_delete"] = *req.IsDelete
	}

	optionsQuery, page, limit := models.GetPagingOption(req.Page, 10, req.Sort)
	var respData []request.ListWattingResponse
	users, err := userModel.PaginationAdmin(c, cond, optionsQuery)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, respond.SuccessPagination([]request.ListWattingResponse{}, page, limit, 0, 0))
		return
	}
	for _, user := range users {
		res := request.ListWattingResponse{
			Uuid:         user.Uuid,
			IsActive:     user.IsActive,
			IdeasName:    user.IdeasName,
			IsDelete:     user.IsDelete,
			Industry:     user.Industry,
			PostDay:      user.PostDay,
			CustomerName: user.CustomerName,
		}
		respData = append(respData, res)
	}
	total, err := userModel.Count(c, cond)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	pages := int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) TotalSoldIdeasPrice(c *gin.Context) {
	ideaModel := new(models.Ideas)
	coll := ideaModel.Model()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "status", Value: int32(1)}, {Key: "is_delete", Value: 0}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total_price", Value: bson.D{{Key: "$sum", Value: "$price"}}},
		}}},
	}

	cursor, err := coll.Aggregate(c, pipeline)
	if err != nil {
		log.Println("Lỗi khi aggregate tổng giá:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tính tổng giá"})
		return
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Println("Lỗi khi đọc dữ liệu aggregate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi đọc dữ liệu"})
		return
	}

	var total int32 = 0
	if len(results) > 0 {
		if val, ok := results[0]["total_price"].(int32); ok {
			total = val
		} else if valFloat, ok := results[0]["total_price"].(float64); ok {
			total = int32(valFloat)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Tổng giá trị các ý tưởng đã bán",
		"total_price": total,
	})
}

// /////////////////////////////////////
func (ideaCtl IdeasController) DetailForUpdate(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var reqUri request.GetDetailUriForUpdate

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	customerUuid, exists := c.Get("customer_uuid")
	if !exists || customerUuid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Tìm ý tưởng theo uuid và customer_uuid (đảm bảo là của người dùng hiện tại)
	cond := bson.M{
		"customeruuid": customerUuid,
		"uuid":         reqUri.Uuid,
		"status":       0,
	}
	idea, err := ideaModel.FindOne(cond)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Idea not found"})
		return
	}
	response := request.GetDetailResponseForUpdate{
		IdeasName:      idea.IdeasName,
		Industry:       idea.Industry,
		Procedure:      idea.Procedure,
		ContentDetail:  idea.ContentDetail,
		Value_Benefits: idea.Value_Benefits,
		Is_Intellect:   idea.Is_Intellect,
		Price:          idea.Price,
		CustomerName:   idea.CustomerName,
		CustomerEmail:  idea.CustomerEmail,
		CustomerUuid:   idea.CustomerUuid,
		PostDay:        idea.PostDay,
		IsActive:       idea.IsActive,
		IsDelete:       idea.IsDelete,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

// //////////////////////////////////////////////////////////////////////////////////////////////
func (ideaCtl IdeasController) CountForAdmin(c *gin.Context) {
	userModel := new(models.Ideas)

	// Đếm tổng số ý tưởng
	total, err := userModel.Count(c, bson.M{})
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	// Đếm số ý tưởng đã mua (status = 1)
	totalIdeasPurchased, err := userModel.Count(c, bson.M{"status": 1})
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	// Trả kết quả gộp
	result := gin.H{
		"total":              total,
		"total_ideas_bought": totalIdeasPurchased,
	}

	c.JSON(http.StatusOK, respond.Success(result, "Thống kê ý tưởng"))
}
