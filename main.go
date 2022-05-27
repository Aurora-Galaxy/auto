package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"time"
)

const (
	username = "18098795753"
	password = "112233"
)


func main(){
	post()
}



func Request1(req *http.Request) {
	req.Header.Set("Host", "student.wozaixiaoyuan.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-us,en")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat")
	req.Header.Set("Referer", "https://servicewechat.com/wxce6d08f781975d91/183/page-frame.html")
	req.Header.Set("Content-Length", "360")
}

// Request2 登录成功后返回用户的JWSESSION,添加入请求头，提交数据
func Request2(req *http.Request, jw string) {
	req.Header.Set("Host", "student.wozaixiaoyuan.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-us,en")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat")
	req.Header.Set("Referer", "https://servicewechat.com/wxce6d08f781975d91/183/page-frame.html")
	req.Header.Set("Content-Length", "360")
	req.Header.Set("JWSESSION",jw)
}

func sendMail(text []byte) {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "1302997173@qq.com"
	// 设置接收方的邮箱
	e.To = []string{"2508571934@qq.com"}
	//设置主题
	e.Subject = "我在校园自动打卡"
	//设置文件发送的内容
	e.Text = text
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("",
		"1302997173@qq.com", "kmibjpfvgamxfggh", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
	//kmibjpfvgamxfggh
}

func post(){
	jw := login()
	client := http.Client{}
	//获取分页内容，URL传参
	api := "https://student.wozaixiaoyuan.com/health/save.json"
	//传入的数据为x-www-form-urlencoded
	data := url.Values{
		"answers": {"[\"0\",\"1\",\"36.6\"]"},
		"areacode": {"610116"},
		"city": {"西安市"},
		"country": {"中国"},
		"district": {"长安区"},
		"latitude": {"34.15775"},
		"longitude": {"108.90688"},
		"province": {"陕西省"},
		"street": {"西长安街"},
		"township": {"韦曲街道"},
		"timestampHeader":{strconv.FormatInt(time.Now().Unix()*1000, 10)},
		//"citycode": {"156610100"},
		//"signatureHeader": {"cbdd9c883eb1f3c61659136fe9060666219c335c2d73cb782255bb01b9722e03"},
	}
	form := data.Encode() //将map转换为x-www-form-urlencoded
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(form))
	//添加请求头
	Request2(req,jw)
	//Do方法发送请求，返回HTTP回复
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求错误 err = ", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var code int
	json.Unmarshal(body,code)
	if code == 0{
		sendMail([]byte("打卡成功"))
	}else {
		sendMail([]byte("打卡失败"))
	}
}

func login() string {
	//1.发送请求
	client := http.Client{}
	//获取分页内容，URL传参
	loginUrl := "https://gw.wozaixiaoyuan.com/basicinfo/mobile/login/username"
	url := loginUrl + "?username=" + username + "&password=" + password
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("请求错误 err = ", err)
	}
	//添加请求头
	Request1(req)
	//Do方法发送请求，返回HTTP回复
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求错误 err = ", err)
	}
	//fmt.Println(resp)
	jw := resp.Header.Get(`Jwsession`)
	//fmt.Println(jw)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var code int
	json.Unmarshal(body,code)
	if code == 0{
		fmt.Println("登录成功")
		//post()
	}else {
		fmt.Println("登录失败")
	}
	return jw
	//fmt.Println("response Body:", string(body))
}


//var l list.List
//l.PushBack("0")
//l.PushBack("1")
//l.PushBack("36.6")
//data := &Data{
//	Answers:   answer,
//	Areacode:  "610116",
//	City:      "西安市",
//	Country:   "中国",
//	District:  "长安区",
//	Latitude:  "34.15775",
//	Longitude: "108.90688",
//	Province:  "陕西省",
//	Street:    "西长安街",
//	Township:  "韦曲街道",
//}
//jsonStr , err := json.Marshal(data)
//if err != nil{
//	panic(err)
//}
//	["0","1","36.6"]

//jsonStr := []byte(
//	`{"answers": answer,
//	"areacode": "610116",
//	"city":      "西安市",
//	"country":   "中国",
//	"district":  "长安区",
//	"latitude":  "34.15775",
//	"longitude": "108.90688",
//	"province":  "陕西省",
//	"street":    "西长安街",
//	"township":  "韦曲街道"}`)
//resp, err := http.PostForm(api,
//type Data struct {
//	Answers   Answers `json:"answers"from:"answers"`
//	Areacode  string   `json:"areacode"from:"areacode"`
//	City      string   `json:"city"from:"city"`
//	Country   string   `json:"country"from:"country"`
//	District  string   `json:"district"from:"district"`
//	Latitude  string   `json:"latitude"from:"latitude"`
//	Longitude string   `json:"longitude"from:"longitude"`
//	Province  string   `json:"province"from:"province"`
//	Street    string   `json:"street"from:"street"`
//	Township  string   `json:"township"from:"township"`
//}
//
//type Answers struct {
//	res string
//	ans string
//	tem string
//}