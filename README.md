wocket 连接测试交互工具

使用命令： -r 后面的参数只是一个参考格式，不能原样复制(默认是有以下请求头，原样输入会报错请求头重复)，一般ws请求时header中通常会根据业务情况会有一些特殊的请求头，请以示例格式输入

`ws -h ws://127.0.0.1:8000 -r '{"Upgrade":"websocket", "Accept-Language": "zh-CN,zh;q=0.9"}'
`
 如果要退出交互命令行,连续4个 control + c ，并按 enter 可退出程序
 
 <img width="723" alt="image" src="https://user-images.githubusercontent.com/73807441/153709066-3574cd28-27b7-4c80-80fe-0af1475c83ea.png">
<img width="763" alt="image" src="https://user-images.githubusercontent.com/73807441/153709084-e689dc16-e53d-4024-a472-1507843d7760.png">
