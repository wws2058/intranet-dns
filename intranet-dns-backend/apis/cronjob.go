package apis

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"github.com/wws2058/intranet-dns/ctx"
	"github.com/wws2058/intranet-dns/models"
	"github.com/wws2058/intranet-dns/service/cronjob"
	"github.com/wws2058/intranet-dns/utils"
)

type newCronjob struct {
	Name        string      `binding:"gt=0" json:"name,omitempty"`        // task name
	Spec        string      `binding:"gte=9" json:"spec,omitempty"`       // spec & s m h D M W
	Description string      `binding:"gt=0" json:"description,omitempty"` // task desc
	TaskType    string      `binding:"gt=0" json:"task_type,omitempty"`   // task type
	TaskArgs    models.Args `json:"task_args,omitempty"`                  // task args
}

// @Summary  add cronjob
// @Tags     cronjob
// @Produce  json
// @Param    token    header  string            false  "jwt token"
// @Param    request  body    newCronjob       true   "new cronjob request"
// @Success  200      object  ctx.StdResponse  "new cronjob id"
// @Router   /api/v1/cronjobs [POST]
func addCronjob(c *gin.Context) {
	// name uniq_key
	req := &newCronjob{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	if _, err := cron.ParseStandard(req.Spec); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	username := ctx.GetLoginUsername(c)

	job := &models.Cronjob{
		Name:        req.Name,
		Spec:        req.Spec,
		Creator:     username,
		Description: req.Description,
		Started:     true,
		LastSucceed: true,
		TaskType:    req.TaskType,
		TaskArgs:    req.TaskArgs,
	}
	if err := models.TemplateCreate(job); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	cronjob.RefreshRegisterJobTs()
	ctx.SucceedRsp(c, job.Id, nil)
}

// @Summary  list cronjob
// @Tags     cronjob
// @Produce  json
// @Param    token  header  string           false  "jwt token"
// @Param    page          query   int              false  "page, min=1"
// @Param    page_size     query   int              false  "page size, min=10, max=1000"
// @Param    name          query   string           false  "cronjob name"
// @Param    creator       query   string           false  "cronjob creator"
// @Param    task_type     query   string           false  "cronjob type"
// @Param    started       query   bool             false  "cronjob running status"
// @Param    last_succeed  query   bool             false  "cronjob last running status"
// @Success  200           object  ctx.StdResponse  "cronjobs"
// @Router   /api/v1/cronjobs [GET]
func listCronjob(c *gin.Context) {
	var req struct {
		Name        string `form:"name"`
		Creator     string `form:"creator"`
		TaskType    string `form:"task_type"`
		Started     *bool  `form:"started"`
		LastSucceed *bool  `form:"last_succeed"`
		models.PageReq
	}
	if err := c.BindQuery(&req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	filter := make(map[string]interface{})
	if len(req.Name) > 0 {
		filter["name"] = req.Name
	}
	if len(req.Creator) > 0 {
		filter["creator"] = req.Creator
	}
	if len(req.TaskType) > 0 {
		filter["task_type"] = req.TaskType
	}
	if req.Started != nil {
		filter["started"] = *req.Started
	}
	if req.LastSucceed != nil {
		filter["last_succeed"] = *req.LastSucceed
	}

	jobs := []models.Cronjob{}
	pageQuery := &models.DaoDBReq{
		PageReq:     models.PageReq{Page: req.Page, PageSize: req.PageSize},
		PageRsp:     models.PageRsp{},
		Dst:         &jobs,
		ModelFilter: filter,
		OrderBy:     "id desc",
	}
	err := models.TemPlatePageQuery(pageQuery)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, jobs, &pageQuery.PageRsp)
}

// @Summary  del cronjob
// @Tags     cronjob
// @Produce  json
// @Param    token  header  string           false  "jwt token"
// @Param    id     path    int              true   "cronjob id"
// @Success  200    object  ctx.StdResponse  "cronjob id"
// @Router   /api/v1/cronjobs/{id} [DELETE]
func delJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}

	opt := models.DaoDBReq{
		Dst: &models.Cronjob{
			BaseModel: models.BaseModel{Id: uint(id)},
		},
	}
	if err := models.TemplateSoftDelete(opt); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	cronjob.RefreshRegisterJobTs()
	ctx.SucceedRsp(c, id, nil)
}

type updateCronjobReq struct {
	Id          uint         `binding:"required,gt=0" json:"id,omitempty"` // user id
	Name        string       `json:"name,omitempty"`                       // task name
	Spec        string       `json:"spec,omitempty"`                       // spec & s m h D M W
	Started     *bool        `json:"started,omitempty"`                    // task switch
	Description string       `json:"description,omitempty"`                // task desc
	TaskArgs    *models.Args `json:"task_args,omitempty"`                  // task args
}

// @Summary  update cronjob
// @Tags     cronjob
// @Produce  json
// @Param    token    header  string           false  "jwt token"
// @Param    request  body    updateCronjobReq  true   "update cronjob request"
// @Success  200      object  ctx.StdResponse   "cronjob id"
// @Router   /api/v1/cronjobs [PUT]
func updateCronjob(c *gin.Context) {
	req := &updateCronjobReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	job := &models.Cronjob{
		BaseModel: models.BaseModel{Id: req.Id},
	}
	fields := []string{}
	if len(req.Name) > 0 {
		job.Name = req.Name
		fields = append(fields, "name")
	}
	if len(req.Spec) > 0 {
		if _, err := cron.ParseStandard(req.Spec); err != nil {
			ctx.FailedRsp(c, err)
			return
		}
		job.Spec = req.Spec
		job.History = nil
		fields = append(fields, "spec")
		fields = append(fields, "history")
	}
	if len(req.Description) > 0 {
		job.Description = req.Description
		fields = append(fields, "description")
	}
	if req.TaskArgs != nil {
		job.TaskArgs = *req.TaskArgs
		job.History = nil
		fields = append(fields, "task_args")
		fields = append(fields, "history")
	}
	if req.Started != nil {
		job.Started = *req.Started
		fields = append(fields, "started")
	}

	if err := models.TemplateUpdate(job, utils.RemoveRepeatedElement[string](fields)); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	cronjob.RefreshRegisterJobTs()
	ctx.SucceedRsp(c, req.Id, nil)
}

// @Summary  list cronjob functions
// @Tags     cronjob
// @Produce  json
// @Param    token         header  string           false  "jwt token"
// @Success  200    object  ctx.StdResponse  "cronjob functions"
// @Router   /api/v1/cronjobs/functions [GET]
func listCronjobFunctions(c *gin.Context) {
	ctx.SucceedRsp(c, cronjob.GetInternalFunctions(), nil)
}

func LoadCronjobApis(r *gin.Engine) {
	apis := []models.Api{
		// system user manage
		{Path: "/cronjobs/:id", Method: http.MethodDelete, Description: "删除定时任务", Handler: delJob},
		{Path: "/cronjobs", Method: http.MethodGet, Description: "列举定时任务", Handler: listCronjob},
		{Path: "/cronjobs", Method: http.MethodPost, Description: "新增定时任务", Handler: addCronjob},
		{Path: "/cronjobs", Method: http.MethodPut, Description: "更新定时任务", Handler: updateCronjob},
		{Path: "/cronjobs/functions", Method: http.MethodGet, Description: "列举内置可用定时任务函数", Handler: listCronjobFunctions},
	}
	loadApi(r, ginGroupApiV1, apis)
}
