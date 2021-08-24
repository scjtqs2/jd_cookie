package web

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
)

var timeout = time.Second * 5

//var ua = "jdapp;android;10.0.5;11;0393465333165363-5333430323261366;network/wifi;model/M2102K1C;osVer/30;appBuild/88681;partner/lc001;eufv/1;jdSupportDarkMode/0;Mozilla/5.0 (Linux; Android 11; M2102K1C Build/RKQ1.201112.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045534 Mobile Safari/537.36"
var user_agents = []string{
	"jdapp;android;10.0.2;10;network/wifi;Mozilla/5.0 (Linux; Android 10; ONEPLUS A5010 Build/QKQ1.191014.012; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;iPhone;10.0.2;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.0.2;9;network/4g;Mozilla/5.0 (Linux; Android 9; Mi Note 3 Build/PKQ1.181007.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045131 Mobile Safari/537.36",
	"jdapp;android;10.0.2;10;network/wifi;Mozilla/5.0 (Linux; Android 10; GM1910 Build/QKQ1.190716.003; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;android;10.0.2;9;network/wifi;Mozilla/5.0 (Linux; Android 9; 16T Build/PKQ1.190616.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;iPhone;10.0.2;13.6;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;13.6;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;13.5;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;14.1;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;13.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;13.7;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;14.1;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;13.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;13.4;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 13_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;14.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.0.2;9;network/wifi;Mozilla/5.0 (Linux; Android 9; MI 6 Build/PKQ1.190118.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;android;10.0.2;11;network/wifi;Mozilla/5.0 (Linux; Android 11; Redmi K30 5G Build/RKQ1.200826.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045511 Mobile Safari/537.36",
	"jdapp;iPhone;10.0.2;11.4;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 11_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15F79",
	"jdapp;android;10.0.2;10;;network/wifi;Mozilla/5.0 (Linux; Android 10; M2006J10C Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;android;10.0.2;10;network/wifi;Mozilla/5.0 (Linux; Android 10; M2006J10C Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;android;10.0.2;10;network/wifi;Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190716.003; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045224 Mobile Safari/537.36",
	"jdapp;android;10.0.2;9;network/wifi;Mozilla/5.0 (Linux; Android 9; MHA-AL00 Build/HUAWEIMHA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;android;10.0.2;8.1.0;network/wifi;Mozilla/5.0 (Linux; Android 8.1.0; 16 X Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;android;10.0.2;8.0.0;network/wifi;Mozilla/5.0 (Linux; Android 8.0.0; HTC U-3w Build/OPR6.170623.013; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/044942 Mobile Safari/537.36",
	"jdapp;iPhone;10.0.2;14.0.1;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.0.2;10;network/wifi;Mozilla/5.0 (Linux; Android 10; LYA-AL00 Build/HUAWEILYA-AL00L; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045230 Mobile Safari/537.36",
	"jdapp;iPhone;10.0.2;14.2;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;14.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;14.2;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.0.2;8.1.0;network/wifi;Mozilla/5.0 (Linux; Android 8.1.0; MI 8 Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045131 Mobile Safari/537.36",
	"jdapp;android;10.0.2;10;network/wifi;Mozilla/5.0 (Linux; Android 10; Redmi K20 Pro Premium Edition Build/QKQ1.190825.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045227 Mobile Safari/537.36",
	"jdapp;iPhone;10.0.2;14.3;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;iPhone;10.0.2;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
	"jdapp;android;10.0.2;11;network/wifi;Mozilla/5.0 (Linux; Android 11; Redmi K20 Pro Premium Edition Build/RKQ1.200826.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045513 Mobile Safari/537.36",
	"jdapp;android;10.0.2;10;network/wifi;Mozilla/5.0 (Linux; Android 10; MI 8 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045227 Mobile Safari/537.36",
	"jdapp;iPhone;10.0.2;14.1;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1",
}

func (s *httpServer) getUa(c *gin.Context) string {
	//lens := len(user_agents)
	//rand.Seed(time.Now().UnixNano())
	////r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//num := rand.Intn(lens)
	//ua := user_agents[num]
	ua, err := s.GetUa(c)
	if err != nil {
		//t := strconv.FormatInt(time.Now().UnixNano()/1e3, 10)
		//User-Agent: jdapp;android;10.1.0;10;3643464346636663-1346663656937316;network/wifi;model/SM-N9600;addressid/4621427242;aid/c4d4d6f61dfce97a;oaid/;osVer/29;appBuild/89583;partner/google;eufv/1;jdSupportDarkMode/0;Mozilla/5.0 (Linux; Android 10; SM-N9600 Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/92.0.4515.159 Mobile Safari/537.36
		//ua = "jdapp;android;10.1.0;10;3643464346636663-1346663656937316;network/wifi;model/SM-N9600;addressid/4621427242;aid/c4d4d6f61dfce97a;oaid/;osVer/29;appBuild/89583;partner/google;eufv/1;jdSupportDarkMode/0;Mozilla/5.0 (Linux; Android 10; SM-N9600 Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/92.0.4515.159 Mobile Safari/537.36"
		//ua = fmt.Sprintf("jdapp;android;10.1.0;10;%s-%s;network/wifi;model/SM-N9600;addressid/4621427242;aid/c4d4d6f61dfce97a;oaid/;osVer/29;appBuild/89583;partner/google;eufv/1;jdSupportDarkMode/0;Mozilla/5.0 (Linux; Android 10; SM-N9600 Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/92.0.4515.159 Mobile Safari/537.36", t, t)
		//ua = fmt.Sprintf("Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 SP-engine/2.14.0 main/1.0 baiduboxapp/11.18.0.16 (Baidu; P2 13.3.1) NABar/0.0 TM/%s", t)
		ua = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5 UCBrowser/13.4.2.1122"
		s.SaveUa(c, ua)
	}
	return ua
}

// 获取二维码
func (s *httpServer) getQrcode(c *gin.Context) {
	s.GetclientIP(c)
	log.Warn("start get qrcode")
	_, err := s.step1(c)
	if err != nil {
		c.JSON(200, MSG{
			"err": 1,
			"msg": err,
		})
		return
	}
	qrurl, err := s.setp2(c)
	if err != nil {
		c.JSON(200, MSG{
			"err": 1,
			"msg": err,
		})
		return
	}
	log.Warnf("get qrcode url = %s", qrurl)
	c.JSON(200, MSG{
		"err":    0,
		"qrcode": qrurl,
	})
}

func (s *httpServer) praseSetCookies(c *gin.Context, rsp string, cookie *cookiejar.Jar) {
	json := gjson.Parse(rsp)
	token := s.getToken(c)
	token.Stoken = json.Get("s_token").String()
	u, _ := url.Parse("https://plogin.m.jd.com")
	a := cookie.Cookies(u)
	for _, v := range a {
		if v.Name == "guid" {
			token.Guid = v.Value
		}
		if v.Name == "lsid" {
			token.Lsid = v.Value
		}
		if v.Name == "lstoken" {
			token.Lstoken = v.Value
		}
	}
	token.Cookies = "guid=" + token.Guid + "; lang=chs; lsid=" + token.Lsid + "; lstoken=" + token.Lstoken + "; "
	s.updateToken(c, token)
	//log.Warnf("cookies=%s", cookies)
}

// 获取二维码第一步
func (s *httpServer) step1(c *gin.Context) (*cookiejar.Jar, error) {
	ip := s.GetclientIP(c)
	jar := s.getCookieJar(c)
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	getUrl := "https://plogin.m.jd.com/cgi-bin/mm/new_login_entrance?lang=chs&appid=300&returnurl=https://wq.jd.com/passport/LoginRedirect?state=" + timeStamp + "&returnurl=https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport"
	client := &http.Client{
		// example, custom cookie jar implements
		Jar: jar,
		// example, ignore self-signed certificate
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	method := "GET"
	req, err := http.NewRequest(method, getUrl, nil)

	if err != nil {
		log.Errorf("get qrcode step1 faild err=%s", err.Error())
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "zh-cn")
	req.Header.Add("Referer", "https://plogin.m.jd.com/cgi-bin/mm/new_login_entrance?lang=chs&appid=300&returnurl=https://wq.jd.com/passport/LoginRedirect?state="+timeStamp+"&returnurl=https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	req.Header.Add("User-Agent", s.getUa(c))
	req.Header.Add("Host", "plogin.m.jd.com")
	req.Header.Set("X-Forwarded-For", ip)
	req.Header.Set("Proxy-Client-IP", ip)
	req.Header.Set("WL-Proxy-Client-IP", ip)
	req.Header.Set("CLIENT-IP", ip)
	req.Header.Set("X-Requested-With", "com.jingdong.app.mall")
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("get qrcode step1 faild err=%s", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("get qrcode step1 faild err=%s", err.Error())
		return nil, err
	}
	s.praseSetCookies(c, string(body), jar)
	//log.Warnf("url=%s,cookiejar=%+v,res=%s", getUrl, jar, string(body))
	s.updateCookieJar(c, jar)
	return jar, nil
}

// 获取 二维码第二步
func (s *httpServer) setp2(c *gin.Context) (string, error) {
	ip := s.GetclientIP(c)
	token := s.getToken(c)
	jar := s.getCookieJar(c)
	if token.Cookies == "" {
		return "", errors.New("empty cookies")
	}
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	getUrl := "https://plogin.m.jd.com/cgi-bin/m/tmauthreflogurl?s_token=" + token.Stoken + "&v=" + timeStamp + "&remember=true"
	client := &http.Client{
		Jar:       jar,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	var res string
	err := gout.New(client).
		POST(getUrl).
		//Debug(true).
		SetJSON(
			gout.H{
				"lang":      "chs",
				"appid":     300,
				"returnurl": "https://wqlogin2.jd.com/passport/LoginRedirect?state=" + timeStamp + "&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action",
				"source":    "wq_passport",
			},
		).
		BindBody(&res).
		SetHeader(gout.H{
			"Connection":   "Keep-Alive",
			"Content-Type": "application/x-www-form-urlencoded; Charset=UTF-8",
			"Accept":       "application/json, text/plain, */*",
			"Cookie":       token.Cookies,
			"Referer":      "https://plogin.m.jd.com/login/login?appid=300&returnurl=https://wqlogin2.jd.com/passport/LoginRedirect?state=" + timeStamp + "&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport",
			//"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
			"User-Agent":         s.getUa(c),
			"X-Forwarded-For":    ip,
			"Proxy-Client-IP":    ip,
			"WL-Proxy-Client-IP": ip,
			"CLIENT-IP":          ip,
			"X-Requested-With":   "com.jingdong.app.mall",
		}).
		SetTimeout(timeout).
		F().Retry().Attempt(5).
		WaitTime(time.Millisecond * 500).MaxWaitTime(time.Second * 5).
		Do()
	s.updateCookieJar(c, jar)
	if err != nil {
		return "", err
	}
	resjson := gjson.Parse(res)
	token.Token = resjson.Get("token").String()
	u, _ := url.Parse("https://plogin.m.jd.com")
	for _, v := range jar.Cookies(u) {
		if v.Name == "okl_token" {
			token.Okl_token = v.Value
		}
	}
	qrUrl := "https://plogin.m.jd.com/cgi-bin/m/tmauth?appid=300&client_type=m&token=" + token.Token
	s.updateToken(c, token)
	//log.Warnf("url=%s,cookiejar=%+v,res=%s", getUrl, jar, res)
	return qrUrl, nil
}

// 获取返回的cookie信息
func (s *httpServer) getCookie(c *gin.Context) {
	//session := sessions.Default(c)
	//cookies=session.Get("cookies").(string)
	jar := s.getCookieJar(c)
	check, err := s.checkLogin(c, jar)
	if err != nil {
		c.JSON(200, MSG{
			"err": 1,
			"msg": err,
		})
		return
	}
	checkJson := gjson.Parse(check)
	if checkJson.Get("errcode").Int() == 0 {
		//获取cookie
		ucookie := s.getJdCookie(check, jar, c)
		c.JSON(200, MSG{
			"err":    0,
			"cookie": ucookie,
		})
		return
	} else {
		c.JSON(200, MSG{
			"err": checkJson.Get("errcode").Int(),
			"msg": checkJson.Get("message").String(),
		})
	}

}

// 校验登录状态
func (s *httpServer) checkLogin(c *gin.Context, jar *cookiejar.Jar) (string, error) {
	token := s.getToken(c)
	ip := s.GetclientIP(c)
	if token.Cookies == "" {
		return "", errors.New("empty cookies")
	}
	res, err := s.checklogin_1(token, jar, ip, s.getUa(c))
	s.updateCookieJar(c, jar)
	if err != nil {
		return "", err
	}
	//log.Warnf("checkLogin res=%s",res)
	return res, nil
}

func (s httpServer) checklogin_1(token *Token, jar *cookiejar.Jar, ip string, ua string) (string, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	getUrl := "https://plogin.m.jd.com/cgi-bin/m/tmauthchecktoken?&token=" + token.Token + "&ou_state=0&okl_token=" + token.Okl_token
	client := &http.Client{
		Jar:       jar,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	var res string
	err := gout.New(client).
		POST(getUrl).
		//Debug(true).
		SetWWWForm(
			gout.H{
				"lang":      "chs",
				"appid":     300,
				"returnurl": "https://wqlogin2.jd.com/passport/LoginRedirect?state=1100399130787&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action",
				"source":    "wq_passport",
			},
		).
		BindBody(&res).
		SetHeader(gout.H{
			"Referer":            "https://plogin.m.jd.com/login/login?appid=300&returnurl=https://wqlogin2.jd.com/passport/LoginRedirect?state=" + timeStamp + "&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport",
			"Cookie":             token.Cookies,
			"Connection":         "Keep-Alive",
			"Content-Type":       "application/x-www-form-urlencoded; Charset=UTF-8",
			"Accept":             "application/json, text/plain, */*",
			"User-Agent":         ua,
			"X-Forwarded-For":    ip,
			"Proxy-Client-IP":    ip,
			"WL-Proxy-Client-IP": ip,
			"CLIENT-IP":          ip,
			"X-Requested-With":   "com.jingdong.app.mall",
		}).
		SetTimeout(timeout).
		F().Retry().Attempt(5).
		WaitTime(time.Millisecond * 500).MaxWaitTime(time.Second * 5).
		Do()
	return res, err
}

// 解析用户的cookie
func (s *httpServer) getJdCookie(resp string, cookie *cookiejar.Jar, c *gin.Context) string {
	token := s.getToken(c)
	tk := s.getJdCookie_1(token, cookie)
	s.updateToken(c, tk)
	return tk.UserCookie
}

func (s *httpServer) getJdCookie_1(token *Token, cookie *cookiejar.Jar) *Token {
	u, _ := url.Parse("https://plogin.m.jd.com")
	var TrackerID, pt_key, pt_pin, pt_token, pwdt_id, s_key, s_pin = "", "", "", "", "", "", ""
	for _, v := range cookie.Cookies(u) {
		if v.Name == "TrackerID" {
			TrackerID = v.Value
		}
		if v.Name == "pt_key" {
			pt_key = v.Value
		}
		if v.Name == "pt_pin" {
			pt_pin = v.Value
		}
		if v.Name == "pt_token" {
			pt_token = v.Value
		}
		if v.Name == "pwdt_id" {
			pwdt_id = v.Value
		}
		if v.Name == "s_key" {
			s_key = v.Value
		}
		if v.Name == "s_pin" {
			s_pin = v.Value
		}
	}
	token.Cookies = "TrackerID=" + TrackerID + "; pt_key=" + pt_key + "; pt_pin=" + pt_pin + "; pt_token=" + pt_token + "; pwdt_id=" + pwdt_id + "; s_key=" + s_key + "; s_pin=" + s_pin + "; wq_skey="
	token.UserCookie = "pt_key=" + pt_key + ";pt_pin=" + pt_pin + ";"
	token.PtPin = pt_pin
	token.PtKey = pt_key
	log.Info("############  登录成功，获取到 Cookie  #############")
	log.Infof("Cookie1=%s", token.UserCookie)
	log.Info("####################################################")
	return token
}

func (s *httpServer) upsave(c *gin.Context) {
	// 清空缓存参数
	token := s.getToken(c)
	//写db
	if s.Conf.DbConf.DbEnable {
		_, err := s.cookiesRepo.UpdateCookie(token.PtPin, token.PtKey, token.UserCookie)
		if err != nil {
			log.Errorf("save cookie to db faild %s", err.Error())
		}
	}
	s.cleanSession(c)
	////发送数据给 挂机服务器
	postUrl := s.Conf.UpSave
	if postUrl != "" {
		var res MSG
		code := 0
		err := gout.POST(postUrl).
			//Debug(true).
			SetWWWForm(
				gout.H{
					"userCookie": token.UserCookie,
				},
			).
			BindJSON(&res).
			SetHeader(gout.H{
				"Connection":   "Keep-Alive",
				"Content-Type": "application/x-www-form-urlencoded; Charset=UTF-8",
				"Accept":       "application/json, text/plain, */*",
				"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
			}).
			Code(&code).
			SetTimeout(timeout).
			F().Retry().Attempt(5).
			WaitTime(time.Millisecond * 500).MaxWaitTime(time.Second * 5).
			Do()
		if err != nil || code != 200 {
			log.Errorf("upsave notify post  usercookie to %s faild", postUrl)
			c.JSON(200, MSG{
				"err":   1,
				"title": "更新到挂机服务器失败",
				"msg":   err.Error(),
			})
		} else {
			errcode := res["err"]
			if errcode == nil {
				errcode = 0
			}
			title := res["title"]
			if title == nil {
				title = "更新到挂机服务成功"
			}
			msg := res["msg"]
			if msg == nil {
				msg = "cookie= " + token.UserCookie
			}
			c.JSON(200, MSG{
				"err":   errcode,
				"title": title,
				"msg":   fmt.Sprintf("%s, cookie= %s", msg, token.UserCookie),
			})
		}
		return
	}

	c.JSON(200, MSG{
		"err":   0,
		"title": "提取cookie成功",
		"msg":   "cookie= " + token.UserCookie,
	})
}

// 获取二维码
func (s *httpServer) getQrcode_jumplogin(c *gin.Context) {
	s.GetclientIP(c)
	session := sessions.Default(c)
	log.Warn("start get qrcode")
	_, err := s.step1(c)
	if err != nil {
		c.JSON(200, MSG{
			"err": 1,
			"msg": err,
		})
		return
	}
	qrurl, err := s.setp2(c)
	if err != nil {
		c.JSON(200, MSG{
			"err": 1,
			"msg": err,
		})
		return
	}
	log.Warnf("get qrcode url = %s", qrurl)
	token := s.getToken(c)
	jpl := jumpLogin{
		CookieJar: s.getCookieJar(c),
		Tk:        token,
		Ip:        c.ClientIP(),
		Ctx:       c,
		Session:   session,
	}
	ckChan <- jpl
	session.Set(cache_key_cookie+"token", []byte(token.Token))
	session.Save()
	c.JSON(200, MSG{
		"err":    0,
		"qrcode": qrurl,
		"token":  token.Token,
	})
}

//jd app登录通过token查cookie
func (s *httpServer) get_cookie_by_token(c *gin.Context) {
	session := sessions.Default(c)
	tokenByte := session.Get(cache_key_cookie + "token")
	if tokenByte == nil {
		c.JSON(200, MSG{
			"err":   404,
			"title": "提取cookie失败",
			"msg":   "请重新提取",
		})
		return
	}
	token := string(tokenByte.([]byte))
	cookie, err := cache.Get(cache_key_cookie + token)
	if err != nil {
		c.JSON(200, MSG{
			"err":   21,
			"title": "提取cookie失败",
			"msg":   "请重新提取，或者再次查询",
		})
		return
	}
	cache.Remove(cache_key_cookie + token)
	session.Clear()
	c.JSON(200, MSG{
		"err":    0,
		"title":  "提取cookie成功",
		"msg":    "cookie:" + cookie.(string),
		"cookie": cookie,
	})
}
