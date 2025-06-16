解决问题：
💡 企业微信机器人要求的是如下格式：
```
{
  "msgtype": "text",
  "text": {
    "content": "你好，这是来自 Metabase 的告警通知"
  }
}

```
但 Metabase 发出的测试 webhook 是这种：
```
{
  "notification_type": "alert",
  "data": {
    "alert_condition": "Results changed",
    "alert_id": 123,
    ...
  }
}
```
🚫 所以直接发给企业微信会被拒收，不会有任何提示或返回消息，这是微信机器人的机制限制。

<img width="1792" alt="image" src="https://github.com/user-attachments/assets/6f71464f-e724-4ad0-a625-0c07a88ce8f0" />

✅ 解决方案：使用中转服务做格式转换


✅ 功能说明
- 接收 Metabase Webhook 请求
- 解析告警信息
- 向企业微信机器人 Webhook 转发符合格式的文本消息
- 日志打印状态
- 可配置（机器人 key、监听端口）

