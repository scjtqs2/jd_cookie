import com.alibaba.fastjson.JSONObject;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.ContentType;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.util.DigestUtils;

import java.util.regex.Pattern;

/**
 * @author Coe
 * @since 2023-06-17 13:26
 */
@Service
public class SMSAction {
    private static final Logger logger = LoggerFactory.getLogger(SMSAction.class);

    private final CloseableHttpClient client = HttpClients.createDefault();
    private final static Pattern phonePattern = Pattern.compile("^1[3-9][0-9]{9}$");
    private final static Pattern codePattern = Pattern.compile("^\\d{6}$");
    private String guid;
    private String lSid;
    private String gSalt;
    private String rsaModulus;

    public String getSmsCode(String phone){
        if(!phonePattern.matcher(phone).matches()){
            return "手机号格式不正确";
        }
        try{
            sendSms(phone);
        }catch(Exception e){
            return "发送验证码失败，" + e.getMessage();
        }
        return "发送验证码成功,请继续输入手机验证码";
    }

    public String checkSmsCode(String phone,String code){
        if (!phonePattern.matcher(phone).matches() || !codePattern.matcher(code).matches()){
            return "手机号或验证码格式不正确";
        }
        String cookie;
        try {
            cookie = checkCode(phone,code);
        }catch (Exception e){
            return "登录失败，原因：" + e.getMessage();
        }
        return cookie;
    }

    /**
     * 输入手机号获取短信
     * @param phone 手机号
     * @throws Exception 运行中错误
     */
    public void sendSms(String phone) throws Exception {
        int appid = 959;
        String version = "1.0.0";
        int countryCode = 86;
        long timestamp = System.currentTimeMillis();
        int cmd = 36;
        int subCmd = 1;
        String gsalt = "sb2cwlYyaCSN1KUv5RHG3tmqxfEb8NKN";
        String gsign = DigestUtils.md5DigestAsHex(String.format("%d%s%d%d%d%s", appid, version, timestamp, cmd, subCmd, gsalt).getBytes());

        String data1 = String.format("client_ver=%s&gsign=%s&appid=%d&return_page=https%%3A%%2F%%2Fcrpl.jd.com%%2Fn%%2Fmine%%3FpartnerId%%3DWBTF0KYY%%26ADTAG%%3Dkyy_mrqd%%26token%%3D&cmd=%d&sdk_ver=%s&sub_cmd=%d&qversion=%s&ts=%d",
                version, gsign, appid, cmd, version, subCmd, version, timestamp);

        logger.info("获取验证码时第一个请求体:{}",data1);
        HttpPost req1 = new HttpPost("https://qapplogin.m.jd.com/cgi-bin/qapp/quick");
        req1.setHeader("User-Agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 Mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({packageName:com.vivo.hybrid,type:deeplink,extra:{}})");
        req1.setHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8");
        req1.setHeader("Content-Type", "application/x-www-form-urlencoded; charset=utf-8");
        req1.setHeader("Accept-Encoding", "");
        req1.setEntity(new StringEntity(data1, ContentType.APPLICATION_FORM_URLENCODED));

        JSONObject res1 = JSONObject.parseObject(EntityUtils.toString(client.execute(req1).getEntity()));
        logger.info("json1:{}",res1);

        if (!res1.containsKey("data")){
            throw new Exception("第一个接口发生错误，错误编码为：" + res1.getInteger("err_code") + ",错误信息为：" + res1.getString("err_msg"));
        }

        JSONObject resData1 = res1.getJSONObject("data");
        String guid = resData1.getString("guid");
        String lSid = resData1.getString("lsid");
        gsalt = resData1.getString("gsalt");
        String rsaModulus = resData1.getString("rsa_modulus");

        this.guid = guid;
        this.lSid = lSid;
        this.gSalt = gsalt;
        this.rsaModulus = rsaModulus;

        subCmd = 2;
        timestamp = System.currentTimeMillis();
        gsign = DigestUtils.md5DigestAsHex(String.format("%d%s%d%d%d%s", appid, version, timestamp, cmd, subCmd, gsalt).getBytes());
        String sign = DigestUtils.md5DigestAsHex(String.format("%d%s%d%s4dtyyzKF3w6o54fJZnmeW3bVHl0$PbXj", appid, version, countryCode, phone).getBytes());

        String data2 = String.format("country_code=%s&client_ver=%s&gsign=%s&appid=%d&mobile=%s&sign=%s&cmd=%d&sub_cmd=%d&qversion=%s&ts=%d",
                countryCode, version, gsign, appid, phone, sign, cmd, subCmd, version, timestamp);

        logger.info("获取验证码时第二个请求体:{}",data2);
        HttpPost req2 = new HttpPost("https://qapplogin.m.jd.com/cgi-bin/qapp/quick");
        req2.setHeader("User-Agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 Mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({packageName:com.vivo.hybrid,type:deeplink,extra:{}})");
        req2.setHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8");
        req2.setHeader("Content-Type", "application/x-www-form-urlencoded; charset=utf-8");
        req2.setHeader("Accept-Encoding", "");
        req2.setHeader("Cookie", getFormat());
        req2.setEntity(new StringEntity(data2, ContentType.APPLICATION_FORM_URLENCODED));

        JSONObject res2 = JSONObject.parseObject(EntityUtils.toString(client.execute(req2).getEntity()));
        logger.info("json2:{}",res2);

        if (!res2.containsKey("data")){
            throw new Exception("第二个接口发生错误，错误编码为：" + res2.getInteger("err_code") + ",错误信息为：" + res2.getString("err_msg"));
        }
    }

    /**
     * 用手机号和验证码换取cookie
     * @param phone 手机号
     * @param code 验证码
     * @return cookie
     */
    private String checkCode(String phone, String code) throws Exception {
        int appid = 959;
        String version = "1.0.0";
        int countryCode = 86;
        long timestamp = System.currentTimeMillis();
        int cmd = 36;
        int subCmd = 3;
        String gSign = DigestUtils.md5DigestAsHex((appid + version + timestamp + cmd + subCmd + this.gSalt).getBytes());

        String reqBody = String.format("country_code=%s&client_ver=%s&gsign=%s&smscode=%s&appid=%d&mobile=%s&cmd=%d&sub_cmd=%d&qversion=%s&ts=%d",
                countryCode, version, gSign, code, appid, phone, cmd, subCmd, version, timestamp);

        logger.info("验证码登录时的请求体：{}",reqBody);
        HttpPost req = new HttpPost("https://qapplogin.m.jd.com/cgi-bin/qapp/quick");
        req.setHeader("Connection", "Keep-Alive");
        req.setHeader("Content-Type", "application/x-www-form-urlencoded; charset=utf-8");
        req.setHeader("Accept", "application/json, text/plain, /");
        req.setHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8");
        req.setHeader("User-Agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 Mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({packageName:com.vivo.hybrid,type:deeplink,extra:{}})");
        req.setHeader("Cookie", getFormat());
        req.setEntity(new StringEntity(reqBody, ContentType.APPLICATION_FORM_URLENCODED));

        JSONObject res = JSONObject.parseObject(EntityUtils.toString(client.execute(req).getEntity()));

        if (!res.containsKey("err_code")){
            throw new Exception("接口访问失败");
        }

        if (res.getInteger("err_code") > 0){
            throw new Exception(res.getString("err_msg"));
        }

        JSONObject data = res.getJSONObject("data");
        logger.info("验证码登录获取到的data：{}",data);
        return "pt_key=" + data.getString("pt_key") + ";pt_pin=" + data.getString("pt_pin") + ";";
    }

    private String getFormat() {
        String format = String.format("guid=%s;lsid=%s;gsalt=%s;rsa_modulus=%s;", guid, lSid, gSalt, rsaModulus);
        logger.info("给京东的加盐Cookie：{}",format);
        return format;
    }
}
