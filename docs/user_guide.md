# User Guide  
## Overview  
This guide walks you through the fundamentals of using Board. You'll learn how to use Board to:  

* User Account
* Role Based Access Control(RBAC)  
* Manage Projects
* Manage Members of a Project
* Manage Services
  * Building images.
  * Building Services.
  * Testing.
  * Deployment.
* Search Projects, Services, Users and Images.
* Monitoring Dashboard
* Administrator Options
* Q&A

## User account
Board supports databse authentication modes, will support LDAP next release:  

* **Database(db_auth)**  

	Users are stored in the local database.  
	
	A user can register himself/herself in Board in this mode. To disable user self-registration, refer to the [installation guide](installation_guide.md) for initial configuration, or disable this feature in [Administrator Options](#administrator-options). When self-registration is disabled, the system administrator can add users into Board. 
	
	When registering or adding a new user, the username and email must be unique in the Board system. The password must contain at least 8 characters with 1 lowercase letter, 1 uppercase letter and 1 numeric character.  
	
	When you forgot your password, you can follow the below steps to reset the password(need SMTP server support):  

	1. Click the link "Forgot Password" in the sign in page.  
	2. Input the email address entered when you signed up, an email will be sent out to you for password reset.  
	3. After receiving the email, click on the link in the email which directs you to a password reset web page.  
	4. Input your new password and click "Save".  
	
* **LDAP/Active Directory (ldap_auth)**  

	Support in next release

## Role Based Access Control(RBAC)  

Board manages services through projects on container service platform. Users can be added into one services as a member with 3 different roles:  

* **Guest**: Guest has read-only privilege for public project and sevices.
* **ProjectMember**: Developer has read and write privileges for a project.
* **ProjectAdmin**: When creating a new project, you will be assigned the "ProjectAdmin" role to the project. Besides read-write privileges, the "ProjectAdmin" also has some management privileges, such as adding and removing members.

Besides the above three roles, there are two system-wide roles:  

* **SysAdmin**: "SysAdmin" has the most privileges. In addition to the privileges mentioned above, "SysAdmin" can also list all projects, set an ordinary user as administrator, delete users. The public project "library" is also owned by the administrator.  
* **Anonymous**: When a user is not logged in, the user is considered as an "Anonymous" user. An anonymous user has no access to private projects and has read-only access to public projects and services.  

## Manage Projects

## Manage Members of a Project

## Manage Services

### Building images

### Building Services

### Testing

### Deployment

## Search Projects, Services, Users and Images

## Monitoring Dashboard

## Administrator Options

## Q&A

