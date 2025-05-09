package controllers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/anhhuy1010/DATN-cms-ideas/helpers/respond"
	"github.com/anhhuy1010/DATN-cms-ideas/models"
	request "github.com/anhhuy1010/DATN-cms-ideas/request/ideas"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
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
	cond := bson.M{}
	if req.Ideasname != nil {
		cond["ideasname"] = req.Ideasname
	}

	if req.IsActive != nil {
		cond["is_active"] = req.IsActive
	}

	optionsQuery, page, limit := models.GetPagingOption(req.Page, req.Limit, req.Sort)
	var respData []request.ListResponse
	ideas, err := ideaModel.Pagination(c, cond, optionsQuery)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	for _, user := range ideas {
		res := request.ListResponse{
			Uuid:     user.Uuid,
			IsActive: user.IsActive,
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

func (ideaCtl IdeasController) Detail(c *gin.Context) {
	ideaModel := new(models.Ideas)
	var reqUri request.GetDetailUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"industry": reqUri.Industry}
	idea, err := ideaModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Industry no found!"))
		return
	}

	response := request.GetDetailResponse{
		IdeasName:      idea.IdeasName,
		Industry:       idea.Industry,
		OrtherIndustry: idea.OrtherIndustry,
		Procedure:      idea.Procedure,
		ContentDetail:  idea.ContentDetail,
		Value_Benefits: idea.Value_Benefits,
		Is_Intellect:   idea.Is_Intellect,
		Price:          idea.Price,
		CustomerName:   idea.CustomerName,
		CustomerEmail:  idea.CustomerEmail,
		Image:          idea.Image,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}
