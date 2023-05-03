package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	"github.com/xyedo/blindate/pkg/domain/online"
	"github.com/xyedo/blindate/pkg/infrastructure"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func New(config infrastructure.Config, onlineUsecase online.Usecase) *onlineH {
	return &onlineH{
		config:   config,
		onlineUC: onlineUsecase,
	}
}

type onlineH struct {
	config   infrastructure.Config
	onlineUC online.Usecase
}

func (h *onlineH) postOnlineHandler(c *gin.Context) {
	err := h.onlineUC.Create(c.GetString(constant.KeyRequestUserId))
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "online created",
	})
}

func (h *onlineH) getOnlineHandler(c *gin.Context) {
	var url struct {
		UserId string `uri:"userId" binding:"required,uuid4"`
	}
	err := c.ShouldBindUri(&url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "must have uuid in uri!",
		})
		return
	}
	online, err := h.onlineUC.GetByUserId(c.GetString(constant.KeyRequestUserId), url.UserId)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"online": getUserOnlineResponse(online),
		},
	})
}
