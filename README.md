# auto
我在校园自动打卡（go语言）
# 我在校园自动打卡（go语言版）

简介：使用go语言实现我在校园打卡，将代码放置到服务器或云服务商云函数上实现每天自动打卡。

## 一.抓取我在校园发送的数据包

使用Fiddler抓包软件抓取我在校园数据包，同时获取post数据时的URL。
    数据包含的字段如图所示
    ![在这里插入图片描述](https://img-blog.csdnimg.cn/b1f9c582ff4c4b1f896317798a32f1d9.png)
具体过程可以搜索Fiddler抓包教程

### 抓取不到我在校园数据包的具体解决方法

 1. Fiddler安装完成后，同时在电脑浏览器上安装完证书后发现PC端微信我在校园小程序打不开

    + 原因：证书安装有问题

    + 解决办法：win+R输入以下命令

      ![image-20220524234409635](https://img-blog.csdnimg.cn/img_convert/4eb738363832c0ea7b62908b1c757659.png)

      然后就会打开以下界面

      ![image-20220524234638924](https://img-blog.csdnimg.cn/img_convert/5d07bb61d5da622669a18ce4ea6ab39a.png)

      

      将个人证书中以下颁发者颁发的证书全部删除，然后重启电脑，在安装所需要安装的证书

      ![image-20220524234535640](https://img-blog.csdnimg.cn/img_convert/2d0f1def9412108b022f92480dd6c795.png)

 2. 之前抓取过微信小程序的包，可能第二次就会抓取不到

    + 解决办法：

     随便打开一个微信小程序，然后打开电脑任务管理器，找到如图所示的内容

    ![image-20220526130921964](https://img-blog.csdnimg.cn/img_convert/cacc7785df7e5e0b9fc81b5b8335bf28.png)

    然后右键打开文件所在位置，找到如下文件夹
    ![在这里插入图片描述](https://img-blog.csdnimg.cn/4c534c9639e440ea97529ebae2e795e9.png)
    将改文件夹全部删除，在打开Fiddler进行抓包



## 二.代码具体实现

1. 先定义常量，存入用户名（手机号）和密码

   ```go
   const (
      username = "xxxxxx" //填入手机号
      password = "xxxxxx" //填入密码
   )
   ```

2. 在登录之前先写一个函数去添加每次发送登录请求的请求头

   按照抓包的内容将请求头相应字段和其所对应值的加入每次的请求中

   ```go
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
   ```

3. 写登录函数（主要是要获取Jwsession，因为每次的登录Jwsession会变化）

   ```go
   func login() string {
      //1.发送请求
      client := http.Client{}
      //获取分页内容，URL传参
      loginUrl := "https://gw.wozaixiaoyuan.com/basicinfo/mobile/login/username"
       //拼接URL，加入用户名和密码
      url := loginUrl + "?username=" + username + "&password=" + password
      req, err := http.NewRequest("POST", url, nil)
      if err != nil {
         fmt.Println("请求错误 err = ", err)
      }
      //添加请求头
      Request1(req) //本文第二步所写函数
      //Do方法发送请求，返回HTTP回复
      resp, err := client.Do(req)
      if err != nil {
         fmt.Println("发送请求错误 err = ", err)
      }
       //将响应信息中的Jwsession对应的值取出
      jw := resp.Header.Get(`Jwsession`)
      defer resp.Body.Close()
      body, _ := ioutil.ReadAll(resp.Body)
      var code int
      json.Unmarshal(body,code)
      if code == 0{
         fmt.Println("登录成功")
      }else {
         fmt.Println("登录失败")
      }
      return jw  //返回Jwsession对应的值，为提交数据做准备
   }
   ```

4. 在写一个请求函数将Jwsession参数加入进去（和第一个请求函数差不多，只是加了Jwsession）

   ```go
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
   ```

5. 写post函数将打卡数据上传至固定的URL（https://student.wozaixiaoyuan.com/health/save.json）

   ```go
   func post(){
      jw := login()
      client := http.Client{}
      //获取分页内容，URL传参
      api := "https://student.wozaixiaoyuan.com/health/save.json"
      //传入的数据为x-www-form-urlencoded
      data := url.Values{
         "answers": {"[\"0\",\"1\",\"36.6\"]"},
          //以下参数自己修改
         "areacode": {"xxxx"},
         "city": {"xxxx"},
         "country": {"xxxx"},
         "district": {"xxxx"},
         "latitude": {"xxxx"},
         "longitude": {"xxxx"},
         "province": {"xxxx"},
         "street": {"xxxx"},
         "township": {"xxxx"},
         //我在校园更新后，提交的内容要加上时间戳
         "timestampHeader":{strconv.FormatInt(time.Now().Unix()*1000, 10)},
         //这两个可加可不加
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
      //发送邮件函数
         sendMail([]byte("打卡成功"))
      }else {
         sendMail([]byte("打卡失败"))
      }
   }
   ```
6. 发送邮件告知打卡成功与失败

```go
func sendMail(text []byte) {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "xxxxx"
	// 设置接收方的邮箱
	e.To = []string{"xxxxxxx"}
	//设置主题
	e.Subject = "我在校园自动打卡" //主题可以自己更改
	//设置文件发送的内容
	e.Text = text
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("",
		"xxxxx", "填入邮箱授权码", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}
```
获取授权码
QQ邮箱-->设置-->账户-->找到以下内容开启
![在这里插入图片描述](https://img-blog.csdnimg.cn/0c77b7ff10634712ac5099c3349d7fd3.png)
生成授权码
![在这里插入图片描述](https://img-blog.csdnimg.cn/d0f89f15edfc4fdeace7405a047b40b5.png)
授权码填入以下位置
![在这里插入图片描述](https://img-blog.csdnimg.cn/1c2f75c0d99a43a49d3a94bc05244f41.png)
## 三.代码部署运行
1.使用云服务商的云函数，将代码放入然后设置定时执行
2.购买一个服务器，定时执行任务
此处以第二种方法为例
1. 在windows下将代码编译成Linux下的可执行文件（如果Linux下配置好go环境的话，可以将代码直接放上去运行）
 + 在编译器终端输入以下命令
 如果是cmd终端，输入以下命令
```go
SET CGO_ENABLED=0  // 禁用CGO
SET GOOS=linux  // 目标平台是linux
SET GOARCH=amd64  // 目标处理器架构是amd64
```
如果是PowerShell终端，输入以下命令
   

```go
$ENV:CGO_ENABLED=0
$ENV:GOOS="linux"
$ENV:GOARCH="amd64"
```

2. 将可执行文件放入Linux，设置定时任务（crontab）


