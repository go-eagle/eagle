package utils

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//Pagination 分页器
type Pagination struct {
	Request  *http.Request
	Total    int
	PageSize int
}

//NewPagination 新建分页器
func NewPagination(req *http.Request, total int, pageSize int) *Pagination {
	return &Pagination{
		Request:  req,
		Total:    total,
		PageSize: pageSize,
	}
}

//Pages 渲染生成html分页标签
func (p *Pagination) Pages() string {
	queryParams := p.Request.URL.Query()
	//从当前请求中获取page
	page := queryParams.Get("page")
	if page == "" {
		page = "1"
	}
	//将页码转换成整型，以便计算
	pageSize, _ := strconv.Atoi(page)
	if pageSize == 0 {
		return ""
	}

	//计算总页数
	var totalPageNum = int(math.Ceil(float64(p.Total) / float64(p.PageSize)))

	//首页链接
	var firstLink string
	//上一页链接
	var prevLink string
	//下一页链接
	var nextLink string
	//末页链接
	var lastLink string
	//中间页码链接
	var pageLinks []string

	//首页和上一页链接
	if pageSize > 1 {
		firstLink = fmt.Sprintf(`<li><a class="pagination-link" href="%s">首页</a></li>`, p.pageURL("1"))
		prevLink = fmt.Sprintf(`<li><a class="pagination-link" href="%s">上一页</a></li>`, p.pageURL(strconv.Itoa(pageSize-1)))
	} else {
		firstLink = `<li><a class="pagination-link" disabled href="javascript:void(0);">首页</a></li>`
		prevLink = `<li><a class="pagination-link" disabled href="javascript:void(0);">上一页</a></li>`
	}

	//末页和下一页
	if pageSize < totalPageNum {
		lastLink = fmt.Sprintf(`<li><a class="pagination-link" href="%s">末页</a></li>`, p.pageURL(strconv.Itoa(totalPageNum)))
		nextLink = fmt.Sprintf(`<li><a class="pagination-link" href="%s">下一页</a></li>`, p.pageURL(strconv.Itoa(pageSize+1)))
	} else {
		lastLink = `<li><a class="pagination-link" disabled href="javascript:void(0);">末页</a></li>`
		nextLink = `<li><a class="pagination-link" disabled href="javascript:void(0);">下一页</a></li>`
	}

	//生成中间页码链接
	pageLinks = make([]string, 0, 10)
	startPos := pageSize - 3
	endPos := pageSize + 3
	if startPos < 1 {
		endPos = endPos + int(math.Abs(float64(startPos))) + 1
		startPos = 1
	}
	if endPos > totalPageNum {
		endPos = totalPageNum
	}
	for i := startPos; i <= endPos; i++ {
		var s string
		if i == pageSize {
			s = fmt.Sprintf(`<li><a class="pagination-link is-current" href="%s">%d</a></li>`, p.pageURL(strconv.Itoa(i)), i)
		} else {
			s = fmt.Sprintf(`<li><a class="pagination-link" href="%s">%d</a></li>`, p.pageURL(strconv.Itoa(i)), i)
		}
		pageLinks = append(pageLinks, s)
	}

	return fmt.Sprintf(`<ul class="pagination-list">%s%s%s%s%s</ul>`,
		firstLink, prevLink, strings.Join(pageLinks, ""), nextLink, lastLink)
}

//pageURL 生成分页url
func (p *Pagination) pageURL(page string) string {
	//基于当前url新建一个url对象
	u, _ := url.Parse(p.Request.URL.String())
	q := u.Query()
	q.Set("page", page)
	u.RawQuery = q.Encode()
	return u.String()
}
