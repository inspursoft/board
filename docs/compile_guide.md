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

### Compiling and Running

You can compile the code by the approache:

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

Target              | Description
--------------------|-------------
all                 | composition of the target : vet, fmt, golint and compile_ui 
prepare             | prepare configuration and env 
start               | building and running board instance
compile_ui          | building and running UI builder 
down                | shutdown board instance
clean_binary        | clean apiserver,tokenserver and collector/cmd binary
install             | compile apiserver,tokenserver and collector/cmd binary
test                | used to test a program written in the Go language 
fmt                 | format the code for all the go language source files in the code package 
vet                 | check for static errors in the go source file 
golint              | used to check for irregularities in the go code  
build               | build apiserver, tokenserver, db, log, collector images
rmimage             | remove apiserver, tokenserver, db, log, collector images 

#### EXAMPLE:


#### compile apiserver,tokenserver and collector/cmd binary 

   ```sh
      $ make install
   ```

   **Note**: the board file path:$GOPATH/src/git/inspursoft/

