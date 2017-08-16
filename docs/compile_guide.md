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


## Step 3: Setting the configuration file of Board
   
./board/make/board.cfg is the configuration file of Board. Configure the parameters in this file according to your board deployment design. 


## Step 4: Building configuration and env file for each components  
   
### Building
Now we have a **Makefile** for building the whole project. So, you can build configuration and env file very simply.

   ```sh
      $ cd board
      $ make prepare
   ```

## Step 5: Building and Running Board

### Building and Running
You can build and start Board very simply.

   ```sh
      $ make start
   ```

These commands will pull/build the images for Board and run them. Depend on your net speed, it will take a few minitus or hours.

### Building and Running(By compose)
   
   ```sh
      $ docker-compose -f make/dev/docker-compose.yml up
   ```

## Step 6: Building  and Running UI components 

### Building and Runing 
Use the following command to building and running  UI components:

   ```sh
      $ make compile_ui
   ```

These commands will also pull/build the images for Board and run them. Depend on your net speed, it will take a few minitus or hours.

### Building and Running(By compose)
Use `docker-compose` to build Board directly.

The UI components (writen in TypeScript) can be built by running with a UI Builder image separately.
   ```sh
      $ docker-compose -f make/dev/docker-compose.uibuilder.yml up
   ```
   

## Step 7: Verify the Board

Refer to [View and test Board REST API via Swagger](configure_swagger.md) for testing the Board REST API.


## Step 8: Stop Board

When you want to stop Board, run:

   ```sh
      $ make down
   ```

To use compose directly, run:

   ```sh
      $ docker-compose -f make/dev/docker-compose.yml down
   ```



## Appendix
For development and test, Board build the source code in container. But you can build source code in host GO environment.
First, create directory in your $GOPATH/src for Board and copy the source code into it:

   ```sh
      $ mkdir -p $GOPATH/src/git/inspursoft
      $ cp -r board $GOPATH/src/git/inspursoft
   ```

Then you can use Makefile to build the apiserver and so on

   ```sh
      $ make compile_apiserver
      $ make compile_tokenserver
      ...
   ```
