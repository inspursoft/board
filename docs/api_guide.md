# API Guide  
## Overview  
This guide walks you through the method of using Board API. You'll learn how to use Board API to:  

* Login
  * /api/v1/sign-in
* Create user 
  * /api/v1/useradd
* Build project
  * /api/v1/projects
* Deployment Service 
  * /api/v1/services/config
  * /api/v1/services/deployment
* Example


## Login

### /api/v1/sign-in

* **Post Method**  

    This endpoint is used to login Board with your account. You can login system with inspur group account if the LDAP system is configed.

  * **Request Parameters** 

          NULL

  * **Request Body** 
  
    ```
      {
        "user_name": "string",
        "user_password": "string"
      }
    ``` 


## Create user

### /api/v1/useradd

* **Post Method**  

    This endpoint is used to create user account for system admin account, when the LDAP system isn't configed in Board. If the LDAP system is configed, this endpoint is invalid.

  * **Request Parameters** 

      Name            | In         | Type        | Description
      ----------------|------------|-------------|--------------
      token           | query      | string      | Current available token   

  * **Request Body** 
  
    ```
      {
        "user_name": "string",            //user name, required field
        "user_email": "string",           //user email, required field
        "user_password": "string",        //user password, password should be at least 8 characters with at least one 
                                          //uppercase, one lowercase and one number, required field
        "user_realname": "string",        
        "user_comment": "string",         
        "user_system_admin": 0,           //system admin account flag, default is 0 
        "user_project_admin": 0,          //project admin account flag, default is 0
        "user_creation_time": "string",   
        "user_update_time": "string"      
      }
    ```                

## Build project

### /api/v1/project

* **Post Method**  

    This endpoint is for user to create a new project.

  * **Request Parameters** 

      Name            | In         | Type        | Description
      ----------------|------------|-------------|--------------
      token           | query      | string      | Current available token   

  * **Request Body** 
  
    ```
      {
        "project_name": "string",            //project name, required field
        "project_comment": "string",         //comment for project, default is null
        "project_public": 0,                 //project public flag, deault is 0 and private project
      }
    ```       

* **Get Method** 

    This endpoint is used to get all projects information in Board. 

  * **Request Parameters**  

      Name            | In         | Type        | Description
      ----------------|------------|-------------|--------------
      token           | query      | string      | Current available token          

  * **Request Body**  

          NULL


## Deployment Service 
Deployment service needs to config the service first.
 
### /api/v1/services/config

* **Post Method** 
    
    This endpoint is used to create service configure.

  * **Request Parameters**  

      Name            | In         | Type        | Description
      ----------------|------------|-------------|--------------
      token           | query      | string      | Current available token          
      phase           | query      | string      | Set phase of config service

  * **Request Body**                                  

      ```
      {
        "project_id": 0,               // project ID in board, type is int64
        "service_id": 0,               // not fill
        "instance": 0,                 // the number of instance, type is int and range is >0
        "service_name": "string",      // service name must be lowercase, type is string 
        "container_list": [            // containers config, type is array of struct 
          {
            "name": "string",          // container name, type is string
            "working_Dir": "string",   // working dir, type is string, it can be null
            "command": "string",       // exec command when container starts, type is string, it can be null
            "container_port": [        // expose port, type is array of int
              0
            ],
            "volume_mounts": {                        // struct about volume mounts, it can be null 
              "target_storage_service": "string",     // service name for volume mounts, type is string
              "target_path": "string",                // service path for volume mounts, type is string 
              "volume_name": "string",                // volume name, type is string
              "container_path": "string"              // path in the container for volumne mounts, type is string
            },
            "image": {                                // struct about image
              "image_name": "string",                 // image name, tpye is string
              "image_tag": "string",                  // image tag, tpye is string
              "project_name": "string"                // project name that image belonged to, type is string
            },
            "env": [                                  // array of environment variable, it can be null 
              {
                "dockerfile_envname": "string",       // the key of environment variable, type is string
                "dockerfile_envvalue": "string"       // the value of environment variable, type is string
              }
            ]
          }
        ],
        "external_service_list": [
          {
            "container_name": "string",    // container name that corresponds to the container_list's container name, type is string 
            "node_config": {               // struct of node config
              "target_port": 0,            // expose port in the container
              "node_port": 0               // default range:30000~32767
            }
          }
        ]
      }
      ```

* **Get Method** 

    This endpoint is used to get service configure. 

  * **Request Parameters**  

      Name            | In         | Type        | Description
      ----------------|------------|-------------|--------------
      token           | query      | string      | Current available token          
      phase           | query      | string      | Set phase of config service    

  * **Request Body**  

          NULL

### /api/v1/services/deployment

* **Post Method** 
    
    This endpoint is used to deploy the service which had been configed by API /api/v1/services/config.

  * **Request Parameters**  

      Name            | In         | Type        | Description
      ----------------|------------|-------------|--------------
      token           | query      | string      | Current available token   

  * **Request Body**  

          NULL


## Example
   
### Examples to create services by API

Deploy a service "demoshow" with swagger.

* Config swagger 

    Refer to [View and test Board REST API via Swagger](configure_swagger.md) for configure swagger. Then the wagger ui can be visited with the url http://host:port/swagger/index.html. In this example, we will use the url http://10.110.18.232:8080/swagger/index.html to visit swagger ui.
    
    Swagger demo:

    <img src="img/apiguide/demoshow-swagger.PNG" width="100" alt="Swagger-ui">

* Login board
    
    The below request body can be used to login board. The user of 'admin' is system admin account, and has the authority to create new user.

    ```
      {
        "user_name": "admin",
        "user_password": "123456a?"
      }
    ``` 
    
    Swagger demo:
    
    <img src="img/apiguide/demoshow-login.PNG" width="100" alt="Board login">

    We will get token from the response
    
    <img src="img/apiguide/demoshow-login-response.PNG" width="100" alt="Get token">  

* Create new user for deployment service

    Now, a new user of 'inspur' for deployment service will be created, and the new user will have the authority to build project and deployment servcie.

    ```
      {
        "user_name": "inspur002",            
        "user_email": "inspur@inspur.com",           
        "user_password": "123456Aa",             
        "user_comment": "deployment service user",              
        "user_project_admin": 1     
      }
    ```    

    Swagger demo:

    <img src="img/apiguide/demoshow-adduser.PNG" width="100" alt="Create user">

* login Board with new user

    The below request body can be used to login board.

    ```
      {
        "user_name": "inspur002",
        "user_password": "123456Aa"
      }
    ``` 
    
    Swagger demo:
    
    <img src="img/apiguide/demoshow-newuserlogin.PNG" width="100" alt="Board relogin">

    The token will be got from the response.

    <img src="img/apiguide/demoshow-newusertoken.PNG" width="100" alt="Get new user token">  

* Create new project for deployment service

    A new private project will be create by the below request body.

    ```
      {
        "project_name": "deploy001",
        "project_public": 0,
        "project_comment": "private project"
      }
    ``` 

    Swagger demo:
    
    <img src="img/apiguide/demoshow-createproject.PNG" width="100" alt="Create project">  

* Get the new project ID for deployment servcie

    The below request body will be used to get projects information.

    Swagger demo:
    
    <img src="img/apiguide/demoshow-getprojectinfo.PNG" width="100" alt="Get project info">  

    And we will get the project ID from the response.

    <img src="img/apiguide/demoshow-getprojectID.PNG" width="100" alt="Get project ID"> 

* Create service configure

    A service named 'demoshowing001' use the new user will be builded, which is belonged to the new project and visited through 30005 port of host. It has one instance. The instance has one container named 'demoshowing001'. The container's exposed port is 5000, image name is 'library/mydemoshowing', the image tag is '1.0', the image is belonged to the public project 'library'. The below request body will be used to create the service and the 'phase' parameter in request is 'ENTIRE_SERVICE'.

    ```
      {
        "project_id": 23,
        "instance": 1,
        "service_name": "demoshowing001",
        "container_list": [
          {
            "name": "demoshowing001",
            "container_port": [
              5000
            ],
            "image": {
            "image_name": "library/mydemoshowing",
            "image_tag": "1.0",
            "project_name": "library"
            }
          }
        ],
        "external_service_list": [
          {
            "container_name": "demoshowing001",
            "node_config": {
              "target_port": 5000,
              "node_port": 32500
            }
          }
        ]
      }
    ``` 

    Swagger demo:

    <img src="img/apiguide/demoshow-configserviceurl.PNG" width="100" alt="config service 1"> 
    <img src="img/apiguide/demoshow-configservice.PNG" width="100" alt="config service 2"> 

* Deploy service
    
    Service configed in the upper section will be deployed.

    Swagger demo:

    <img src="img/apiguide/demoshow-deployment.PNG" width="100" alt="deploy service"> 




