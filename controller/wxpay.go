package Payment

import (
	"encoding/base64"
	"fmt"
	"odeke-em/qr"
	"os"
	"payment/models/Wxpay"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type WxpayController struct {
	beego.Controller
}

func (this *WxpayController) Native() {
	orderNumber := this.Ctx.Input.Param(":id") //获取订单号
	payAmount := this.GetString("price")       //获取价格
	params := make(map[string]interface{})
	params["body"] = "****company-" + orderNumber //显示标题
	params["out_trade_no"] = orderNumber
	params["total_fee"] = payAmount
	params["product_id"] = orderNumber
	params["attach"] = "abc" //自定义参数

	var modwx Wxpay.UnifyOrderReq
	res := modwx.CreateOrder(this.Ctx, params)

	this.Data["data"] = res
	//拿到数据之后，需要生成二维码。
	this.Data["Image"] = Img(res.Code_url)

	this.TplName = "Wxpay/index.tpl"
}

func (this *WxpayController) Notify() {
	var notifyReq Wxpay.WXPayNotifyReq
	res := notifyReq.WxpayCallback(this.Ctx)
	//beego.Debug("res",res)
	if res != nil {
		//这里可以组织res的数据 处理自己的业务逻辑：
		sendData := make(map[string]interface{})
		sendData["id"] = res["out_trade_no"]
		sendData["trade_no"] = res["transaction_id"]
		paid_time, _ := time.Parse("20060102150405", res["time_end"].(string))
		paid_timestr := paid_time.Format("2006-01-02 15:04:05")
		sendData["paid_time"] = paid_timestr
		sendData["payment_type"] = "wxpay"
		intfee := res["cash_fee"].(int)
		floatfee := float64(intfee)
		cashfee := floatfee / 100
		sendData["payment_amount"] = strconv.FormatFloat(cashfee, 'f', 2, 32)

		//api(sendData)...自己的逻辑处理
		//

	}
}

func Img(url string) string {
	code, err := qr.Encode(url, qr.H)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	imgByte := code.PNG()
	str := base64.StdEncoding.EncodeToString(imgByte)

	return str
}
