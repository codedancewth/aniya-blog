---
title: OpenClaw 飞书配置指南
description: 详细配置指南，教你如何将 OpenClaw 与飞书集成，实现通过飞书机器人使用 AI 助手功能
author: 小吴同学
publishDate: 2026-03-25
tags: ['OpenClaw', '飞书', '配置指南', 'AI']
category: '教程'
draft: false
---

## 📝 简介

本文档将指导您如何配置 OpenClaw 与飞书集成，实现通过飞书机器人使用 AI 助手功能。

参考官方配置：https://docs.openclaw.ai/channels/feishu

## 📋 前置准备

1. 已注册的飞书个人企业账号
2. 飞书管理员权限（用于创建应用）
3. OpenClaw 基础环境已部署

## 🚀 配置步骤

访问 [飞书开发平台](https://open.feishu.cn/app) 登录注册一个。

### 1. 创建应用

1. 点击"创建企业自建应用"
2. 填写应用名称"OpenClaw"及相关信息
3. 保存并获取以下凭证，保存完毕后先回到文档启动（调试/训练/在线推理）项目后再执行第四步流程：
   - **App ID**
   - **App Secret**

### 2. 配置应用权限

程序启动后配置权限：

1. 在"权限管理"中添加以下权限：

```json
{
  "scopes": {
    "tenant": [
      "aily:file:read",
      "aily:file:write",
      "application:application.app_message_stats.overview:readonly",
      "application:application:self_manage",
      "application:bot.menu:write",
      "cardkit:card:read",
      "cardkit:card:write",
      "contact:user.employee_id:readonly",
      "corehr:file:download",
      "event:ip_list",
      "im:chat.access_event.bot_p2p_chat:read",
      "im:chat.members:bot_access",
      "im:message",
      "im:message.group_at_msg:readonly",
      "im:message.p2p_msg:readonly",
      "im:message:readonly",
      "im:message:send_as_bot",
      "im:resource"
    ],
    "user": ["aily:file:read", "aily:file:write", "im:chat.access_event.bot_p2p_chat:read"]
  }
}
```

### 3. 配置事件订阅

1. 在事件订阅页面中
2. 选择**使用长连接接收事件（WebSocket）**
3. 添加事件：`im.message.receive_v1`

### 4. 发布应用

1. 在版本管理中创建版本并发布
2. 提交审核并发布
3. 等待管理员审批

## ❓ 常见问题 Q&A

**Q: 报错 "OpenClaw: access not configured."**

A: 查看当前输入的日志链接，输入【`openclaw pairing approve feishu [你的验证码]`】，同时排查是否添加了对应的通讯录权限，如果没有可以权限管理中搜索 "contact:contact.base:readonly" 添加权限试试，若还没有解决，请反馈。

---

如有问题，欢迎在 AI 社区提问 🦞
