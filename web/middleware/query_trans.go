package middleware

import (
	"gcrontab/constant"
	"gcrontab/custom"
	"gcrontab/model"
	"gcrontab/model/requestmodel"
	"gcrontab/utils"
	"gcrontab/web/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func QueryTrans() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageQ := c.Query(constant.PAGE)
		pageSizeQ := c.Query(constant.PAGESIZE)
		nameQ := c.Query(constant.NAME)
		statusQ := c.Query(constant.STATUS)
		createrQ := c.Query(constant.CREATER)
		startTimeQ := c.Query(constant.STARTTIME)
		endTimeQ := c.Query(constant.ENDTIME)
		sortedByQ := c.Query(constant.SORTEDBY)
		orderQ := c.Query(constant.ORDER)

		if orderQ != constant.ASC && orderQ != constant.DESC && orderQ != "" {
			logrus.Errorf("query order invalid:%v", orderQ)
			c.AbortWithStatusJSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ErrorInvalideRequest.Error()))
			return
		}
		createrNameQ := c.Query(constant.CREATERNAME)
		page, pageSize, err := atoiHandler(pageQ, pageSizeQ)
		if err != nil {
			logrus.Errorf("convert page[%s],pageSize[%s] to int failed:%v", pageQ, pageSizeQ, err)
			c.AbortWithStatusJSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ErrorInvalideRequest.Error()))
			return
		}
		logTaskIDQ := c.Query(constant.LOGTASKID)
		timeStampQ := c.Query(constant.TIMESTAMP)

		params := requestmodel.Params{}
		params.Page = page
		params.PageSize = pageSize
		params.Name = nameQ
		params.Status = statusQ
		params.Creater = createrQ
		params.LogTaskID = logTaskIDQ

		if timeStampQ != "" {
			timeStampint, err := strconv.Atoi(timeStampQ)
			if err != nil {
				logrus.Errorf("convert string[%s] to int failed:%v", timeStampQ, err)
				c.AbortWithStatusJSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ParamErrorReturn(constant.TIMESTAMP).Error()))
				return
			}
			params.TimeStamp = int64(timeStampint)
		}

		if startTimeQ != "" {
			ts, err := time.ParseInLocation(constant.TIMELAYOUT, startTimeQ, utils.DefaultLocation)
			if err != nil {
				logrus.Errorf("%s error", startTimeQ)
				c.AbortWithStatusJSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ParamErrorReturn(constant.STARTTIME).Error()))
				return
			}
			params.StartTime = ts.In(utils.DefaultLocation)
		}
		if endTimeQ != "" {
			te, err := time.ParseInLocation(constant.TIMELAYOUT, endTimeQ, utils.DefaultLocation)
			if err != nil {
				logrus.Errorf("%s error", endTimeQ)
				c.AbortWithStatusJSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ParamErrorReturn(constant.ENDTIME).Error()))
				return
			}
			params.EndTime = te.In(utils.DefaultLocation)
		}
		params.SortedBy = sortedByQ
		params.Order = orderQ
		params.CreaterName = createrNameQ

		c.Set("params", params)
	}
}

func atoiHandler(pageQ, pageSizeQ string) (int, int, error) {
	var page, pageSize int
	var err error
	if pageQ != "" {
		page, err = strconv.Atoi(pageQ)
		if err != nil {
			logrus.Errorf("page's type need be int : %v", pageQ)
			return page, pageSize, custom.ParamErrorReturn("page")
		}
	}

	if page <= 0 {
		page = 1
	}

	if pageSizeQ != "" {
		pageSize, err = strconv.Atoi(pageSizeQ)
		if err != nil || pageSize > model.PageSizeLimit {
			logrus.Errorf("pageSize's type need be int and not greater than %d: %v", model.PageSizeLimit, pageSizeQ)
			return page, pageSize, custom.ParamErrorReturn("pageSize")
		}
	}

	if pageSize <= 0 {
		pageSize = model.PageSizeDefault
	}
	return page, pageSize, nil
}
