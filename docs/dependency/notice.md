# 构建通知

## 第三方软件
- [钉钉](https://ding-doc.dingtalk.com/document#/org-dev-guide/qf2nxq)
- [企业微信](https://work.weixin.qq.com/help?doc_id=13376)
- [飞书](https://www.feishu.cn/hc/zh-CN/articles/360024984973-%E6%9C%BA%E5%99%A8%E4%BA%BA-%E5%A6%82%E4%BD%95%E5%9C%A8%E7%BE%A4%E8%81%8A%E4%B8%AD%E4%BD%BF%E7%94%A8%E6%9C%BA%E5%99%A8%E4%BA%BA-)

## 自定义

1. 填写需要通知的API
2. 相应事件完成后，goploy推送信息至API

```
# 项目构建完成
{
    "code": "int",
    "message": "string",
    "data":  {
        "projectId": "int",
        "projectName": "string",
        "branch": "string",
    }
}
```

