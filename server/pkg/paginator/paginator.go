package paginator

import (
	"fmt"
	"math"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Paging 分页数据
type Paging struct {
	CurrentPage int    // 当前页
	PerPage     int    // 每页条数
	TotalPage   int    // 总页数
	TotalCount  int64  // 总条数
	NextPageURL string // 下一页的链接
	PrevPageURL string // 上一页的链接
}

// Paginator 分页操作类
type Paginator struct {
	BaseURL    string // 用以拼接 URL
	PerPage    int    // 每页条数
	Page       int    // 当前页
	Offset     int    // 数据库读取数据时 Offset 的值
	TotalCount int64  // 总条数
	TotalPage  int    // 总页数 = TotalCount/PerPage
	Sort       string // 排序规则
	Order      string // 排序顺序

	query *gorm.DB     // db query 句柄
	ctx   *gin.Context // gin context，方便调用
}

// Paginate 分页
// c —— gin.context 用来获取分页的 URL 参数
// db —— GORM 查询句柄，用以查询数据集和获取数据总数
// baseURL —— 用以分页链接
// data —— 模型数组，传址获取数据
// perPage —— 每页条数，优先从 url 参数里取，否则使用 perPage 的值
// filters —— 筛选条件
// 用法:
//
// query := database.DB.Model(Topic{}).Where("category_id = ?", cid)
// var topics []Topic
// paging := paginator.Paginate(
//
//	c,
//	query,
//	&topics,
//	app.APIURL(database.TableName(&Topic{})),
//	perPage,
//	filters,
//
// )
func Paginate(c *gin.Context, db *gorm.DB, data interface{}, baseURL string, perPage int, filters map[string]interface{}) Paging {
	// 初始化 Paginator 实例
	p := &Paginator{
		query: db,
		ctx:   c,
	}
	p.initProperties(perPage, baseURL, filters)
	// 查询数据库
	err := p.query.Preload(clause.Associations). // 读取关联
							Order(p.Sort + " " + p.Order). // 排序
							Limit(p.PerPage).
							Offset(p.Offset).
							Find(data).
							Error

	// 数据库出错
	if err != nil {
		logger.LogIf(err)
		return Paging{}
	}

	return Paging{
		CurrentPage: p.Page,
		PerPage:     p.PerPage,
		TotalPage:   p.TotalPage,
		TotalCount:  p.TotalCount,
		NextPageURL: p.getNextPageURL(),
		PrevPageURL: p.getPrevPageURL(),
	}
}

func PaginateCustomQuery(c *gin.Context, query *gorm.DB, data interface{}, baseURL string, perPage int) Paging {
	// 初始化 Paginator 实例
	p := &Paginator{
		query: query,
		ctx:   c,
	}

	// 获取基础参数
	p.BaseURL = p.formatBaseURL(baseURL)
	p.PerPage = p.getPerPage(perPage)
	p.Order = c.DefaultQuery(config.Get("paging.url_query_order"), "desc")
	p.Sort = c.DefaultQuery(config.Get("paging.url_query_sort"), "id")

	// 获取总数
	p.TotalCount = p.getTotalCount()
	p.TotalPage = p.getTotalPage()
	p.Page = p.getCurrentPage()
	p.Offset = (p.Page - 1) * p.PerPage

	// 执行查询（不再使用 clause.Associations）
	err := p.query.Order(p.Sort + " " + p.Order).
		Limit(p.PerPage).
		Offset(p.Offset).
		Find(data).
		Error

	if err != nil {
		logger.LogIf(err)
		return Paging{}
	}

	return Paging{
		CurrentPage: p.Page,
		PerPage:     p.PerPage,
		TotalPage:   p.TotalPage,
		TotalCount:  p.TotalCount,
		NextPageURL: p.getNextPageURL(),
		PrevPageURL: p.getPrevPageURL(),
	}
}

// 初始化分页必须用到的属性，基于这些属性查询数据库
func (p *Paginator) initProperties(perPage int, baseURL string, filters map[string]interface{}) {

	p.BaseURL = p.formatBaseURL(baseURL)
	p.PerPage = p.getPerPage(perPage)

	// 排序参数（控制器中以验证过这些参数，可放心使用）
	p.Order = p.ctx.DefaultQuery(config.Get("paging.url_query_order"), "desc")
	p.Sort = p.ctx.DefaultQuery(config.Get("paging.url_query_sort"), "id")

	p.applyFilters(filters) // Apply filters to query

	p.TotalCount = p.getTotalCount()
	p.TotalPage = p.getTotalPage()
	p.Page = p.getCurrentPage()
	p.Offset = (p.Page - 1) * p.PerPage
}

// applyFilters applies the filters to the query
func (p *Paginator) applyFilters(filters map[string]interface{}) {
	for key, value := range filters {
		p.query = p.query.Where(key, value)
	}
}

// getPerPage returns the perPage value
func (p Paginator) getPerPage(perPage int) int {
	// 优先使用请求 per_page 参数
	queryPerpage := p.ctx.Query(config.Get("paging.url_query_per_page"))
	if len(queryPerpage) > 0 {
		perPage = cast.ToInt(queryPerpage)
	}

	// 没有传参，使用默认
	if perPage <= 0 {
		perPage = config.GetInt("paging.perpage")
	}

	return perPage
}

// getCurrentPage 返回当前页码
//func (p Paginator) getCurrentPage() int {
//	// 优先取用户请求的 page
//	page := cast.ToInt(p.ctx.Query(config.Get("paging.url_query_page")))
//	if page <= 0 {
//		// 默认为 1
//		page = 1
//	}
//	// TotalPage 等于 0 ，意味着数据不够分页
//	if p.TotalPage == 0 {
//		return 0
//	}
//	// 请求页数大于总页数，返回总页数
//	if page > p.TotalPage {
//		return p.TotalPage
//	}
//	return page
//}

func (p Paginator) getCurrentPage() int {
	page := cast.ToInt(p.ctx.Query(config.Get("paging.url_query_page")))
	if page <= 0 {
		page = 1
	}
	// 不再对 TotalPage 做特殊处理
	return page
}

// getTotalCount 返回的是数据库里的条数
func (p *Paginator) getTotalCount() int64 {
	var count int64
	if err := p.query.Count(&count).Error; err != nil {
		return 0
	}
	return count
}

// getTotalPage 计算总页数
func (p Paginator) getTotalPage() int {
	if p.TotalCount == 0 {
		return 0
	}
	nums := int64(math.Ceil(float64(p.TotalCount) / float64(p.PerPage)))
	if nums == 0 {
		nums = 1
	}
	return int(nums)
}

// 兼容 URL 带与不带 `?` 的情况
func (p *Paginator) formatBaseURL(baseURL string) string {
	if strings.Contains(baseURL, "?") {
		baseURL = baseURL + "&" + config.Get("paging.url_query_page") + "="
	} else {
		baseURL = baseURL + "?" + config.Get("paging.url_query_page") + "="
	}
	return baseURL
}

// 拼接分页链接
func (p Paginator) getPageLink(page int) string {
	return fmt.Sprintf("%v%v&%s=%s&%s=%s&%s=%v",
		p.BaseURL,
		page,
		config.Get("paging.url_query_sort"),
		p.Sort,
		config.Get("paging.url_query_order"),
		p.Order,
		config.Get("paging.url_query_per_page"),
		p.PerPage,
	)
}

// getNextPageURL 返回下一页的链接
func (p Paginator) getNextPageURL() string {
	if p.TotalPage > p.Page {
		return p.getPageLink(p.Page + 1)
	}
	return ""
}

// getPrevPageURL 返回上一页的链接
func (p Paginator) getPrevPageURL() string {
	if p.Page <= 1 || p.Page > p.TotalPage {
		return ""
	}
	return p.getPageLink(p.Page - 1)
}
