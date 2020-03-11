## 用户指南  
## 概要  
本指南介绍了浪潮容器云服务平台（Board）的基本功能，指导用户如何使用Board系统  

* 用户账户管理
* 基于项目的访问控制
* 管理项目
* 管理项目的成员
* 管理镜像
* 管理服务
  * 构建服务的镜像
  * 构建服务
  * 部署服务
* 查找项目，服务，用户和镜像
* 监控仪表板
* 管理员选项
* 常见问题

## 用户账户管理
Board支持database认证模式，同时也支持LDAP模式 

* **基于数据库(db_auth)**  

    用户存储在本地数据库中。
    
    用户可以在此模式下自己注册。用户点击注册， 输入相关的信息即可注册属于自己的账户。
	
    如果需要禁用用户自注册功能，请参考初始配置安装指南，或者禁用管理员选项中的该特性。在禁用自注册时，系统管理员可以将用户添加到Board中。
	
    在注册或添加新用户时，用户名和电子邮件必须在Board系统中是唯一的， 如果新注册的用户名已经存在， 系统会提示该用户不可用， 请选择其他用户名。密码必须包含至少8个字符，包括1个小写字母、1个大写字母和1个数字字符。


* **基于LDAP (ldap_auth)**  

	在这种身份验证模式下，存储在外部LDAP或AD服务器中的用户可以直接登录Board。系统默认为database模式， 如果需要使用LDAP模式，需要做必要的配置。配置方法参考配置手册。  
	
	当LDAP / AD用户以用户名和密码登录时，使用“LDAP搜索DN”和安装指南中描述的“LDAP搜索密码”来绑定到LDAP / AD服务器。如果成功，Board在LDAP条目“LDAP基本DN”中查找用户，包括substree。“LDAP uid”指定的属性(如uid、cn)用于将用户与用户名匹配。如果找到匹配，则将用户的密码通过绑定请求验证到LDAP / AD服务器。
	
	在LDAP / AD认证模式下，不支持自注册、更改密码和重新设置密码，因为用户是由LDAP或AD来管理的。

## 基于项目的访问控制  

通过容器服务平台上的项目管理服务。用户可作为不同角色的成员添加到系统中:

* **匿名用户:** 当用户未登录时，用户被视为“匿名”用户。匿名用户无法访问私有项目，并且只能访问公共项目和服务。
	
* **注册用户**: 当用户登录后, 用户就会拥有创建新项目的权限或者可以被拉入一个已经存在的项目。
* *** 项目管理者 ***: 当用户创建一个项目后，将会获得这个项目“项目管理者”的角色。项目管理者可以邀请其他用户加入自己创建的项目。
* *** 项目成员 ***: 当被邀请加入一个新项目后，用户会获得这个项目的“项目成员”。项目成员可以在项目中创建或者删除服务，但不能删除项目本身。如果一个用户不是一个项目的成员，是不能在这个项目中创建或删除的项目，也不能访问它的私有服务。
	
* **系统管理员:** “系统管理员”具有最多的特权。除了上面提到的权限之外，“系统管理员”还可以列出所有项目，将普通用户设置为管理员，删除用户。公共项目“库”也由管理员拥有。

## 管理项目
一个项目包含所有的服务，图片等。在Board中有两种类型的项目，他们分别是公有项目和私有项目:

*  **公有:** 所有用户都拥有对公共项目的read权限，您可以通过这种方式共享一些服务或获得其他服务。
*  **私有:** 私有项目只有拥有适当特权的用户访问和使用。

您可以在登录后创建一个项目。检查“公共/私人”复选框将使这个项目公开。


<img src="img/userguide/create-project.png" width="100" alt="Board create project">

在创建项目之后，您可以使用左侧的导航栏浏览服务、用户和镜像。

## 管理项目的成员
### 添加成员

您可以向现有项目添加不同成员， 被添加的成员拥有对该项目的读，写等权限。

<img src="img/userguide/add-members.png" width="100" alt="Board add members">

### 更新和删除成员

您可以通过单击左箭头来更新或删除成员，以便在用户和成员列表的中间添加成员， 用户被删除后失去原有的权限。

<img src="img/userguide/add-remove-members.png" width="100" alt="Board add remove members">

## 管理服务

Board支持创建容器服务。所有服务必须按项目分组。点击“创建服务”。 第一步是选择项目。如果没有项目，请先创建项目。

### 创建镜像

在“选择镜像”界面， 选择“创建自定义镜像”， 从下拉菜单建造新的镜像， 或者在“镜像”界面点击“创建镜像”。
支持以下三种方法创建镜像：“利用模板创建”， “使用Dockerfile文档创建”，“DevOps方式创建”， 
使用模板创建：
会弹出一个窗口为用户提供以下镜像参数输入使用

* 新建镜像名称
* 镜像标签
* 基础镜像
* 镜像进入点
* 镜像环境变量
* 镜像卷
* 镜像执行
* 镜像外部端口
* 上传外部文件
* 命令

使用Dockerfile文件创建：
会弹出一个窗口为用户提供以下镜像参数输入使用

* 新建镜像名称
* 镜像标签
* 选择创建进行使用的doker-file文件

添加所需要的参数后， 点击“构建镜像”开始构建新的镜像。 如果构建成功，新的镜像将被添加到Board的仓库中

### 构建服务
“选择镜像”是构建服务的第一步
选择需要的镜像和镜像标签， 选择多个你需要的镜像

### 选择镜像

“选择镜像”是构建服务的开始。

选择所需的镜像和它的镜像标签，如果需要，选择多个镜像。

可以为这个服务的容器定制以下参数。
* 容器名称
* 工作路径
* 卷挂入点
* 环境变量
* 容器端口
* 命令

下一步是配置服务，服务提供了以下参数
* 服务名称
* 外部服务
* 实例

在高级配置项中， 可以给外部服务分配节点端口

下一步是在配置服务完成之后

### 测试服务
这一步是用了测试服务的配置。 下一步会跳过测试

### 部署服务

单击“部署”部署新服务。
在成功部署服务之后，用户可以从服务列表监视服务状态。


### 创建服务实例

#### “demoshow” 实例
部署服务 "demoshow"

* 登录到Board

<img src="img/userguide/demoshow-a.PNG" width="100" alt="Board login">

* 选择项目

<img src="img/userguide/demoshow-d.PNG" width="100" alt="Select project">

* 选择 library/mydemoshow 镜像

<img src="img/userguide/demoshow-e.PNG" width="100" alt="Select image">

* 配置容器

<img src="img/userguide/demoshow-f.PNG" width="100" alt="Container image">

* 选择容器“mydemoshow” 

<img src="img/userguide/demoshow-g.PNG" width="100" alt="Container name">

* 设置容器的端口为 5000

<img src="img/userguide/demoshow-h.PNG" width="100" alt="Container port">

* 设置服务名称

<img src="img/userguide/demoshow-i.PNG" width="100" alt="Service name">

* 为外部服务设置节点端口

<img src="img/userguide/demoshow-j.PNG" width="100" alt="Service port">

* 部署“demoshow”服务

<img src="img/userguide/demoshow-todeploy.PNG" width="100" alt="Service deploy">

* “demoshow”服务部署成功

<img src="img/userguide/demoshow-deploy.PNG" width="100" alt="Service success">

* 服务能在服务列表中显示
<img src="img/userguide/demoshow-ok.PNG" width="100" alt="Service success">

#### 实例 ：浪潮 “bigdata” 服务
部署“bigdata”服务到项目中：

* 开始创建服务

<img src="img/userguide/bigdata-a.PNG" width="100" alt="create a service">

* 选择一个项目

<img src="img/userguide/bigdata-b.PNG" width="100" alt="Select a project">

* 添加镜像到服务中

<img src="img/userguide/bigdata-c.PNG" width="100" alt="add images">

* 给服务选择一个镜像

<img src="img/userguide/bigdata-d.PNG" width="100" alt="select images">

* 为这个服务选择两个镜像

<img src="img/userguide/bigdata-e.PNG" width="100" alt="select two images">

* 配置容器

<img src="img/userguide/bigdata-f.PNG" width="100" alt="Configure containers">

* 为mysql容器配置存储卷

<img src="img/userguide/bigdata-g.PNG" width="100" alt="Configure storage volume">

* 配置环境变量参数

<img src="img/userguide/bigdata-i.PNG" width="100" alt="Configure environment parameters">

* 配置容器端口

<img src="img/userguide/bigdata-j.PNG" width="100" alt="Configure container ports">

* 配置 bigdata 服务、 配置节点外部端口

<img src="img/userguide/bigdata-k.PNG" width="100" alt="Configure bigdata service">


* 部署 bigdata 服务

<img src="img/userguide/bigdata-m.PNG" width="100" alt="Deploy the bigdata service">

<img src="img/userguide/bigdata-n.PNG" width="100" alt="Deployed the bigdata service">

* bigdata 服务部署完成

<img src="img/userguide/bigdata-p.PNG" width="100" alt="bigdata service deployed">

<img src="img/userguide/bigdata-q.PNG" width="100" alt="bigdata service deployed">

* 监控 bigdata 服务在Board上的状态

<img src="img/userguide/bigdata-o.PNG" width="100" alt="Monitor the bigdata service status">

## 查询

搜索引擎可以搜索项目、服务、用户和镜像。

### 查询的分类

* **项目**:
用户可以通过一些限制搜索项目:


* 普通用户只能搜索这些项目，这些项目是他们共同的项目
* 系统管理员可以搜索所有的项目


* **服务**:
用户可以通过一些限制搜索服务: 

* 普通用户只能搜索服务服务的所有者,或属于同一项目。
* 系统管理员可以搜索所有的项目


* **用户**
可以使用一些约束来搜索用户: 

* 项目管理员可以搜索属于这个项目的用户。
* 系统管理员可以搜索所有用户

* **镜像**

* 普通用户只能搜索属于同一项目或普通图像的镜像。
* 系统管理员可以搜索所有的镜像

### 查询的结果
* **查询的结果** 如下所示
![search](img/userguide/search_result.png)

## 监控仪表板

监视仪表板从k8s主节点和节点收集日志。它涵盖了机器指标，如CPU、内存使用、文件系统和k8s服务运行时。

服务运行时收集所有服务并对应于pod和容器的标签。在仪表板中显示统计实时和平均数字
![search](img/userguide/dashboard_service.png) 

机器指示器收集所有节点的所有CPU和内存节点指示器
![search](img/userguide/dashboard_node.png) 

文件系统收集节点的所有存储指标
![search](img/userguide/dashboard_storage.png) 

## 管理员选项

管理员选项提供用户管理，可由管理员用户添加、更改或删除用户。

* **注意**: 此选项只提供给用户系统管理员角色。

### 查看用户

这个列表显示所有注册用户。
<img src="img/userguide/list-all-users.png" width="100" alt="Board list all users">

* **注意**: 管理员用户是系统默认的第一个用户，不能被修改。

* 具有系统管理员角色的用户可以更改其他用户的权限。

### 管理用户

用户可以通过点击“添加用户”按钮创建。

<img src="img/userguide/add-user.png" width="100" alt="Board add user">

用户可通过点击编辑按钮进行编辑

<img src="img/userguide/edit-user.png" width="100" alt="Board edit user">

从列表删除用户。

<img src="img/userguide/delete-user.png" width="100" alt="Board delete user">


## Q&A

