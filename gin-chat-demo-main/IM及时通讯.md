# README

IM及时通讯项目是一个使用 Go 语言和 Gin 框架开发的聊天应用示例项目。该项目提供了基础的聊天功能，包括用户注册、登录、单聊、群聊，群管理等。

## 项目特点

- 用户认证系统(jwt 实现登录保持)
- 私密数据采用AES加密算法的CBC模式进行密钥加密
- 私聊和群聊功能
- 邮箱验证
- 基于 WebSocket 的实时消息传递
- 用redis储存在线用户列表，MongoDB储存历史对话消息，MySQL将用户，群，好友信息持久化
- RESTful API 设计
- 已部署到阿里云服务器上

## 项目结构

```
gin-chat-demo-main/
|-- api/                # API 路由和控制器
|-- cache/               # 缓存操作
|-- conf/                # 配置文件
|-- main.go              # 主入口
|-- model/               # 数据模型
|-- pkg/                 # 工具包
|-- router/              # 路由配置
|-- service/             # 服务层
|-- serializer/          # 数据序列化
`-- ws/                  # WebSocket 相关
```

## 环境要求

- Go 1.15 或以上版本
- MySQL 数据库
- Redis 缓存服务
- MongoDB 文档数据库（用于存储聊天记录）

tip:项目已部署，https://blog.xiaordong.cn在本文档中以{{url}}表示![image-20240728131612632](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728131612632.png)



数据库表设计：

~~~sql
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password_digest` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL UNIQUE,
  `avatar` varchar(1000),
  `phone` varchar(255),
  `status` varchar(255),
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `friends` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `owner_id` int(11) NOT NULL,
  `target_id` int(11) NOT NULL,
  `type` int NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `owner_id_target_id_unique` (`owner_id`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_name` varchar(255) NOT NULL,
  `group_owner` int(11) NOT NULL,
  `desc` varchar(100),//群简介
  `size` int NOT NULL DEFAULT 200,//大小，默认200人
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `group_members` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `role` int NOT NULL,//角色，群主为1，管理员2，普通成员0
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
~~~

**返回值设定**

~~~json
{
    "status":200,
    "Data":"返回前端需要的数据",
    "Msg":"携带消息",
    "Error":"可能的错误信息"
}
~~~

**状态码**

~~~go
	SUCCESS               = 200  "ok"
	UpdatePasswordSuccess = 201  "修改密码成功"
	NotExistInentifier    = 202  "该第三方账号未绑定"
	ERROR                 = 500  "fail"
	InvalidParams         = 400  "请求参数错误"
	ErrorDatabase         = 40001  "数据库操作出错,请重试"

	ErrorHeaderData = 30001  "缺少Authorization数据"
	ErrorValidToken = 30002 "无效的Token"
	ErrorGetFile    = 30003  "从网络获取文件失败"
	ErrorCreateFile = 30004  "创建文件目录失败"
	ErrorPlaceFile  = 30005  "文件放置失败"
	DeleteError     = 30006  "删除数据失败"
	ErrorCreateData = 30007  "创建数据失败"
	ErrorGroupFull  = 30008   "群人数已满"
	ErrorNoData     = 30009  "无法获取群数据"
	ExistUser       = 30010  "用户已在群中"

	WebsocketSuccessMessage = 50001  "解析content内容信息"
	WebsocketSuccess        = 50002  "发送信息，请求历史纪录操作成功"
	WebsocketEnd            = 50003  "请求历史纪录，但没有更多记录了"
	WebsocketOnlineReply    = 50004  "针对回复信息在线应答成功"
	WebsocketOfflineReply   = 50005  "针对回复信息离线回答成功"
	WebsocketLimit          = 50006  "请求收到限制"
~~~

## 接口文档

### 用户注册

**请求路径：**

```http
GET {{url}}/user/register
```

**请求参数：**

- `username` (string, required): 用户名
- `password` (string, required): 密码
- `key`(string,len =  16):邀请码

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": null,
    "msg": "ok",
    "error": ""
}
~~~

### 用户登录

**请求路径：**

~~~http
GET {{url}}/user/login
~~~

**请求参数**：

- `username` (string, required): 用户名
- `password` (string, required): 密码

**成功返回示例**:

~~~json
{
    "status": 200,
    "data": {
        "aToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdHVkZW50X2lkIjoiMiIsImV4cCI6MTcyMjE3MzgxMSwiaXNzIjoid2F0ZXJTeXN0ZW0ifQ.zfXQAmbM2gzzHpfHaTBgaMYKom8LehRjjU2xob_RsA0",
        "rToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjMzNTQ2MTEsImlzcyI6Im15LXByb2plY3QifQ.o-rlchKGcorbofWlAV_Xg_qtTlu_I9Uy7bLtnsUZoes"
    },
    "msg": "ok",
    "error": ""
}
~~~

### 查找用户

**请求路径**：

~~~http
GET {{url}}/user/find
~~~

**请求参数：**

- `info` (string, required): 搜索用户id

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述
- `data` (array): 用户列表

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": {
        "ID": 1,
        "CreatedAt": "2024-07-25T08:42:14Z",
        "UpdatedAt": "2024-07-25T08:42:14Z",
        "DeletedAt": null,
        "UserName": "test1",
        "PasswordDigest": "$2a$12$GZFtrMzU5fouBCVYGM9qfOuVD1A2GzfiO0N5iWT98V3T8iyt0IaY2",
        "Email": "",
        "Avatar": "http://q1.qlogo.cn/g?b=qq&nk=294350394&s=640",
        "Phone": "",
        "Status": "active"
    },
    "msg": "ok",
    "error": ""
}
~~~

## 注销用户

**请求路径：**

```http
DELETE /user/delete
```

**请求头：**

- `Authorization` (string, required): 当前的访问令牌

**请求参数：**（为了确保安全，注销用户时需要在登录情况下再次提供用户名及密码）

- `username` (string, required): 用户名
- `password` (string, required): 密码

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": null,
    "msg": "ok",
    "error": ""
}
~~~

## 用户邮箱验证

**请求路径：**

```http
POST /user/emailCheck
```

**请求参数：**

- `email` (string, required): 邮箱地址

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": null,
    "msg": "ok",
    "error": ""
}
~~~

验证码截图：![image-20240728135336023](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728135336023.png)

## 完善用户资料

**请求路径：**

```http
PUT /user/update
```

**请求头：**

- `Authorization` (string, required): 当前的访问令牌

**请求参数：**

- `email` (string, required): 新的邮箱地址
- `verification_code` (string, required): 邮箱验证码
- `phone` (string, optional): 新的电话号码
- `avatar` (file, optional): 新的头像图片

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述
- `data` (string): 更新后的头像文件路径

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": "./pkg/upload/personal/172214607514362432.png",
    "msg": "ok",
    "error": ""
}
~~~

## 添加好友

**请求路径：**

```http
POST /friend/add
```

**请求头：**

- `Authorization` (string, required): 当前的访问令牌

**请求参数：**

- `target_id` (uint, required): 好友的用户 ID

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": null,
    "msg": "ok",
    "error": ""
}
~~~

## 获取好友列表

**请求路径：**

```
GET /friend/list
```

**请求头：**

- `Authorization` (string, required): 当前的访问令牌

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述
- `data` (array): 好友列表

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": [
        {
            "ID": 4,
            "CreatedAt": "2024-07-28T05:56:28Z",
            "UpdatedAt": "2024-07-28T05:56:28Z",
            "DeletedAt": null,
            "owner_id": 2,//我id
            "target_id": 3,//好友id
            "Type": 1
        },
        {
            "ID": 5,
            "CreatedAt": "2024-07-28T05:58:54Z",
            "UpdatedAt": "2024-07-28T05:58:54Z",
            "DeletedAt": null,
            "owner_id": 2,
            "target_id": 1,
            "Type": 1
        }
    ],
~~~

## 删除好友

**请求路径：**

```
POST /friend/delete
```

**请求头：**

- `Authorization` (string, required): 当前的访问令牌

**请求参数：**

- `target_id` (uint, required): 好友的用户 ID

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": null,
    "msg": "ok",
    "error": ""
}
~~~

### 私聊

**思路**：

- `WsHandler` 函数处理 WebSocket 连接升级。
- 开启两个协程监听消息，每一个客户端都将缓存到channel
- `Read` 方法读取客户端发送的消息，并根据消息类型进行相应处理（例如，发送消息或拉取历史消息）。
- `Write` 方法处理服务器向客户端发送的消息。
- `Start` 方法在 `ClientManager` 中监听连接状态变化和消息广播。
- 客户端，消息结构：

~~~go
// SendMsg 发送消息的类型
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// ReplyMsg 回复的消息
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// Client 用户类
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// Broadcast 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// ClientManager 用户管理
type ClientManager struct {
	Clients   map[string]*Client
	Broadcast chan *Broadcast
	Reply     chan *Client
	Online    chan *Client
	Offline   chan *Client
}

// Message 信息转JSON (包括：发送者、接收者、内容)
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:   make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast: make(chan *Broadcast),
	Online:    make(chan *Client),
	Reply:     make(chan *Client),
	Offline:   make(chan *Client),
}
~~~



**请求路径**：

~~~http
ws://{{url}}/chat/ws
~~~

**请求参数**：

* `uid`(int,required):用户ID
* `toUid`(int,required):对话用户ID

**聊天截图**:![image-20240728140810293](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728140810293.png)

当另一个用户在线时发送消息反馈：

![image-20240728140908301](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728140908301.png)

另一个用户能够实时接收消息：

![image-20240728141000805](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728141000805.png)

另一个用户也可以查看历史消息：（查看历史消息的type:2）

![image-20240728141420206](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728141420206.png)

## 群组相关接口

### 创建群组

**请求路径：**

```http
POST /group/create
```

**请求头：**

- `Authorization` (string, required): 访问令牌

**请求参数：**

- `group_name` (string, required): 群组名称
- `desc` (string, required): 群组描述

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述
- `data` (object): 群组信息

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": {
        "ID": 2,
        "CreatedAt": "2024-07-28T14:23:31.5502506+08:00",
        "UpdatedAt": "2024-07-28T14:23:31.5502506+08:00",
        "DeletedAt": null,
        "group_name": "debug",
        "group_owner": 2,
        "desc": "用于测试的群",
        "size": 200
    },
    "msg": "ok",
    "error": ""
}
~~~

### 加入群组

**请求路径：**

```http
POST /group/join
```

**请求头：**

- `Authorization` (string, required): 访问令牌

**请求参数：**

- `group_id` (uint, required): 群组 ID

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述

**成功返回示例：**

```json
{
    "status": 200,
    "data": null,
    "msg": "ok",
    "error": ""
}
```

### 获取用户创建的群组

**请求路径：**

```http
GET /group/mine
```

**请求头：**

- `Authorization` (string, required): 访问令牌

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述
- `data` (array): 群组列表

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": [
        {
            "ID": 1,
            "CreatedAt": "2024-07-27T18:38:03Z",
            "UpdatedAt": "2024-07-27T18:39:14Z",
            "DeletedAt": null,
            "group_name": "debug",
            "group_owner": 2,
            "desc": "用于测试的群",
            "size": 200
        },
        {
            "ID": 2,
            "CreatedAt": "2024-07-28T06:23:32Z",
            "UpdatedAt": "2024-07-28T06:23:32Z",
            "DeletedAt": null,
            "group_name": "debug",
            "group_owner": 2,
            "desc": "用于测试的群",
            "size": 200
        }
    ],
    "msg": "ok",
    "error": ""
}
~~~

### 获取用户加入的群组

**请求路径：**

```http
GET /group/joined
```

**请求头：**

- `Authorization` (string, required): 访问令牌

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述
- `data` (array): 群组列表

**成功返回示例：**

```json
{
    "status": 200,
    "data": [
        {
            "ID": 1,
            "CreatedAt": "2024-07-27T18:38:03Z",
            "UpdatedAt": "2024-07-27T18:39:14Z",
            "DeletedAt": null,
            "group_name": "debug",
            "group_owner": 2,
            "desc": "用于测试的群",
            "size": 200
        },
        {
            "ID": 2,
            "CreatedAt": "2024-07-28T06:23:32Z",
            "UpdatedAt": "2024-07-28T06:23:32Z",
            "DeletedAt": null,
            "group_name": "debug",
            "group_owner": 2,
            "desc": "用于测试的群",
            "size": 200
        }
    ],
    "msg": "ok",
    "error": ""
}
```

### 退出群组

**逻辑**:退出群时，如果时群主退出，会将群主移交给下一个管理员或入群时间第二的普通成员，否则，解散群

**请求路径：**

```http
DELETE /group/delete
```

**请求头：**

- `Authorization` (string, required): 访问令牌

**请求参数：**

- `group_id` (uint, required): 群组 ID

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述

**成功返回示例：**

```json
{
    "status": 200,
    "data": null,
    "msg": "ok",
    "error": ""
}
```

### 群聊 WebSocket 连接

**请求路径：**

```http
GET ws://{{url}}chat/group_ws?group_id=群组ID&uid=用户ID
```

**请求头：**

- `Authorization` (string, required): 访问令牌

**返回参数：**

- WebSocket 连接

一位成员发送消息

![image-20240728143537500](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728143537500.png)

其他两位成员同样接收到消息：

![image-20240728143627150](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728143627150.png)

![image-20240728143651236](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728143651236.png)

### 发送图片

**请求路径**：

~~~http
POST {{url}}/user/upload
~~~

**请求头：**

- `Authorization` (string, required): 访问令牌

**请求参数：**

- `file` (File, required): 文件

**返回参数：**

- `status` (int): 状态码
- `msg` (string): 消息描述
- `data`:（string）:文件的url

**成功返回示例：**

~~~json
{
    "status": 200,
    "data": {
        "url": "./pkg/upload/data/172214876618229056.png"
    },
    "msg": "ok",
    "error": ""
}
~~~

这样，前端就拿到了对应文件的路径，同样的视频和语音也是如此

![image-20240728144115431](C:\Users\Keith\AppData\Roaming\Typora\typora-user-images\image-20240728144115431.png)