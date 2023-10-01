# vgo框架-后端部分
```
开发者: 
@codelunaticer 
https://github.com/codelunaticer
企鹅号: 2282249330

	后端: Gin+redis+mysql+gorm
	前端(Vue3版): vue3+vite+ts+tailwindcss+pinia+vueRouter+axios

版本
mysql 5.7.25
node 18.17.1
redis 6.2.10
go 1.21
```


# 错误状态码

```
422 参数验证错误(请求参数格式不对,或者缺少等)
401 无权限(一般是token校验失败)
404 路径错误,资源不存在(请求路径有问题,资源不存在)
400 其他错误(被捕获的错误)
500 程序错误(被全局捕获的错误,不影响程序)
```

# 鉴权

## jwt级鉴权

```
token有无
token是否有效和过期
密码版本与redis中的是否一致(修改密码,冻结账号,角色更改 密码版本都会被更新)
校验账号是否存在,是否停用
```



## 超管级鉴权

```
必须先通过jwt鉴权
有1号超管角色或者root账户名
```

## SSS级鉴权

```
必须是root账号
```







# 初始数据

```
1. root账户
2. 1-admin 角色
```





## 一 : 开放级

```
无需token验证,任何来源都可访问
```

### 用户模块

#### 1.10 注册账号

```go
post

/user/register

{
    "username": "admin1", // 账号
    "password": "123456" // 密码
}
```



```
接收账号密码
将密码散列hash加密,再将加密过后的账号密码入库
此时密码版本为1 (用处:在鉴权token时,还要对比token中的密码版本是否一致,这样改过密码后,旧token将全部失效)
```

#### 1.11 登录账号

```go
post

/user/login

{
    "username": "admin",// 账号
    "password": "123456"// 密码
}
```



```
接收账号密码
校验账号是否存在,密码是否一致,如果一致将返回token
token含用户id,密码版本
密码版本存放在redis中(5号)
```

## 二 : 用户级(拥有账号)



### 用户模块

#### 2.10 修改某用户密码

```go
put

/user/auth/setPass

{
    "username":"admin", // 账号
    "old_pass":"123456", // 旧密码
    "new_pass":"123456",// 新密码
    "confirm_new_pass":"123456" // 确认密码(新)
}
```



```
接收原账号密码和新密码
校验token
校验旧密码
通过后更新密码
更新完密码更新密码版本号(更新完版本之后旧的token将全部失效,必须重新登录)
```



#### 2.11 获取当前用户信息

```go
get

/user/auth/getUserInfo
```



```
通过校验
返回当前用户信息
```

#### 2.12 修改当前用户信息

```go
put

/user/auth/setUserInfo

{
    "nike_name":"新用户", // 别名
    "avatar":"sdfsdf", // 头像网址
    "qq":"sdfsdf",// qq
    "wechat":"sdfsdf",// 微信
    "email":"sdfsdf",// 邮箱
    "github":"sdfsdf"// github
}
```



```
通过校验
获取字段信息(其中用户id选用token中的id,防止用户上传id字段)(密码等一系列隐私字段排除不修改)
进行修改
```



#### 2.13 修改头像

```go
put

/user/auth/setAvatar

formdata字段: avatar
```



```
通过校验
保存文件到临时目录
文件访问路径入数据库
入库成功后将文件移动到正确的目录-入库失败移回到临时目录
修改成功后将就头像文件删除
返回新的路径
```





### 角色模块

#### 2.21 获取当前用户的角色

```
get

/role/auth/takeCRole
```

```
通过校验
返回角色信息
```



#### 2.22 获取当前用户的菜单

```
get

/role/auth/takeCMenu
```

```
通过校验
返回菜单信息
```





## 三 : 超管级(超管角色持有):

```
        '删除/停用 对root账户或者有超管角色的操作无效,会被拒绝,因为这是超管账号'
```

```
超管账号分两种:
	1. root:拥有撤销和授权某个账号为超管的权限 (SSS权限)
	2. 普通账号+超管角色:没有撤销和授权某个账号为超管的权限,但其他超管权限都有
```



### 用户模块

#### 3.10 查看账号列表

```go
get

/user/root/findAllUser
```



```
验证通过
返回所有账号信息(密码不返回)
```

#### 3.11 冻结某个账号

```go
post

/user/root/disUser

{
    "user_id": 65 // 用户id
}
```



```
验证通过
查找账号是否存在
禁止操作root
修改冻结字段
修改密码版本
```

#### 3.12 解冻某个账号

```go
post

/user/root/openUser

{
    "user_id":65 // 用户id
}
```



```
验证通过
查找账号是否存在
禁止操作root
修改冻结字段
修改密码版本
```



### 角色模块



#### 3.20 添加一个角色分类

```go
post

/role/root/addRole

{
    "role_name": "超级管理员" // 角色名称
}
```



```
校验通过
角色是否已经存在
不存在添加
```



#### 3.21 删除一个角色分类

```go
delete

/role/root/delRole

{
    "role_id": 1 // 角色id
}
```





```
校验通过
角色是否为超管角色(是的话禁止操作)
角色是否存在,存在则删除
```

#### 3.22 修改一个角色分类

```go
put

/role/root/putRole

{
    "role_id": 5, // 角色id
    "role_name": "1" // 角色名称
}
```



```
校验通过
角色是否为超管角色(是的话禁止操作)
角色是否存在,存在则修改
```

#### 3.23 查找单个角色信息

```go
get

/role/root/takeRole/:role_id
```



```
校验通过
返回单个角色信息
```

#### 3.24 查找所有角色信息

```go
get 

/role/root/findAllRole
```



```
校验通过
返回所有角色信息
```

#### 3.25 为某账号添加某个角色

```go
post

/role/root/userAddRole

{
    "user_id":66, // 账号id
    "role_id":6 // 角色id
}
```



```
校验通过
账户是否存在,角色是否存在
相同记录是否存在
都存在就添加
```

#### 3.26 删除某账号的某个角色

```go
post

/role/root/delUserRole

{
    "user_id":66, // 账号id
    "role_id":6 // 角色id
}
```



```
校验通过
存在就删除
```

#### 3.27 查找某用户拥有的角色

```go
post

/role/root/findUserRole

{
    "user_id":71 // 账户id
}
```



```
校验通过
直接返回
```

### 菜单模块

#### 4.10 添加一条菜单

```go
post

/menu/root/add

{
    "menu_name":"df" // 菜单名称
}
```



```
通过校验
菜单是否存在
不存在-添加菜单
```

#### 4.11 删除一条菜单

```go
delete

/menu/root/del/:menu_id

```

```
校验通过
直接删除
```

#### 4.12 修改一条菜单

```go
put

/menu/root/put

{
    "menu_id":2, // 菜单id
    "menu_name":"1" // 菜单名称
}
```

```
校验通过
直接修改
```

#### 4.13 查看一条菜单

```go
get

/menu/root/take/:menu_id
```

```
校验通过
直接查看
```

#### 4.14 查看所有菜单

```go
get

/menu/root/findAll
```

```
校验通过
直接查看
```

#### 4.15 为某角色添加某菜单

```go
post

/menu/root/roleAddMenu

{
    "role_id": 3, // 角色id
    "menu_id": 6 // 菜单id
}
```

```
校验通过
角色是否存在,菜单是否存在
是否存在相同记录
进行添加
```

#### 4.16 删除某角色的某菜单

```go
post

/menu/root/delRoleMenu

{
    "role_id": 3, // 角色id
    "menu_id": 6 // 菜单id
}
```

```
校验通过
直接删除
```

#### 4.17 查找某角色拥有的菜单

```go
get

/menu/root/findRoleMenus/:role_id
```

```
校验通过
直接返回
```



## 四:  SSS 级(root账号独有)



### 用户模块

#### 4.10 为某账户添加超管角色

```go
post

/user/SSS/addRootAccount

{
    "user_id": 71  // 账号id
}
```



```
校验通过
禁止操作root账号
用户是否存在
相同记录是否已经存在
将超管角色添加到账户
更新该账户的密码版本
```

#### 4.11 撤销某账户超管角色

```go
delete

/user/SSS/delRootAccount

{
    "user_id": 71 // 账号id
}
```



```
校验通过
禁止操作root账号
如果不是root则该条记录进行删除
更新该账户的密码版本
```



#### 4.12 重置某个用户的密码

```go
put

/user/SSS/resetUserPass/:user_id
```



```
校验通过
直接重置
更新密码版本
```