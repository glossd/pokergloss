package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/bank/services"
	"github.com/glossd/pokergloss/bank/services/model"
	"net/http"
)

// @ID get rating page
// @Param pageSize query int true "How many entities should by in page"
// @Param pageNumber query int false "The number of page. Starts from 1. if not defined the request will return page with user"
// @Success 200 {object} model.PageRating
// @Failure 400 {object} ErrorRes
// @Router /ratings/page [get]
func GetRatings(c *gin.Context) {
	iden, _ := auth.IdSafe(c)
	pageSize := c.Request.URL.Query().Get("pageSize")
	if pageSize == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Efmt("pageSize is required"))
		return
	}
	pageNumber := c.Request.URL.Query().Get("pageNumber")
	pr, err := model.NewPageRequest(pageSize, pageNumber)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	page, err := services.GetRatingPage(c.Request.Context(), iden, *pr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	c.JSON(http.StatusOK, page)
}

// @ID get user rating
// @Success 200 {object} model.Rating
// @Failure 400 {object} ErrorRes
// @Router /ratings/me [get]
func GetUserRating(c *gin.Context) {
	rating, err := services.GetRating(c.Request.Context(), auth.Id(c).UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	c.JSON(http.StatusOK, model.ToRating(rating))
}
