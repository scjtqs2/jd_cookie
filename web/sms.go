package web

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/scjtqs/jd_cookie/util"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SmsSession 验证码登录 连接上下文
type SmsSession struct {
	Guid       string `json:"guid"`
	Lsid       string `json:"lsid"`
	Gsalt      string `json:"gsalt"`
	RsaModulus string `json:"rsa_modulus"`
	ErrCode    int64  `json:"err_code"`
	ErrMsg     string `json:"err_msg"`
	Phone      string `json:"phone"`
}

// toQuickCookie 转换成quick接口的cookie
func (s *SmsSession) toQuickCookie() string {
	return fmt.Sprintf("guid=%s; lsid=%s; gsalt=%s; rsa_modulus=%s; ", s.Guid, s.Lsid, s.Gsalt, s.RsaModulus)
}

func (s *SmsSession) toJson() []byte {
	b, _ := json.Marshal(s)
	return b
}

// sendSms 通过手机号获取短信
func (s *httpServer) sendSms(c *gin.Context) (*SmsSession, error) {
	session := sessions.Default(c)
	var (
		appid       = 959
		version     = "1.0.0"
		countryCode = 86
		timestamp   = time.Now().UnixMilli()
		cmd         = 36
		subCmd      = 1
		gsalt       = "sb2cwlYyaCSN1KUv5RHG3tmqxfEb8NKN"
	)
	phone := c.Query("phone")
	// jar := s.getCookieJar(c)
	client := &http.Client{
		// Jar:       jar,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   time.Second * 10,
	}
	gsign := util.Md5(fmt.Sprintf("%d%s%d%d%d%s", appid, version, timestamp, cmd, subCmd, gsalt))
	dataVal1 := url.Values{}
	dataVal1.Add("client_ver", version)
	dataVal1.Add("gsign", gsign)
	dataVal1.Add("appid", strconv.Itoa(appid))
	dataVal1.Add("return_page", "https%3A%2F%2Fcrpl.jd.com%2Fn%2Fmine%3FpartnerId%3DWBTF0KYY%26ADTAG%3Dkyy_mrqd%26token%3D")
	dataVal1.Add("cmd", strconv.Itoa(cmd))
	dataVal1.Add("sdk_ver", "1.0.0")
	dataVal1.Add("sub_cmd", strconv.Itoa(subCmd))
	dataVal1.Add("qversion", version)
	dataVal1.Add("ts", fmt.Sprintf("%d", timestamp))
	req1, err := http.NewRequest("POST", "https://qapplogin.m.jd.com/cgi-bin/qapp/quick", strings.NewReader(dataVal1.Encode()))
	if err != nil {
		return nil, err
	}
	req1.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 Mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({packageName:com.vivo.hybrid,type:deeplink,extra:{}})")
	req1.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req1.Header.Add("Accept-Encoding", "")
	res1, err := client.Do(req1)
	if err != nil {
		return nil, err
	}
	if res1 != nil {
		//goland:noinspection GoDeferInLoop
		defer res1.Body.Close()
	}
	r, err := io.ReadAll(res1.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("semdSms r1=%s", string(r))
	g1 := gjson.ParseBytes(r).Get("data")
	ck := SmsSession{
		Guid:       g1.Get("guid").String(),
		Lsid:       g1.Get("lsid").String(),
		Gsalt:      g1.Get("gsalt").String(),
		RsaModulus: g1.Get("rsa_modulus").String(),
		Phone:      phone,
	}
	subCmd = 2
	timestamp = time.Now().UnixMilli()
	gsalt = ck.Gsalt
	gsign = util.Md5(fmt.Sprintf("%d%s%d%d%d%s", appid, version, timestamp, cmd, subCmd, gsalt))
	sign := util.Md5(fmt.Sprintf("%d%s%d%s4dtyyzKF3w6o54fJZnmeW3bVHl0$PbXj", appid, version, countryCode, phone))
	dataVal2 := url.Values{}
	dataVal2.Add("country_code", strconv.Itoa(countryCode))
	dataVal2.Add("client_ver", version)
	dataVal2.Add("gsign", gsign)
	dataVal2.Add("appid", strconv.Itoa(appid))
	dataVal2.Add("mobile", phone)
	dataVal2.Add("sign", sign)
	dataVal2.Add("cmd", strconv.Itoa(cmd))
	dataVal2.Add("sub_cmd", strconv.Itoa(subCmd))
	dataVal2.Add("qversion", version)
	dataVal2.Add("ts", strconv.FormatInt(timestamp, 10))
	req2, err := http.NewRequest("POST", "https://qapplogin.m.jd.com/cgi-bin/qapp/quick", strings.NewReader(dataVal2.Encode()))
	if err != nil {
		return nil, err
	}
	req2.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 Mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({packageName:com.vivo.hybrid,type:deeplink,extra:{}})")
	req2.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req2.Header.Add("Accept-Encoding", "")
	req2.Header.Add("Cookie", ck.toQuickCookie())
	res2, err := client.Do(req2)
	if err != nil {
		return nil, err
	}
	if res2 != nil {
		defer res2.Body.Close()
	}
	r2, err := io.ReadAll(res2.Body)
	log.Printf("semdSms r2=%s", string(r))
	if err != nil {
		return nil, err
	}
	if !gjson.ValidBytes(r2) {
		return nil, errors.New("接口访问失败")
	}
	g2 := gjson.ParseBytes(r2)
	// log.Println("code code response :", string(r2))
	ck.ErrCode = g2.Get("err_code").Int()
	ck.ErrMsg = g2.Get("err_msg").String()
	if g2.Get("data.rsa_modulus").Exists() {
		ck.RsaModulus = g2.Get("data.rsa_modulus").String()
	}
	if g2.Get("data.gsalt").Exists() {
		ck.Gsalt = g2.Get("data.gsalt").String()
	}
	if g2.Get("data.lsid").Exists() {
		ck.Lsid = g2.Get("data.lsid").String()
	}
	if g2.Get("data.guid").Exists() {
		ck.Guid = g2.Get("data.guid").String()
	}
	// _ = s.updateCookieJar(c, jar)
	if ck.ErrCode != 0 {
		return nil, errors.New(ck.ErrMsg)
	}
	session.Set("SmsSession", ck.toJson())
	_ = session.Save()
	return &ck, nil
}

// checkCode 通过验证码 换取cookie
func (s *httpServer) checkCode(c *gin.Context) (*Token, error) {
	session := sessions.Default(c)
	var smsSession SmsSession
	smsSessionByte := session.Get("SmsSession")
	if smsSessionByte == nil {
		return nil, errors.New("empty session")
	}
	err := json.Unmarshal(smsSessionByte.([]byte), &smsSession)
	if err != nil {
		return nil, err
	}
	log.Printf("smsSession : %+v", smsSession)
	code := c.Query("code")
	var (
		appid       = 959
		version     = "1.0.0"
		countryCode = 86
		timestamp   = time.Now().UnixMilli()
		cmd         = 36
		subCmd      = 3
		gsign       = util.Md5(fmt.Sprintf("%d%s%d%d%d%s", appid, version, timestamp, cmd, subCmd, smsSession.Gsalt))
	)
	dataVal := url.Values{}
	dataVal.Add("country_code", strconv.Itoa(countryCode))
	dataVal.Add("client_ver", version)
	dataVal.Add("gsign", gsign)
	dataVal.Add("smscode", code)
	dataVal.Add("appid", strconv.Itoa(appid))
	dataVal.Add("mobile", smsSession.Phone)
	dataVal.Add("cmd", strconv.Itoa(cmd))
	dataVal.Add("sub_cmd", strconv.Itoa(subCmd))
	dataVal.Add("gversion", version)
	dataVal.Add("ts", strconv.FormatInt(timestamp, 10))
	// jar := s.getCookieJar(c)
	client := &http.Client{
		// Jar:       jar,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   time.Second * 10,
	}
	req, err := http.NewRequest("POST", "https://qapplogin.m.jd.com/cgi-bin/qapp/quick", strings.NewReader(dataVal.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 Mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({packageName:com.vivo.hybrid,type:deeplink,extra:{}})")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("Accept-Encoding", "")
	req.Header.Add("Cookie", smsSession.toQuickCookie())
	log.Printf("cookie=%s", smsSession.toQuickCookie())
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	r, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if !gjson.ValidBytes(r) {
		return nil, errors.New("接口访问失败")
	}
	g := gjson.ParseBytes(r)
	log.Println("check code response :", string(r))
	if g.Get("err_code").Int() > 0 {
		return nil, errors.New(g.Get("err_msg").String())
	}
	token := &Token{
		Guid:  smsSession.Guid,
		Lsid:  smsSession.Lsid,
		PtKey: g.Get("data.pt_key").String(),
		PtPin: g.Get("data.pt_pin").String(),
	}
	token.UserCookie = "pt_key=" + token.PtKey + ";pt_pin=" + token.PtPin + ";"
	log.Info("############  登录成功，获取到 Cookie  #############")
	log.Infof("Cookie1=%s", token.UserCookie)
	log.Info("####################################################")
	s.cleanSession(c)
	return token, nil
}

// getSmsCode http路由 get_sms_code
func (s *httpServer) getSmsCode(ctx *gin.Context) {
	// 校验手机号是否合法
	phone := ctx.Query("phone")
	if phone == "" {
		ctx.JSON(http.StatusOK, MSG{
			"err":   400,
			"title": "参数错误",
			"msg":   "手机号不能为空",
		})
		return
	}
	if !regexp.MustCompile(`^\d{11}$`).MatchString(phone) {
		ctx.JSON(http.StatusOK, MSG{
			"err":   400,
			"title": "参数错误",
			"msg":   "手机号格式错误",
		})
		return
	}
	_, err := s.sendSms(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, MSG{
			"code":  500,
			"title": "发送验证码失败",
			"msg":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, MSG{
		"err":   200,
		"title": "发送验证码成功",
		"msg":   "请继续输入手机验证码",
	})
}

// checkSmsCode http路由 check_sms_code
func (s *httpServer) checkSmsCode(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusOK, MSG{
			"err":   400,
			"title": "参数错误",
			"msg":   "请输入手机验证码",
		})
		return
	}
	if !regexp.MustCompile(`^\d{6}$`).MatchString(code) {
		ctx.JSON(http.StatusOK, MSG{
			"err":   400,
			"title": "参数错误",
			"msg":   "请输入6位验证码",
		})
	}
	tk, err := s.checkCode(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, MSG{
			"err":   500,
			"title": "校验验证码失败",
			"msg":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, MSG{
		"err":    200,
		"cookie": tk.UserCookie,
	})
}
