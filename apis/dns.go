package apis

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/tswcbyy1107/intranet-dns/ctx"
	"github.com/tswcbyy1107/intranet-dns/models"
	"github.com/tswcbyy1107/intranet-dns/service/dnslib"
)

type newZone struct {
	Zone        string `json:"zone,omitempty" binding:"gt=0"`        // ns name zone name FQDN
	NsAddress   string `json:"ns_address,omitempty" binding:"gt=0"`  // ns server ip:port
	TsigName    string `json:"tsig_name,omitempty" binding:"gt=0"`   // tsig key name
	TsigSecret  string `json:"tsig_secret,omitempty" binding:"gt=0"` // ns dynamic update key, to be encoded
	Description string `json:"description,omitempty" binding:"gt=0"` // zone description
}

// @Summary  add intranet dns zone
// @Tags     dns
// @Produce  json
// @Param    token    header  string               false  "min=1"
// @Param    request  body    newZone          true   "zone request"
// @Success  200      object  ctx.StdResponse      "user id"
// @Router   /api/v1/dns/zones [POST]
func addZone(c *gin.Context) {
	req := &newZone{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	username := ctx.GetLoginUsername(c)

	newZone := &models.DnsZone{
		Zone:        req.Zone,
		NsAddress:   req.NsAddress,
		TsigName:    req.TsigName,
		TsigSecret:  req.TsigSecret,
		Description: req.Description,
		Creator:     username,
	}

	if err := newZone.PreCheck(); err != nil {
		ctx.FailedRsp(c, err)
		return
	}

	if err := models.TemplateCreate(newZone); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SetSensitiveApi(c)
	ctx.SucceedRsp(c, newZone.Id, nil)
}

// @Summary  del intranet dns zone
// @Tags     dns
// @Produce  json
// @Param    token   header  string           false  "min=1"
// @Param    id     path    int              true   "user id"
// @Success  200    object  ctx.StdResponse  "user id"
// @Router   /api/v1/dns/zones/{id} [DELETE]
func delZone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	opt := models.DaoDBReq{
		Dst: &models.DnsZone{
			BaseModel: models.BaseModel{Id: uint(id)},
		},
	}

	if err := models.TemplateSoftDelete(opt); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, id, nil)
}

type updateZoneReq struct {
	Id          uint   `json:"id" binding:"required,gt=0"`
	Zone        string `json:"zone,omitempty"`
	NsAddress   string `json:"ns_address,omitempty"`
	TsigName    string `json:"tsig_name,omitempty"`
	TsigSecret  string `json:"tsig_secret,omitempty"`
	Description string `json:"description,omitempty"`
}

// @Summary  update intranet dns zone
// @Tags     dns
// @Produce  json
// @Param    token    header  string           false  "min=1"
// @Param    request  body    updateZoneReq    true   "update zone request"
// @Success  200      object  ctx.StdResponse  "user id"
// @Router   /api/v1/dns/zones [PUT]
func updateZone(c *gin.Context) {
	req := &updateZoneReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	zone := &models.DnsZone{
		BaseModel: models.BaseModel{Id: req.Id},
	}
	fields := []string{}
	if len(req.Zone) > 0 {
		zone.Zone = dns.Fqdn(req.Zone)
		fields = append(fields, "zone")
	}
	if len(req.NsAddress) > 0 {
		zone.NsAddress = req.NsAddress
		fields = append(fields, "ns_address")
	}
	if len(req.TsigName) > 0 {
		zone.TsigName = dns.Fqdn(req.TsigName)
		fields = append(fields, "tsig_name")
	}
	if len(req.TsigSecret) > 0 {
		zone.TsigSecret = req.TsigSecret
		fields = append(fields, "tsig_secret")
	}
	if len(req.Description) > 0 {
		zone.Description = req.Description
		fields = append(fields, "description")
	}
	if err := zone.PreCheck(); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	if err := models.TemplateUpdate(zone, fields); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, req.Id, nil)
}

// @Summary  list intranet dns zone
// @Tags     dns
// @Produce  json
// @Param    token  header  string           false  "min=1"
// @Param    page            query   int              false  "min=1"
// @Param    page_size       query   int              false  "min=10, max=1000"
// @Success  200     object  ctx.StdResponse  "roles"
// @Router   /api/v1/dns/zones [GET]
func listDnsZone(c *gin.Context) {
	var req models.PageReq
	if err := c.BindQuery(&req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}

	zones := []models.DnsZone{}
	pageQuery := &models.DaoDBReq{
		PageReq: models.PageReq{Page: req.Page, PageSize: req.PageSize},
		PageRsp: models.PageRsp{},
		Dst:     &zones,
		OrderBy: "id desc",
	}
	err := models.TemPlatePageQuery(pageQuery)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, zones, &pageQuery.PageRsp)
}

// @Summary  list intranet dns
// @Tags     dns
// @Produce  json
// @Param    token           header  string           false  "min=1"
// @Param    page       query   int              false  "min=1"
// @Param    page_size  query   int              false  "min=10, max=1000"
// @Param    record_name     query   string           false  "domain"
// @Param    record_type     query   string           false  "type"
// @Param    record_content  query   string           false  "data"
// @Success  200             object  ctx.StdResponse  "roles"
// @Router   /api/v1/dns/records [GET]
func listDns(c *gin.Context) {
	type r struct {
		RecordName    string `form:"record_name" json:"record_name,omitempty"`
		Zone          string `form:"zone" json:"zone,omitempty"`
		RecordType    string `form:"record_type" json:"record_type,omitempty"`
		RecordContent string `form:"record_content" json:"record_content,omitempty"`
		models.PageReq
	}
	req := r{}
	if err := c.BindQuery(&req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}

	filter := make(map[string]interface{})
	if len(req.RecordName) > 0 {
		filter["record_name"] = dns.Fqdn(req.RecordName)
	}
	if len(req.Zone) > 0 {
		filter["zone"] = dns.Fqdn(req.Zone)
	}
	if len(req.RecordType) > 0 {
		filter["record_type"] = strings.ToUpper(req.RecordType)
	}
	if len(req.RecordContent) > 0 {
		if req.RecordType == models.CnameType {
			req.RecordContent = dns.Fqdn(req.RecordContent)
		}
		filter["record_content"] = req.RecordContent
	}

	records := []models.DnsRecord{}
	pageQuery := &models.DaoDBReq{
		PageReq:     models.PageReq{Page: req.Page, PageSize: req.PageSize},
		PageRsp:     models.PageRsp{},
		Dst:         &records,
		ModelFilter: filter,
		OrderBy:     "id desc",
	}
	err := models.TemPlatePageQuery(pageQuery)
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, records, &pageQuery.PageRsp)
}

// @Summary  del intranet dns record
// @Tags     dns
// @Produce  json
// @Param    token   header  string           false  "min=1"
// @Param    id     query   int              true   "record id"
// @Param    clean  query   bool             false  "del all of the same type rrs"
// @Success  200    object  ctx.StdResponse  "user id"
// @Router   /api/v1/dns/records [DELETE]
func delDns(c *gin.Context) {
	req := struct {
		Id    uint `json:"id,omitempty" form:"id" binding:"required,gt=0"`
		Clean bool `json:"clean,omitempty" form:"clean"`
	}{}
	if err := c.BindQuery(&req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	if err := dnslib.DelIntranetDns(req.Id, req.Clean); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, req.Id, nil)
}

// @Summary  update intranet dns
// @Tags     dns
// @Produce  json
// @Param    token    header  string           false  "min=1"
// @Param    request  body    dnslib.UpdateDnsReq  true   "update dns request"
// @Success  200      object  ctx.StdResponse  "user id"
// @Router   /api/v1/dns/records [PUT]
func updateDns(c *gin.Context) {
	req := &dnslib.UpdateDnsReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	if err := dnslib.UpdateIntranetDns(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, req.Id, nil)
}

type newDns struct {
	RecordName    string `binding:"gt=0" json:"record_name,omitempty"`
	RecordTtl     int    `binding:"gt=0" json:"record_ttl,omitempty"`
	Zone          string `binding:"gt=0" json:"zone,omitempty"`
	RecordType    string `binding:"gt=0" json:"record_type,omitempty"`
	RecordContent string `binding:"gt=0" json:"record_content,omitempty"`
}

// @Summary  add intranet dns
// @Tags     dns
// @Produce  json
// @Param    token    header  string           false  "min=1"
// @Param    request  body    newDns           true   "dns request"
// @Success  200      object  ctx.StdResponse  "user id"
// @Router   /api/v1/dns/records [POST]
func addDns(c *gin.Context) {
	req := &newDns{}
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	username := ctx.GetLoginUsername(c)

	if err := dnslib.AddIntranetDns(&models.DnsRR{
		RecordName:    req.RecordName,
		RecordTtl:     req.RecordTtl,
		Zone:          req.Zone,
		RecordType:    req.RecordType,
		RecordContent: req.RecordContent,
	}, username); err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, "", nil)
}

// @Summary  intranet dns query
// @Tags     dns
// @Produce  json
// @Param    token  header  string           false  "min=1"
// @Param    domain  query   string           true   "domain"
// @Success  200     object  ctx.StdResponse  "roles"
// @Router   /api/v1/dns/rrs [GET]
func dnsQuery(c *gin.Context) {
	domain := strings.TrimSpace(c.Query("domain"))
	if len(domain) == 0 {
		ctx.FailedRsp(c, fmt.Errorf("domain requiered"))
		return
	}
	rrs, err := dnslib.IntranetRRQueryAll(domain, "test.com")
	if err != nil {
		ctx.FailedRsp(c, err)
		return
	}
	ctx.SucceedRsp(c, rrs, nil)
}

// @Summary  edns query
// @Tags     dns
// @Produce  json
// @Param    token  header  string           false  "min=1"
// @Param    domain  query   string           true   "domain"
// @Success  200        object  ctx.StdResponse  "roles"
// @Router   /api/v1/dns/edns [GET]
func ednsQuery(c *gin.Context) {
	domain := strings.TrimSpace(c.Query("domain"))
	if len(domain) == 0 {
		ctx.FailedRsp(c, fmt.Errorf("domain requiered"))
		return
	}

	ch := make(chan dnslib.EdnsRRs, len(dnslib.PublicDnsIP))
	data := []dnslib.EdnsRRs{}

	wg := &sync.WaitGroup{}
	for _, isp := range dnslib.PublicDnsIP {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rrs, _ := dnslib.PublicEDnsQueryRR(domain, isp.DnsIP)
			edns := dnslib.EdnsRRs{
				ISP:      isp.ISP,
				ClientIP: isp.DnsIP,
				DnsRRs:   rrs,
			}
			ch <- edns
		}()
	}
	wg.Wait()
	close(ch)

	for rrs := range ch {
		data = append(data, rrs)
	}
	ctx.SucceedRsp(c, data, nil)
}

// @Summary  province isp ns ip
// @Tags     dns
// @Produce  json
// @Param    token      header  string           false  "min=1"
// @Success  200    object  ctx.StdResponse  "isps"
// @Router   /api/v1/dns/isps [GET]
func getIsps(c *gin.Context) {
	ctx.SucceedRsp(c, dnslib.PublicDnsIP, nil)
}

func LoadDnsApis(r *gin.Engine) {
	apis := []models.Api{
		// intranet dns crud, A AAAA CNAME
		{Path: "/dns/zones", Method: http.MethodPost, Description: "内网dns zone新增", Handler: addZone},
		{Path: "/dns/zones/:id", Method: http.MethodDelete, Description: "内网dns zone清理", Handler: delZone},
		{Path: "/dns/zones", Method: http.MethodPut, Description: "内网dns zone更新", Handler: updateZone},
		{Path: "/dns/zones", Method: http.MethodGet, Description: "内网dns zone查询", Handler: listDnsZone},

		{Path: "/dns/records", Method: http.MethodGet, Description: "内网dns记录枚举", Handler: listDns},
		{Path: "/dns/records", Method: http.MethodDelete, Description: "内网dns记录删除", Handler: delDns},
		{Path: "/dns/records", Method: http.MethodPut, Description: "内网dns更新", Handler: updateDns},
		{Path: "/dns/records", Method: http.MethodPost, Description: "内网dns新增", Handler: addDns},
		{Path: "/dns/rrs", Method: http.MethodGet, Description: "内部dns查询", Handler: dnsQuery},

		{Path: "/dns/edns", Method: http.MethodGet, Description: "公网edns查询", Handler: ednsQuery},
		{Path: "/dns/isps", Method: http.MethodGet, Description: "各省份运营商dns地址", Handler: getIsps},

		{Path: "/dns/probe/:id", Method: http.MethodDelete, Description: "dns探测删除", Handler: delUser},
	}
	loadApi(r, ginGroupApiV1, apis)
}
