# 基于Qt和Golang的即时通讯系统

## 开发环境

操作系统:  OpenSuse
数据库:    MariaDB-10.2.25
编程语言:  Golang1.15 , Qt5.15.2

## 服务端部署流程

**服务端**

1. 安装Golang;
2. 开启 Go Moudle;
4. 创建数据库,不需创建数据表,程序运行后,数据表的创建由程序自动管理;
5. 进入config文件夹rocket.ini配置数据库,服务器ip,port等参数;
6. 进入服务器目录,执行`go build litrocket.go`,编译成功将生成名为litrocket的可执行程序,执行即可。


## 主要功能
### 服务端主要功能
	1. 用户登录、注册(登录密码加密)    ok  
	2. 消息转发(个人消息，群聊消息)      
	3. 搜索好友与群,添加删除好友与群,好友推荐
	4. 群文件
	5. 动态
	7. 分享屏幕

### 客户端
	1. 截图
	2. 查询常用单词
	3. 消息翻译


### 待做功能:
	1. 合理关闭服务端,而不是强制关闭   :  向服务端发送关机消息，接收后，服务端关闭全部listener(即用户不再可以登录)，待处理客户端数据的协程全部结束后，关闭服务器
	2. 在线升级程序
	3. 对json进行加密传输
	4. 语音消息，视频通话,等



## 数据库表的设计   (操作数据库使用Gorm库)

使用GORM框架，Gorm是协程安全的，不需要手动处理并发问题.

### 用户表,User
gorm.Model  //自动创建三个字段(ID(主键且自增),createtime,updatetime,deletetime)
UserName    //用户名
PassWord    //密码
UserSex     //性别,  0男1女
UserAge     //年龄
UserImage   //用户头像(存储图片的路径)

### 用户信息表,UserInfo
UserID      //用户ID
MyGroupID   //与我有关的群聊
GroupRole   //0表示我创建的群聊,1表示我加入的群聊

### 用户好友表,Friend
UserID      //用户id
FriendID    //其朋友的id
FriendName  //好友备注

### 群组成员表,Group
GroupID     //群聊ID
UserRole    //用户角色,0为管理员,1为普通人
UserID      //加入此群聊的人的id

### 群组信息表,GroupInfo
gorm.Model  //自动创建三个字段
GroupName   //群聊名字
GroupImage  //群聊头像(存储图片的路径)
GroupRootID //群主
GroupUserNum  //群组人数

### 离线消息表,Message
MsgDate     //消息的发送日期
Message    //消息内容
SrcID      //发送者ID
DestID     //接收者ID
MessFormat //消息格式,可能是文字(0)或者是图片(1),如果是图片则消息内容为图片路径
MessDest   //目标是群聊还是个人,0为群聊，1为个人

### 文件记录表，File
Filetime   // 文件创建时间
FileSize   // 文件大小(bytes)
Filename   // 文件名
FileLocation // 文件在服务器的路径
SrcID      // 发送者ID
DestID     // 目标ID
MessDest   // 目标是群聊还是个人

### 存储用户动态表,Article      (看gblog的教程)
Date   //动态时间
UserID //发动态的人
Content//动态的文字
Image  //动态的图片


## log日志模块的设计

   handlelog包下封装函数Handlelog()
   
   1. 函数将所有异常信息都写入到以天命名的文件内,不使用log包的Fatal开头的函数,他会直接退出程序,本程序在可以预料的可能发生Fatal错误的地方使用Handlelog函数打印log到文件和标准输出,然后使用手动return,做到正常关闭和正常销毁资源,唯一使用log.Fatalln的地方就是Handlelog的实现部分,没办法,按理说log功能不正常时要关闭服务,所以使用了Fatalln

   2. log分为三个等级,INFO,WARNING,FATAL,  
		INFO只是输出正常信息,
		WARNING代表有异常发生需要检查程序,
		FATAL表示发生致命错误,程序立即退出,所有服务立即停止
