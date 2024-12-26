## 例子

#### 直接放在登录后面即可

```vue
<script>
  let ws;
  let token = "your-jwt-token"; // 假设这是用户的 JWT token
  let heartBeatInterval = 3000; // 3秒发送一次心跳包 一般承载上限 3w/s

  // 初始化 WebSocket 连接
  function connectWebSocket() {
    ws = new WebSocket("ws://localhost:8000/echo"); //和后端对齐

    ws.onopen = function (evt) {
      console.log("WebSocket连接已打开");
// 连接成功后发送 token 作为第一个消息
      sendToken(token);
      startHeartbeat();  // 连接成功后开始发送心跳包
    };

    ws.onclose = function (evt) {
      console.log("WebSocket连接关闭");
    };

    ws.onmessage = function (evt) {
// 处理来自后端的消息
      console.log("收到消息: " + evt.data);
    };

    ws.onerror = function (evt) {
      console.log("WebSocket错误:", evt);
    };
  }

  // 发送 token 作为第一个消息
  function sendToken(token) {
    const tokenMessage = JSON.stringify({
      type: "authenticate",  // 自定义消息类型，表示这是认证消息
      token: token
    });
    ws.send(tokenMessage);
    console.log("发送token:", tokenMessage);
  }

  // 发送心跳包
  function sendHeartbeat() {
    const heartBeatMessage = JSON.stringify({
      type: "heartbeat",  // 这是心跳包的类型
    });
    ws.send(heartBeatMessage);
    console.log("发送心跳包", heartBeatMessage);
  }

  // 启动心跳包发送
  function startHeartbeat() {
    setInterval(sendHeartbeat, heartBeatInterval);  // 每隔30秒发送一次心跳包
  }

  connectWebSocket();  // 建立连接
</script>
```


### 发送请求例子
```vue
<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Invitation</title>
</head>
<body>
    <script>
        var socket = new WebSocket("ws://localhost:8080/ws");

        // 连接打开时的回调
        socket.onopen = function(event) {
            console.log("Connected to WebSocket server.");
            // 模拟发送邀请消息
            var inviteMessage = {
                action: "invite",
                userID: 1,
                teamID: 2,
                groupID: 3
            };
            socket.send(JSON.stringify(inviteMessage));
        };

        // 监听消息
        socket.onmessage = function(event) {
            var message = event.data;
            console.log("Received message:", message);
            // 在此处理消息，可能是展示通知或其他UI更新
        };

        // 连接关闭时的回调
        socket.onclose = function(event) {
            console.log("Disconnected from WebSocket server.");
        };

        // 连接错误时的回调
        socket.onerror = function(event) {
            console.log("WebSocket error:", event);
        };
    </script>
</body>
</html>

```