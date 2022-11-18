# linebot-go

### API格式
| Method       |       URL                         |    說明             | 備註    |
| -------------| --------------------------------- | ------------------ | -------- |
|  POST        | /pushMessage                      |    推送 line 訊息給user   |         |
|  POST        | /callback                         |    儲存 user 訊息到DB     |         |
|  GET         | /queryMessage                     |    搜尋 user 回傳的訊息    |         |

----
#### 搜尋 user 回傳的訊息 

URL：
```
[GET] http://localhost:9090/queryMessage
```

回傳：
```
{
    "data": [
        {
            "type": "message",
            "text": "你好",
            "reply_token": "aaaaaaaaaaaaaaa"
        },
        {
            "type": "message",
            "text": "Hi",
            "reply_token": "bbbbbbbbbbbbbbb"
        },
        {
            "type": "message",
            "text": "Hello",
            "reply_token": "ccccccccccccccc"
        }
    ]
}
	
```
