package controllers

import (
	"bytes"
	"common/base"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"math/rand"
	"net/http"

	"time"

	"os"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}
type Rq struct {
}
type Request struct {
	App_id string `json:"app_id"`

	Img_data     string    `json:"img_data"`
	Rsp_img_type string    `json:"rsp_img_type"`
	Opdata       []*Opdata `json:"opdata"`
}
type Opdata struct {
	Cmd    string  `json:"cmd"`
	Params *Params `json:"params"`
}
type Params struct {
	Model_id string `json:"model_id"`
}
type Image struct {
	Img_base64 string `json:"img_base64"`
}

func (this *MainController) Get() {

	SecretID := "AKIDznNjqhqjsmRt634ESRBk47uQilX2lNAg"
	secretKey := "RfdKmzyiDSSEl36yc1C8qIUQ9euud7SZ"
	appID := "10117716"
	userID := "656712495"

	imgData, err := ioutil.ReadFile("./static/img/xx.jpg")

	if err != nil {
		fmt.Fprintf(os.Stderr, "ReadFile() failed: %s\n", err)
		return
	}

	ntime := base.GetCurrentDataUnix()
	//

	sign := fmt.Sprintf("a=%s&k=%s&e=%d&t=%d&r=%d&u=%s&f=",
		appID,
		SecretID,
		ntime+1000,
		ntime,
		rand.Int31(),
		userID)
	//

	h := hmac.New(sha1.New, []byte(secretKey))
	beego.Debug(sign)
	h.Write([]byte(sign))
	hm := h.Sum(nil)
	//attach orig_sign to hm
	dstSign := []byte(string(hm) + sign)
	b64 := base64.StdEncoding.EncodeToString(dstSign)
	////请求body
	params := new(Params)
	params.Model_id = "hezuo_junzhuangzhao_1999w_20170919141627"
	od := new(Opdata)
	od.Cmd = "doFaceMerge"
	od.Params = params
	rt := new(Request)
	rt.App_id = string(appID)
	rt.Img_data = base64.StdEncoding.EncodeToString(imgData)
	rt.Rsp_img_type = "base64"
	rt.Opdata = append(rt.Opdata, od)

	byt, err := json.Marshal(rt)
	beego.Debug(string(byt))
	br := bytes.NewReader(byt)

	////
	hr, err := http.NewRequest("POST", "http://api.youtu.qq.com/cgi-bin/pitu_open_access_for_youtu.fcg", br)
	if err != nil {
		beego.Error(err)
		return
	}
	hr.Header.Add("Authorization", b64)
	hr.Header.Add("Content-Type", "text/json")
	hr.Header.Add("Host", "api.youtu.qq.com")
	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	res, err := client.Do(hr)
	if err != nil {
		beego.Error(err)
		return
	}
	resbyt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		beego.Error(err)
		return
	}
	img := new(Image)
	err = json.Unmarshal(resbyt, img)
	if err != nil {
		beego.Error(err)
		return
	}
	imgbyt, err := base64.StdEncoding.DecodeString(img.Img_base64)
	if err != nil {
		beego.Error(err)
		return
	}
	err = this.Ctx.Output.Body(imgbyt)
	if err != nil {
		beego.Error(err)
		return
	}
}
