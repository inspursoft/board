## Introduction

This guide provides instructions for developers to build and run Board from source code.


## Step 1: Prepare for a build environment for Board

Board is deployed as several Docker containers and most of the code is written in Go language. The build environment requires Python, Docker, Docker Compose and golang development environment. Please install the below prerequisites:


Software              | Required Version
----------------------|--------------------------
docker-ce             | 17.03 +
docker-compose        | 1.14.0 +
python                | 2.7 +
git                   | 1.8.3 +
make                  | 3.81 +
golang                | 1.8.1 +


## Step 2: Getting the source code

   ```sh
      $ git clone http://10.110.18.40:10080/inspursoft/board.git
   ```


## Step 3: Building and installing Board

### Configuration

Edit the file **make/board.cfg** and make necessary configuration changes such as hostname, admin password ,mysql server and so on. Refer to **[Installation and Configuration Guide](installation_guide.md#configuring-board)** for more info.

   ```sh
      $ cd board
      $ vi make/board.cfg
   ```

### Compiling and Running in development

Make sure DEVFLAG=dev in Makefile, then you can compile and running by the approache:

#### I. Prepare configuration and env 

   ```sh
      $ make prepare
   ```

#### Ⅱ. Compile UI  

   ```sh
      $ make compile_ui
   ```
   
#### Ⅲ. Building and Running Board

   ```sh
      $ make start
   ```

## Step 4: Verify the Board

Refer to [View and test Board REST API via Swagger](configure_swagger.md) for testing the Board REST API.


## Step 5: Stop Board Instance

When you want to stop Board instance, run:

   ```sh
      $ make down
   ```


## Appendix
* Using the Makefile

The `Makefile` contains predefined targets:

Target                           | Description
---------------------------------|-------------
all                              | Check board source file by fmt vet golint and compile  
prepare                          | Prepare configuration and env 
compile                          | Check board source file by fmt vet golint and compile
compile_ui                       | Building ui
start                            | Building images containers and running
test                             | Runs the tests
fmt                              | Formats board source files
vet                              | Examines board source code and reports suspicious constructs 
golint                           | Linter for source code
clean                            | Print help infomation about clean
cleanall                         | Clean binary and images 
cleanbinary                      | Clean binary 
cleanimage                       | Clean images
down                             | Stop and remove board instance 
install                          | Compile board binary
build                            | Build board images
container/mysql_build            | Build mysql image
container/apiserver_build        | Build apiserver image
container/collector_build        | Build collector image
container/gitserver_build        | Build gitserver image
container/jenkins_build          | Build jenkins image
container/log_build              | Build log image
container/nginx_build            | Build nginx image
container/tokenserver_build      | Build tokenserver image
container/mysql_rmi              | Clean mysql image
container/apiserver_rmi          | Clean apiserver image
container/collector_rmi          | Clean collector image
container/gitserver_rmi          | Clean gitserver image
container/jenkins_rmi            | Clean jenkins image
container/log_rmi                | Clean log image
container/nginx_rmi              | Clean nginx image
container/tokenserver_rmi        | Clean tokenserver image
src/apiserver_clean              | Clean apiserver binary
src/apiserver_fmt                | Formats apiserver source files
src/apiserver_install            | Compile apiserver binary
src/apiserver_vet                | Examines apiserver source code and reports suspicious constructs
src/apiserver_compile            | Check apiserver source file by fmt vet golint and compile
src/apiserver_golint             | Linter for apiserver source code
src/apiserver_test               | Runs apiserver tests
src/tokenserver_clean            | Clean tokenserver binary
src/tokenserver_fmt              | Formats tokenserver source files
src/tokenserver_install          | Compile tokenserver binary
src/tokenserver_vet              | Examines tokenserver source code and reports suspicious constructs
src/tokenserver_compile          | Check tokenserver source file by fmt vet golint and compile
src/tokenserver_golint           | Linter for tokenserver source code
src/tokenserver_test             | Runs tokenserver tests
src/collector/cmd_clean          | Clean collector binary
src/collector/cmd_fmt            | Formats collector source files
src/collector/cmd_install        | Compile collector binary
src/collector/cmd_vet            | Examines collector source code and reports suspicious constructs
src/collector/cmd_compile        | Check collector source file by fmt vet golint and compile
src/collector/cmd_golint         | Linter for collector source code
src/collector/cmd_test           | Runs collector tests




#### EXAMPLE:


#### Compile apiserver,tokenserver and collector/cmd binary 

   ```sh
      $ make install
   ```
   
#### Formats apiserver source code

   ```sh
      $ make src/apiserver_fmt 
   ```

   **Note**: the board file path:$GOPATH/src/git/inspursoft/