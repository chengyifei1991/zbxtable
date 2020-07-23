package controllers

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	"github.com/canghai908/zbxtable/models"
)

// ImagesController operations for Host
type ImagesController struct {
	//	BaseController
	beego.Controller
}

// URLMapping ...
func (c *ImagesController) URLMapping() {
	//c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
}

// GetOne ...
// @Title Get One
// @Description get Alarm by id
// @Param	X-Token		header  string			true		"x-token in header"
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Alarm
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ImagesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	var StartTime, EndTime string
	all := c.Ctx.Input.Query("from")
	b := strings.Split(all, "?")
	StartTime = b[0]
	EndTime = strings.Split(b[1], "=")[1]
	GraphID := idStr
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	c.Ctx.ResponseWriter.Header().Set("Pragma", "no-cache, value")

	client1 := &http.Client{nil, nil,
		models.JAR, 99999999999999}
	ZabbixServer := beego.AppConfig.String("zabbix_server")
	//imgurl
	imgurl := ZabbixServer + "/chart2.php?"
	data := url.Values{}
	URL, err := url.Parse(imgurl)
	if err != nil {
		beego.Error("Fatal error ", err.Error())
	}
	data.Set("graphid", GraphID)
	data.Set("StartTime", StartTime)
	data.Set("EndTime", EndTime)
	data.Set("profileIdx", "web.graphs.filter")
	data.Set("profileIdx2", "200")
	data.Set("width", "400")
	//Encode rul
	URL.RawQuery = data.Encode()
	urlPath := URL.String()
	reqest1, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		beego.Error("Fatal error ", err.Error())
	}
	reqest1.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest1.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest1.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest1.Header.Add("Connection", "keep-alive")
	reqest1.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	response1, err := client1.Do(reqest1)
	if err != nil {
		beego.Error("Fatal error ", err.Error())
	}
	defer response1.Body.Close()
	if response1.StatusCode == 200 {
		var reader io.Reader
		switch response1.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ = gzip.NewReader(response1.Body)
		default:
			reader = response1.Body
		}
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			beego.Error("读取响应数据失败: %+v", err)
		}
		c.Ctx.ResponseWriter.Write(data)
	}
}
