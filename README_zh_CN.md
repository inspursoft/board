# Board

[English](README.md) | [中文](README_zh_CN.md)

**注意**：开发过程中，`master`分支可能处于*不稳定的甚至中断的状态*。
请使用`releases`分支，而不是`master`的分支，来获得稳定的二进制文件。

<img alt="Board" src="docs/img/board_logo.png">

"Board"产品是基于docker+ kubernetes的容器服务平台，为浪潮软件提供云解决方案，包括轻量级容器的虚拟化，微服务，DevOps，持续的交付，帮助企业和开发团队实现快速的业务应用交付和不断创新。

## 特性
* **用户账户**：Board支持数据库的认证方式和LDAP模式。
* **基于项目的访问控制**：通过容器服务平台上的项目管理服务，用户可作为不同角色的成员添加到系统中。
* **管理项目**：一个项目包含所有的服务，图片等。
* **管理服务**：Board支持创建容器服务。所有服务必须按项目分组。
* **查询**：Board搜索引擎可以搜索项目、服务、用户和镜像。
* **监控仪表板**：监视仪表板从k8s主节点和节点收集日志，它涵盖了机器指标，如CPU、内存使用、文件系统和k8s服务运行时。
* **管理员选项**：管理员选项提供用户管理，可由管理员用户添加、更改或删除用户。

## 安装
**系统要求：**
Board作为几个Docker容器部署，因此，Board可以部署在任何支持Docker的Linux发行版。
* Python的版本应为2.7或更高。请注意，您可能需要对不附带默认安装Python解释器的Linux发行版（Gentoo，ARCH）安装Python。
* Docker engine的版本应为1.11.2或更高。有关安装说明，请参考：https://docs.docker.com/engine/installation/
* Docker Compose的版本应为1.7.1或更高。有关安装说明，请参考：https://docs.docker.com/compose/install/

### 在线安装
（即将推出）

### 离线安装
安装步骤如下

1. 下载安装程序；
2. 配置**board.cfg**；
3. 运行**install.sh**安装，然后启动Board
注意：如果你需要准备Kubernetes和注册表环境，请参阅附录部分。


#### 下载安装程序
安装程序的二进制文件可以从`release`页面下载，选择在线或离线安装程序，使用*tar*命令提取包。
在线安装程序:
（即将推出）
离线安装程序:
```sh
    $ tar xvf board-offline-installer-latest.tgz.tgz
```

#### 配置Board
配置参数所在文件为**board.cfg**。

board.cfg中有两种类型的参数，**所需的参数**和**可选参数**。

* **所需的参数**：这些参数需要在配置文件中设置。他们将在用户更新```board.cfg```，并运行```install.sh```脚本重新安装Board后生效。
* **可选参数**：这些参数是可选的更新。后续更新这些参数```board.cfg```将被忽略。

参数说明如下 - 请注意，至少，你需要更改**hostname**属性。

##### 所需的参数
* **hostname**：目标主机的主机名，这是用来访问用户界面和API服务器服务。它应该是IP地址或者目标机器的完全限定域名（FQDN），例如，`192.168.1.10`或`reg.yourdomain.com`。 _不要使用`localhost`或`127.0.0.1`作为主机名 - 该API服务器服务需要由外部客户端访问！_
* **db_password**：MySQL数据库root密码用于**db_auth** 。 _任何生产使用请修改该密码！_

##### 可选参数
* *****：
* *****：
* *****：
* *****：
* *****：
* *****：
* *****：






















