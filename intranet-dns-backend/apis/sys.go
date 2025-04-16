package apis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tswcbyy1107/intranet-dns/ctx"
	"github.com/tswcbyy1107/intranet-dns/models"
	"github.com/tswcbyy1107/intranet-dns/service/jwt"
)

// @Summary  list api
// @Tags     system
// @Produce  json
// @Param    token      header  string           false  "jwt token"
// @Param    page       query   int              false  "page, min=1"
// @Param    page_size  query   int              false  "page size, min=10, max=1000"
// @Param    path       query   string           false  "api path"
// @Param    method     query   string           false  "api method"
// @Param    active     query   bool             false  "api activated"
// @Success  200        object  ctx.StdResponse  "apis"
// @Router   /api/v1/apis [GET]
func listApis(c *gin.Context) {
	var req struct {
		Path   string `form:"path"`
		Method string `form:"method"`
		Active *bool  `form:"active"`
		models.PageReq
	}
	if err := c.BindQuery(&req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}

	filter := make(map[string]interface{})
	if len(req.Path) > 0 {
		filter["path"] = req.Path
	}
	if len(req.Method) > 0 {
		filter["method"] = req.Method
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
		OrderBy:     "id desc",
	}
	err := models.TemPlatePageQuery(pageQuery)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, apis, &pageQuery.PageRsp)
}

type updateApiReq struct {
	Id     uint  `json:"id" binding:"required"`
	Audit  *bool `json:"audit,omitempty"`
	Active *bool `json:"active,omitempty"`
}

// @Summary  update api
// @Tags     system
// @Produce  json
// @Param    token    header  string           false  "jwt token"
// @Param    request  body    updateApiReq     false  "req"
// @Success  200      object  ctx.StdResponse  "api updated"
// @Router   /api/v1/apis [PUT]
func updateApi(c *gin.Context) {
	req := &updateApiReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
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
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, req.Id, nil)
}

// @Summary  list system roles
// @Tags     system
// @Produce  json
// @Param    token      header  string           false  "jwt token"
// @Param    page       query   int              false  "min=1"
// @Param    page_size  query   int              false  "min=10, max=1000"
// @Param    name_cn    query   string           false  "role chinese name"
// @Success  200        object  ctx.StdResponse  "roles"
// @Router   /api/v1/roles [GET]
func listSysRoles(c *gin.Context) {
	var req struct {
		NameCN string `json:"name_cn,omitempty" form:"name_cn"`
		models.PageReq
	}
	if err := c.BindQuery(&req); err != nil {
		ctx.FailedRsp(c, err)
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
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, roles, &pageQuery.PageRsp)
}

// @Summary  system role accessible apis
// @Tags     system
// @Produce  json
// @Param    token  header  string           false  "jwt token"
// @Param    id     path    int              true   "role id"
// @Success  200    object  ctx.StdResponse  "role detail with accessible apis"
// @Router   /api/v1/roles/{id}/apis [GET]
func roleDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.FailedRsp(c, err)
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
		ctx.FailedRsp(c, err)
		return
	}
	if err := role.ApiDetails(); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, role, nil)
}

// @Summary  del role
// @Tags     system
// @Produce  json
// @Param    token  header  string           false  "jwt token"
// @Param    id     path    int              true   "role id"
// @Success  200    object  ctx.StdResponse  "role id"
// @Router   /api/v1/roles/{id} [DELETE]
func delRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	if err := models.DelSysRole(uint(id)); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, id, nil)
}

// @Summary  add system role
// @Tags     system
// @Produce  json
// @Param    token    header  string           false  "jwt token"
// @Param    request  body    models.SysRole   true   "role request"
// @Success  200      object  ctx.StdResponse  "role id"
// @Router   /api/v1/roles [POST]
func addRole(c *gin.Context) {
	// name uniq_key
	req := &models.SysRole{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	if err := models.TemplateCreate(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, req.Id, nil)
}

type updateRoleReq struct {
	Id     uint    `json:"id" binding:"required,gt=0"` // role id
	Name   *string `json:"name,omitempty"`             // role en name
	NameCn *string `json:"name_cn,omitempty"`          // role cn name
	ApiIds *[]uint `json:"api_ids,omitempty"`          // role accessible apis id
}

// @Summary  update system role
// @Tags     system
// @Produce  json
// @Param    token    header  string           false  "jwt token"
// @Param    request  body    updateRoleReq    true   "role request"
// @Success  200      object  ctx.StdResponse  "role id"
// @Router   /api/v1/roles [PUT]
func updateRole(c *gin.Context) {
	req := &updateRoleReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	role := &models.SysRole{
		BaseModel: models.BaseModel{Id: req.Id},
	}
	fields := []string{}
	if req.Name != nil && len(*req.Name) > 0 {
		role.Name = *req.Name
		fields = append(fields, "name")
	}
	if req.NameCn != nil && len(*req.NameCn) > 0 {
		role.NameCn = *req.NameCn
		fields = append(fields, "name_cn")
	}
	if req.ApiIds != nil && len(*req.ApiIds) > 0 {
		role.ApiIds = *req.ApiIds
		fields = append(fields, "api_ids")
	}
	if err := models.TemplateUpdate(role, fields); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, req.Id, nil)
}

// @Summary  list system user
// @Tags     system
// @Produce  json
// @Param    token      header  string           false  "jwt token"
// @Param    page       query   int              false  "page, min=1"
// @Param    page_size  query   int              false  "page size, min=10, max=1000"
// @Param    role_id    query   int              false  "user role's id"
// @Param    name_cn    query   string           false  "user chinese name"
// @Param    active     query   bool             false  "system role activated"
// @Success  200        object  ctx.StdResponse  "users"
// @Router   /api/v1/users [GET]
func listUser(c *gin.Context) {
	var req struct {
		NameCN string `form:"name_cn"`
		Active *bool  `form:"active"`
		RoleID *uint  `form:"role_id"`
		models.PageReq
	}
	if err := c.BindQuery(&req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	filter := make(map[string]interface{})
	where := make(map[string]interface{})
	if req.Active != nil {
		filter["active"] = *req.Active
	}
	if req.RoleID != nil {
		where[models.ColumnContains("role_ids")] = *req.RoleID
	}
	if len(req.NameCN) > 0 {
		filter["name_cn"] = req.NameCN
	}
	users := []models.SysUser{}
	pageQuery := &models.DaoDBReq{
		PageReq:     models.PageReq{Page: req.Page, PageSize: req.PageSize},
		PageRsp:     models.PageRsp{},
		Dst:         &users,
		ModelFilter: filter,
		Where:       where,
		OrderBy:     "id desc",
	}
	err := models.TemPlatePageQuery(pageQuery)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, users, &pageQuery.PageRsp)
}

// @Summary  del system user
// @Tags     system
// @Produce  json
// @Param    token  header  string           false  "jwt token"
// @Param    id     path    int              true   "user id"
// @Success  200    object  ctx.StdResponse  "user id"
// @Router   /api/v1/users/{id} [DELETE]
func delUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	opt := models.DaoDBReq{
		Dst: &models.SysUser{
			BaseModel: models.BaseModel{Id: uint(id)},
		},
	}

	if err := models.TemplateSoftDelete(opt); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, id, nil)
}

type updateUserReq struct {
	Id      uint    `json:"id" binding:"required,gt=0"` // user id
	Email   *string `json:"email,omitempty"`            // user email address
	Active  *bool   `json:"active,omitempty"`           // user is active
	RoleIds *[]uint `json:"role_ids,omitempty"`         // user roles
}

// @Summary  update system role
// @Tags     system
// @Produce  json
// @Param    token    header  string           false  "jwt token"
// @Param    request  body    updateUserReq    true   "update user request"
// @Success  200      object  ctx.StdResponse  "user id"
// @Router   /api/v1/users [PUT]
func updateUser(c *gin.Context) {
	req := &updateUserReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	user := &models.SysUser{
		BaseModel: models.BaseModel{Id: req.Id},
	}
	fields := []string{}
	if req.Email != nil {
		user.Email = *req.Email
		fields = append(fields, "email")
	}
	if req.Active != nil {
		user.Active = *req.Active
		fields = append(fields, "active")
	}
	if req.RoleIds != nil && len(*req.RoleIds) > 0 {
		user.RoleIds = *req.RoleIds
		fields = append(fields, "role_ids")
	}
	if err := models.TemplateUpdate(user, fields); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, req.Id, nil)
}

type newUser struct {
	Name     string               `json:"name" binding:"min=1,max=20"`     // user name
	Password string               `json:"password" binding:"min=8,max=20"` // user password, sha256 encode save
	NameCn   string               `json:"name_cn" binding:"max=20"`        // user cn name
	Email    string               `json:"email" binding:"min=5"`           // user email address
	RoleIds  models.MySlice[uint] `json:"role_ids" binding:"min=1"`        // user roles
}

// @Summary  add system user
// @Tags     system
// @Produce  json
// @Param    token    header  string           false  "jwt token"
// @Param    request  body    newUser          true   "user request"
// @Success  200      object  ctx.StdResponse  "user id"
// @Router   /api/v1/users [POST]
func addUser(c *gin.Context) {
	// name uniq_key
	req := &newUser{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	newUser := &models.SysUser{
		Name:     req.Name,
		NameCn:   req.NameCn,
		Password: req.Password,
		Email:    req.Email,
		RoleIds:  req.RoleIds,
		Active:   true,
	}
	if err := models.TemplateCreate(newUser); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SetSensitiveApi(c)
	ctx.SucceedRsp(c, newUser.Id, nil)
}

type userLoginInfo struct {
	Name     string `json:"name" binding:"min=1,max=20"`     // user name
	Password string `json:"password" binding:"min=8,max=20"` // user password, sha256 encode save
}

// @Summary  user login
// @Tags     system
// @Produce  json
// @Param    request  body    userLoginInfo    true  "user request"
// @Success  200      object  ctx.StdResponse  "user id"
// @Router   /api/v1/users/login [POST]
func userLogin(c *gin.Context) {
	// name uniq_key
	req := &userLoginInfo{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	user := &models.SysUser{
		Name:     req.Name,
		Password: req.Password,
	}
	user.Sha256Password()
	dao := &models.DaoDBReq{
		Dst:         user,
		ModelFilter: user,
	}
	if err := models.TemplateQuery(dao); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SetSensitiveApi(c)
	if !user.Active {
		ctx.FailedRsp(c, fmt.Errorf("user is inactive"))
		return
	}
	jwtToken, _ := jwt.GenJwtToken(user.Name)
	ctx.SucceedRsp(c, map[string]interface{}{
		"name":      user.Name,
		"jwt_token": jwtToken,
	}, nil)
	ctx.SetLoginUsername(c, user.Name)
	go func() {
		models.UpdateUserLoginInfo(user.Name)
	}()
}

// @Summary  list system audit logs
// @Tags     system
// @Produce  json
// @Param    token       header  string           false  "jwt token"
// @Param    page        query   int              false  "min=1"
// @Param    page_size   query   int              false  "min=10, max=1000"
// @Param    user_name   query   string           false  "user name"
// @Param    request_id  query   string           false  "request uid"
// @Param    client_ip   query   string           false  "remote ip"
// @Param    start_time  query   string           false  "2006-01-02 15:04:05"
// @Param    end_time    query   string           false  "2006-01-02 15:04:05"
// @Success  200         object  ctx.StdResponse  "audit logs"
// @Router   /api/v1/audit_logs [GET]
func listAuditLogs(c *gin.Context) {
	var req struct {
		UserName  string `form:"user_name"`
		RequestID string `form:"request_id"`
		ClientIP  string `form:"client_ip"`
		StartTime string `form:"start_time"`
		EndTime   string `form:"end_time"`
		models.PageReq
	}
	if err := c.BindQuery(&req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}

	logs := []models.AuditLog{}
	filter := make(map[string]interface{})
	where := make(map[string]interface{})
	if len(req.UserName) > 0 {
		filter["user_name"] = req.UserName
	}
	if len(req.RequestID) > 0 {
		filter["request_id"] = req.RequestID
	}
	if len(req.ClientIP) > 0 {
		filter["client_ip"] = req.ClientIP
	}
	if len(req.StartTime) > 0 {
		where["created_at >= ?"] = req.StartTime
	}
	if len(req.EndTime) > 0 {
		where["created_at <= ?"] = req.EndTime
	}

	pageQuery := &models.DaoDBReq{
		PageReq:     models.PageReq{Page: req.Page, PageSize: req.PageSize},
		PageRsp:     models.PageRsp{},
		Dst:         &logs,
		Where:       where,
		ModelFilter: filter,
		OrderBy:     "id desc",
	}
	err := models.TemPlatePageQuery(pageQuery)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, logs, &pageQuery.PageRsp)
}

func LoadSysApis(r *gin.Engine) {
	apis := []models.Api{
		// system api manage
		{Path: "/apis", Method: http.MethodGet, Description: "列举服务本身所有api接口", Handler: listApis},
		{Path: "/apis", Method: http.MethodPut, Description: "更新相关api接口", Handler: updateApi},

		// system role manage
		{Path: "/roles/:id/apis", Method: http.MethodGet, Description: "获取角色api权限详情", Handler: roleDetail},
		{Path: "/roles/:id", Method: http.MethodDelete, Description: "删除角色", Handler: delRole},
		{Path: "/roles", Method: http.MethodGet, Description: "列举系统角色", Handler: listSysRoles},
		{Path: "/roles", Method: http.MethodPost, Description: "新增角色", Handler: addRole},
		{Path: "/roles", Method: http.MethodPut, Description: "更新角色", Handler: updateRole},

		// system user manage
		{Path: "/users/:id", Method: http.MethodDelete, Description: "删除用户", Handler: delUser},
		{Path: "/users", Method: http.MethodGet, Description: "列举系统用户", Handler: listUser},
		{Path: "/users", Method: http.MethodPost, Description: "新增用户", Handler: addUser},
		{Path: "/users", Method: http.MethodPut, Description: "更新用户", Handler: updateUser},
		{Path: "/users/login", Method: http.MethodPost, Description: "用户登录", Handler: userLogin},

		// system audit log
		{Path: "/audit_logs", Method: http.MethodGet, Description: "接口审计日志查询", Handler: listAuditLogs},
	}
	loadApi(r, ginGroupApiV1, apis)
}
