package fuDaSpider

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"spiderGo/logger"
	"strings"
	"sync"
)

/*
1. 爬取**福大要文**（https://news.fzu.edu.cn/fdyw.htm）

**要求：**
- 包含发布时间，作者，标题，阅读数以及正文。
- 可自动翻页。
- 范围：2020年1月1号 - 2021年9月1号（不要爬太多了）。
*/

const (
	homeUrl     = "https://news.fzu.edu.cn/fdyw.htm"
	homePageNum = 1
	detailPrefix = "https://news.fzu.edu.cn/"
)

var (
	f              *fuDaSpider
	fuDaSpiderOnce sync.Once
)

type fuDaSpider struct {
	homeUrl     string
	homePageNum int
}

func (f *fuDaSpider) Spider() (res [][]string, err error) {
	res = make([][]string, 0)
	for i := 0; i < f.homePageNum; i++ {
		r, err := getPageValidInfo(f.homeUrl, i)
		if err != nil {
			return res, err
		}
		res=append(res, r...)
	}
	logger.Logger.Println("res List:", res)
	return
}

func getPageValidInfo(u string, p int) (rList [][]string, err error) {
	var resString string
	resString, err = getPageInfo(u,p)
	if err != nil {
		err = fmt.Errorf("获取整页信息错误: %w", err)
		return
	}
	rList, err = getValidInfo(resString)
	if err != nil {
		err = fmt.Errorf("过滤有效信息错误: %w", err)
		return
	}
	return
}

func getPageInfo(u string ,p int) (resBodyString string, err error) {
	logger.Logger.Printf("start Get page: %v", p)
	res, err := http.DefaultClient.Get(u)
	if err != nil {
		err = fmt.Errorf("请求错误: %w", err)
		return
	}
	body, err := io.ReadAll(res.Body)
	err = res.Body.Close()
	if err != nil {
		err = fmt.Errorf("关闭响应体错误: %w", err)
		return
	}
	if res.StatusCode > 299 {
		err = fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return
	}
	resBodyString = string(body)
	return
}

func getValidInfo(b string) (resStringList [][]string, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(b))
	if err != nil{
		err=fmt.Errorf("解析html错误:%w", err)
		return
	}
	//   //section[@class='n_news']//li//div[@class='time fl']/h3
	list, err := htmlquery.QueryAll(doc, "//section[@class='n_news']//li/a")
	if err != nil {
		err=fmt.Errorf("查询标签section错误:%w", err)
		return
	}
	for _, n := range list {
		info := make([]string,0)
		detailUrl:= detailPrefix+"/"+ htmlquery.SelectAttr(n, "href")
		detailInfo, err := getPageInfo(detailUrl,0)
		if err != nil{
			err=fmt.Errorf("获取详情页:%s,信息错误:%w", detailUrl, err)
			return resStringList, err
		}
		detailDoc, err := htmlquery.Parse(strings.NewReader(detailInfo))
		if err != nil{
			err=fmt.Errorf("详情页解析html错误:%w", err)
			return resStringList, err
		}
		dateTime, err := htmlquery.Query(detailDoc,"//div[@class='ar_article_box']/div[1]/h6/span[1]")
		if err != nil {
			err=fmt.Errorf("详情页查询发布日期错误:%w", err)
			return resStringList, err
		}
		info = append(info, htmlquery.InnerText(dateTime))
		logger.Logger.Println(info)

		author, err := htmlquery.Query(detailDoc,"//div[@class='ar_article_box']/div[1]/h6/span[2]")
		if err != nil {
			err=fmt.Errorf("详情页查询作者错误:%w", err)
			return resStringList,err
		}
		info = append(info, htmlquery.InnerText(author))
		logger.Logger.Println(info)

		title, err := htmlquery.Query(detailDoc,"//div[@class='ar_article_box']/div[1]/h3")
		if err != nil {
			err=fmt.Errorf("详情页查询标题错误:%w", err)
			return resStringList, err
		}
		info = append(info, htmlquery.InnerText(title))
		logger.Logger.Println(info)
		resStringList=append(resStringList, info)
	}
	return
}

func newFuDaSpider() *fuDaSpider {
	return &fuDaSpider{
		homeUrl:     homeUrl,
		homePageNum: homePageNum}
}

type FuDaSpider interface {
	Spider() ([][]string, error)
}

func GetFuDaSpider() FuDaSpider {
	fuDaSpiderOnce.Do(func() {
		f = newFuDaSpider()
	})
	return f
}
