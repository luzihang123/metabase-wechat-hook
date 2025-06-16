package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type MetabaseAlert struct {
	Data struct {
		AlertCondition string `json:"alert_condition"`
		QuestionName   string `json:"question_name"`
	} `json:"data"`
}

type WechatText struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func sendToWechatWebhook(content string, webhookKey string) error {
	msg := WechatText{
		MsgType: "text",
	}
	msg.Text.Content = content

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", webhookKey),
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("Wechat response: %s\n", bodyBytes)
	return nil
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var alert MetabaseAlert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	content := fmt.Sprintf("📢 [Metabase 告警]：%s\n📊 问题：%s",
		alert.Data.AlertCondition,
		alert.Data.QuestionName,
	)

	webhookKey := os.Getenv("WECHAT_WEBHOOK_KEY")
	if webhookKey == "" {
		http.Error(w, "WECHAT_WEBHOOK_KEY not set", http.StatusInternalServerError)
		return
	}

	err := sendToWechatWebhook(content, webhookKey)
	if err != nil {
		http.Error(w, "发送企业微信失败", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/webhook", webhookHandler)
	log.Printf("✅ Listening on :%s ...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

