package apis

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tswcbyy1107/dns-service/models"
)

// @Summary      list api
// @Description  page query apis by params
// @Tags     system
// @Produce  json
// @Param        page       query   int          false  "page, min=1"
// @Param        page_size  query   int          false  "page size, min=10, max=1000"
// @Param        path       query   string       false  "api path"
// @Param        method     query   string       false  "api method"
// @Param        active     query   bool         false  "api activated"
// @Success      200        object  StdResponse  "apis"
// @Router       /api/v1/apis [GET]
func listApis(c *gin.Context) {
	var req struct {
		Path   string `form:"path"`
		Method string `form:"method"`
		Active *bool  `form:"active"`
		models.PageReq
	}
	if err := c.BindQuery(&req); err != nil {
		failedRsp(c, models.FormatErr(models.ErrParams, err))
		return
	}

	filter := make(map[string]interface{})
	if len(req.Path) > 0 {
		filter["path"] = req.Path
	}
	if len(req.Method) > 0 {
		filter["method"] = req.Path
	}
	if req.Active != nil {
		filter["active"] = *req.Active
	}
	apis := []models.Api{}
	pageQuery := &models.DaoDBReq{
		PageReq:     models.PageReq{Page: req.Page, PageSize: req.PageSize},
		PageRsp:     models.PageRsp{},
		Dst:         &apis,
		ModelFilter: filter,
		OrderBy:     "path",
	}
	err := models.TemPlatePageQuery(pageQuery)
	if err != nil {
		failedRsp(c, models.ErrDbQuery)
		return
	}
	succeedRsp(c, apis, &pageQuery.PageRsp)
}

type updateApiReq struct {
	Id     uint  `json:"id" binding:"required"`
	Audit  *bool `json:"audit,omitempty"`
	Active *bool `json:"active,omitempty"`
}

// @Summary      update api
// @Description  api's active and audit attributes
// @Tags     system
// @Produce  json
// @Param        request  body    updateApiReq  false  "id: api db id; audit: true 1, request logged; active: true 1, in use"
// @Success      200      object  StdResponse   "api updated"
// @Router       /api/v1/api [PUT]
func updateApi(c *gin.Context) {
	req := &updateApiReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		failedRsp(c, models.FormatErr(models.ErrParams, err))
		return
	}
	api := models.Api{
		BaseModel: models.BaseModel{Id: req.Id},
	}
	fields := []string{}
	if req.Audit != nil {
		api.Audit = *req.Audit
		fields = append(fields, "audit")
	}
	if req.Active != nil {
		api.Active = *req.Active
		fields = append(fields, "active")
	}
	if err := models.TemplateUpdate(api, fields); err != nil {
		failedRsp(c, models.ErrDbUpdate)
		return
	}
	succeedRsp(c, "ok", nil)
}

// @Summary      list system roles
// @Description  get all system roles in pages
// @Tags     system
// @Produce  json
// @Param        page       query   int          false  "min=1"
// @Param        page_size  query   int          false  "min=10, max=1000"
// @Param        name_cn    query   string       false  "role chinese name"
// @Success      200        object  StdResponse  "roles"
// @Router       /api/v1/roles [GET]
func listSysRoles(c *gin.Context) {
	var req struct {
		NameCN string `json:"name_cn,omitempty" form:"name_cn"`
		models.PageReq
	}
	if err := c.BindQuery(&req); err != nil {
		failedRsp(c, models.FormatErr(models.ErrParams, err))
		return
	}

	filter := make(map[string]interface{})
	if len(req.NameCN) > 0 {
		filter["name_cn"] = req.NameCN
	}

	roles := []models.SysRole{}
	pageQuery := &models.DaoDBReq{
		PageReq:     models.PageReq{Page: req.Page, PageSize: req.PageSize},
		PageRsp:     models.PageRsp{},
		Dst:         &roles,
		ModelFilter: filter,
		OrderBy:     "id desc",
	}
	err := models.TemPlatePageQuery(pageQuery)
	if err != nil {
		failedRsp(c, models.ErrDbQuery)
		return
	}
	succeedRsp(c, roles, &pageQuery.PageRsp)
}

// @Summary  role detail
// @Tags         system
// @Produce      json
// @Param    id   path    int          true  "role id"
// @Success  200  object  StdResponse  "role detail with accessible apis"
// @Router   /api/v1/roles/{id}/apis [GET]
func roleDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		failedRsp(c, models.FormatErr(models.ErrParams, "id required"))
		return
	}
	filter := make(map[string]interface{})
	filter["id"] = id
	role := &models.SysRole{}
	opt := &models.DaoDBReq{
		Dst:         role,
		ModelFilter: filter,
	}
	if err := models.TemplateQuery(opt); err != nil {
		failedRsp(c, models.ErrDbQuery)
		return
	}
	if err := role.ApiDetails(); err != nil {
		failedRsp(c, models.ErrDbQuery)
		return
	}
	succeedRsp(c, role, nil)
}

// @Summary  del role
// @Tags         system
// @Produce      json
// @Param    id   path    int          true  "role id"
// @Success  200  object  StdResponse  "role id"
// @Router   /api/v1/roles/{id} [DELETE]
func delRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		failedRsp(c, models.FormatErr(models.ErrParams, "id required"))
		return
	}
	opt := models.DaoDBReq{
		Dst: &models.SysRole{BaseModel: models.BaseModel{Id: uint(id), Deleted: 1}},
	}

	if err := models.TemplateSoftDelete(opt); err != nil {
		failedRsp(c, models.ErrDbUpdate)
		return
	}
	succeedRsp(c, nil, nil)
}

// @Summary  add role
// @Tags         system
// @Produce      json
// @Param    request  body    models.SysRole  true  "role request"
// @Success  200      object  StdResponse     "role id"
// @Router   /api/v1/roles [POST]
func addRole(c *gin.Context) {
	req := &models.SysRole{}
	if err := c.ShouldBindJSON(req); err != nil {
		failedRsp(c, models.FormatErr(models.ErrParams, err))
		return
	}
	if err := models.TemplateCreate(req); err != nil {
		failedRsp(c, models.ErrDbInsert)
		return
	}
	succeedRsp(c, req.Id, nil)
}

type updateRoleReq struct {
	Id     uint    `json:"id" binding:"required"`            // role id
	Name   *string `json:"name,omitempty" binding:"gt=0"`    // role en name
	NameCn *string `json:"name_cn,omitempty" binding:"gt=0"` // role cn name
	ApiIds *[]uint `json:"api_ids,omitempty" binding:"gt=0"` // role accessible apis id
}

// @Summary  update role
// @Tags         system
// @Produce      json
// @Param    request  body    updateRoleReq  true  "role request"
// @Success  200      object  StdResponse     "role id"
// @Router   /api/v1/roles [PUT]
func updateRole(c *gin.Context) {
	req := &updateRoleReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		failedRsp(c, models.FormatErr(models.ErrParams, err))
		return
	}
	role := &models.SysRole{
		BaseModel: models.BaseModel{Id: req.Id},
	}
	fields := []string{}
	if req.Name != nil {
		role.Name = *req.Name
		fields = append(fields, "name")
	}
	if req.NameCn != nil {
		role.NameCn = *req.NameCn
		fields = append(fields, "name_cn")
	}
	if req.ApiIds != nil {
		role.ApiIds = *req.ApiIds
		fields = append(fields, "api_ids")
	}
	if err := models.TemplateUpdate(role, fields); err != nil {
		failedRsp(c, models.ErrDbUpdate)
		return
	}
	succeedRsp(c, req.Id, nil)
}

func listUser(c *gin.Context) {

}

func delUser(c *gin.Context) {

}

func updateUser(c *gin.Context) {

}

func addUser(c *gin.Context) {

}

// api相关接口
func LoadSysApis(r *gin.Engine) {
	apis := []models.Api{
		// system api manage
		{Path: "/apis", Method: http.MethodGet, Description: "列举服务本身所有api接口", Handler: listApis},
		{Path: "/api", Method: http.MethodPut, Description: "更新相关api接口", Handler: updateApi},

		// system role manage
		{Path: "/roles/:id/apis", Method: http.MethodGet, Description: "获取角色api权限详情", Handler: roleDetail},
		{Path: "/roles/:id", Method: http.MethodDelete, Description: "删除角色", Handler: delRole},
		{Path: "/roles", Method: http.MethodGet, Description: "列举系统角色", Handler: listSysRoles},
		{Path: "/roles", Method: http.MethodPost, Description: "新增角色", Handler: addRole},
		{Path: "/roles", Method: http.MethodPut, Description: "更新角色", Handler: updateRole},

		// system use manage
		{Path: "/users/:id", Method: http.MethodDelete, Description: "删除用户", Handler: delUser},
		{Path: "/users", Method: http.MethodGet, Description: "列举系统用户", Handler: listUser},
		{Path: "/users", Method: http.MethodPost, Description: "新增用户", Handler: addUser},
		{Path: "/users", Method: http.MethodPut, Description: "更新用户", Handler: updateUser},
	}
	loadApi(r, ginGroupApiV1, apis)
}
