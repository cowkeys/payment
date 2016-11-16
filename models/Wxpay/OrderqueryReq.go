package Wxpay

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"strings"

	"github.com/astaxie/beego"
)

type OrderqueryReq struct {
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Transaction_id string `xml:"transaction_id"`
	Nonce_str      string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
}

type OrderqueryResp struct {
	Return_code      string `xml:"return_code"`
	Return_msg       string `xml:"return_msg"`
	Appid            string `xml:"appid"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Sign             string `xml:"sign"`
	Result_code      string `xml:"result_code"`
	Openid           string `xml:"prepay_id"`
	Trade_type       string `xml:"trade_type"`
	Trade_state      string `xml:"trade_state"`
	Bank_type        string `xml:"bank_type"`
	Total_fee        string `xml:"total_fee"`
	Cash_fee         int    `xml:"cash_fee"`
	Transaction_id   string `xml:"transaction_id"`
	Out_trade_no     string `xml:"out_trade_no"`
	Time_end         string `xml:"time_end"`
	Trade_state_desc string `xml:"trade_state_desc"`
}

func (o *OrderqueryReq) WxQueryOrder(transId string) OrderqueryResp {
	xmlResp := OrderqueryResp{}

	query_order_req := "https://api.mch.weixin.qq.com/pay/orderquery"

	var qeuryReq OrderqueryReq
	qeuryReq.Appid = "wx74a6f50de79ccd40" //微信开放平台我们创建出来的app的app id
	qeuryReq.Mch_id = "1319106901"
	qeuryReq.Transaction_id = transId
	qeuryReq.Nonce_str = randStr(32, "alphanum")

	//beego.Debug("qeuryReqEntity",qeuryReq)
	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = qeuryReq.Appid
	m["mch_id"] = qeuryReq.Mch_id
	m["transaction_id"] = qeuryReq.Transaction_id
	m["nonce_str"] = qeuryReq.Nonce_str
	appkey := beego.AppConfig.String("wxappkey")
	qeuryReq.Sign = WxpayCalcSign(m, appkey) //这个是计算wxpay签名的函数上面已贴出

	//beego.Debug("qeuryReq",qeuryReq)
	bytes_req, err := xml.Marshal(qeuryReq)
	if err != nil {
		fmt.Println("以xml形式编码发送错误, 原因:", err)
		return xmlResp
	}

	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "OrderqueryReq", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", query_order_req, bytes.NewReader(bytes_req))
	if err != nil {
		fmt.Println("New Http Request发生错误，原因:", err)
		return xmlResp
	}
	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}
	resp, _err := c.Do(req)
	if _err != nil {
		fmt.Println("请求微信支付查询发送错误, 原因:", _err)
		return xmlResp
	}

	//xmlResp :=UnifyOrderResp{}
	body, _ := ioutil.ReadAll(resp.Body)
	_err = xml.Unmarshal(body, &xmlResp)
	if xmlResp.Return_code == "FAIL" {
		fmt.Println("微信支付统查询不成功，原因:", xmlResp.Return_msg)
		return xmlResp
	}
	//beego.Debug("xmlResp",xmlResp)

	return xmlResp

}
