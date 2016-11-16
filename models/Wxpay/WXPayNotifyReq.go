package Wxpay

import (
	_ "bytes"
	_ "crypto/md5"
	_ "crypto/rand"
	_ "encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	_ "net/http"
	_ "reflect"
	_ "sort"
	_ "strconv"
	"strings"
	_ "time"

	_ "KiWiFi/gate.cloud/models"
	_ "gatecloud"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_ "github.com/bitly/go-simplejson"
	//"encoding/json"
)

type WXPayNotifyReq struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce          string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Openid         string `xml:"openid"`
	Is_subscribe   string `xml:"is_subscribe"`
	Trade_type     string `xml:"trade_type"`
	Bank_type      string `xml:"bank_type"`
	Total_fee      int    `xml:"total_fee"`
	Fee_type       string `xml:"fee_type"`
	Cash_fee       int    `xml:"cash_fee"`
	Cash_fee_Type  string `xml:"cash_fee_type"`
	Transaction_id string `xml:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	Time_end       string `xml:"time_end"`
}

type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

func (o *WXPayNotifyReq) WxpayCallback(ctx *context.Context) map[string]interface{} {

	body, err := ioutil.ReadAll(ctx.Input.Context.Request.Body)
	if err != nil {
		fmt.Println("读取http body失败，原因!", err)
	}
	fmt.Println("微信支付异步通知，HTTP Body:", string(body))
	var mr WXPayNotifyReq
	err = xml.Unmarshal(body, &mr)
	if err != nil {
		fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
	}

	//beego.Debug("body",body)

	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["return_code"] = mr.Return_code
	reqMap["return_msg"] = mr.Return_msg
	reqMap["appid"] = mr.Appid
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce
	reqMap["result_code"] = mr.Result_code
	reqMap["openid"] = mr.Openid
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["trade_type"] = mr.Trade_type
	reqMap["bank_type"] = mr.Bank_type
	reqMap["total_fee"] = mr.Total_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["cash_fee_type"] = mr.Cash_fee_Type
	reqMap["transaction_id"] = mr.Transaction_id
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["attach"] = mr.Attach
	reqMap["time_end"] = mr.Time_end

	var resp WXPayNotifyResp

	//进行签名校验
	if WxpayVerifySign(reqMap, mr.Sign) {
		//这里就可以更新我们的后台数据库了，其他业务逻辑同理。
		//beego.Debug("succes","succes")
		resp.Return_code = "SUCCESS"
		resp.Return_msg = "OK"
		ctx.WriteString("SUCCESS")
		return reqMap
	} else {
		resp.Return_code = "FAIL"
		resp.Return_msg = "failed to verify sign, please retry!"
	}

	//结果返回，微信要求如果成功需要返回return_code "SUCCESS"
	bytes, _err := xml.Marshal(resp)
	strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
	if _err != nil {
		fmt.Println("xml编码失败，原因：", _err)
		//return
	}
	////beego.Debug("sa",strResp)
	ctx.WriteString(strResp)

	//c.Ctx.ResponseWriter.WriteHeader(200)
	//fmt.Fprint(c.Ctx.ResponseWriter, strResp)

	return reqMap
}

//微信支付签名验证函数
func WxpayVerifySign(needVerifyM map[string]interface{}, sign string) bool {
	appkey := beego.AppConfig.String("wxappkey")
	signCalc := WxpayCalcSign(needVerifyM, appkey)

	//beego.Debug("计算出来的sign: ", signCalc)
	//beego.Debug("微信异步通知sign: ", sign)
	if sign == signCalc {
		fmt.Println("签名校验通过!")
		return true
	}

	fmt.Println("签名校验失败!")
	return false
}