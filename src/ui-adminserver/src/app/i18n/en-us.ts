export const LANG_EN_US = {
  BUTTON: {
    OK: 'OK',
    YES: 'Yes',
    NO: 'No',
    CONFIRM: 'Confirm',
    CANCEL: 'Cancel',
    DELETE: 'Delete',
    CREATE: 'Create',
    NEXT: 'Next',
    BACK: 'Back',
    TEST: 'Test',
  },
  DATAGRID: {
    ITEMS: 'items',
    TOTAL: ' total:'
  },
  GLOBAL_ALERT: {
    TITLE: 'Title',
    HINT: 'Hint',
    ASK: 'Ask',
    DELETE: 'Delete',
    SIGN_IN: 'Sign In',
    WARNING: 'Warning',
    ERROR_DETAIL: 'Error Detail'
  },
  ERROR: {
    METHOD_NOT_ALLOWED: 'Method not allowed.',
    CONFLICT_INPUT: 'A conflict with the current state of the resource.',
    INPUT_ONLY_NUMBER: 'Only input number',
    INPUT_REQUIRED: 'This field is required',
    INPUT_NOT_REPEAT: 'This field can not repeat',
    INPUT_MAX_LENGTH: 'Max length',
    INPUT_MIN_LENGTH: 'Min length',
    INPUT_MAX_VALUE: 'Max value',
    INPUT_MIN_VALUE: 'Min value',
    INPUT_PATTERN: 'Input does not confirm to the rules',
    HTTP_504: 'Gateway Timeout:504',
    HTTP_500: 'Unexpected internal error:500',
    HTTP_502: 'Bad gateway :502',
    HTTP_400: 'Bad Request:400',
    HTTP_401: 'User\'s token has expired, please sign in again:401',
    HTTP_403: 'Insufficient privilege to requested operation:403',
    HTTP_404: 'Resource not found:404',
    HTTP_412: 'Precondition given in one or more of the request are evaluated to false:412',
    HTTP_422: 'Unprocessable Entity:422',
    HTTP_UNK: 'Unknown error',
    HTTP_TIME_OUT: 'Http request timeout',
    INSUFFICIENT_PRIVILEGE: 'Insufficient privilege to operate it.'
  },
  HEAD_NAV: {
    DASHBOARD: 'Dashboard',
    CONFIGURATIONS: 'Configurations',
    RESOURCES: 'Resources',
    LANG_EN_US: 'English',
    LANG_ZH_CN: '中文',
    LOGOUT: 'Logout',
  },
  DASHBOARD: {
    BUTTONS: {
      RESTART_BORAD: 'Restart Board',
      STOP_BORAD: 'Stop Board',
      DETAIL: 'Detial',
      RESTART_CONTAINER: 'Restart Container',
      STOP_CONTAINER: 'Stop Container',
    },
    TITLES: {
      SYSTEM_INFO: 'System Info:',
      SYSTEM_CONTENT: 'NULL',
      CONTAINER_INFO: 'Container Info:',
      CONTAINER_CONTENT: 'Containers are Running',
    },
    CONTAINER_INFO: {
      NAME: 'NAME: ',
      ID: 'ID: ',
      IMAGE: 'Image: ',
      CREATED: 'Created: ',
      STATUS: 'Status: ',
      PORTS: 'Ports: ',
      CPU_RATE: 'CPU Rate: ',
      MEMORY_USAGE: 'Memory Usage: ',
      MEMORY_RATE: 'Memory Rate: ',
      NETWORK_IO: 'Network I/O: ',
      BLOCK_IO: 'Block I/O: ',
      PID: 'PID: ',
    },
    CONFIRM: {
      RESTART_BORAD: 'Are you sure to RESTART the Board? If so, please enter the account and password of the host machine:',
      RESTART_CONTAINER: 'Please enter the account and password of the host machine to Restart the Container:',
      STOP_BORAD: 'Are you sure to STOP the Board? If so, please enter the account and password of the host machine:',
    },
  },
  Resource: {
    Nav_Compute: 'Compute resource',
    Nav_Storage: 'Storage resource',
    Nav_Network: 'Network resource',
    Nav_Node_List: 'Node list',
    Nav_Node_Log_List: 'Node operating log list',
  },
  Node: {
    Node_List_Add: 'Add node',
    Node_List_Empty: 'Node list is empty',
    Node_List_Remove_Node: 'Remove node',
    Node_List_Remove_Ask: 'Are you sure to remove the node?',
    Node_List_Status_Unknown: 'Unknown node',
    Node_List_Status_Schedulable: 'Schedulable',
    Node_List_Status_Unschedulable: 'Unschedulable',
    Node_Logs_Ip: 'IP',
    Node_Logs_Operation_Time: 'Operation time',
    Node_Logs_PID: 'PID',
    Node_Logs_Type: 'Add/Remove',
    Node_Logs_Success: 'Success',
    Node_Logs_Operation: 'Operation',
    Node_Logs_Empty: 'Node log list is empty',
    Node_Logs_Delete_Ask: 'Are you sure to delete this log?',
    Node_Logs_Delete_Log_Success: 'Deleted log successfully.',
    Node_Logs_Delete_Log_In_Use: 'This log has been in use, can not be deleted!',
    Node_Logs_Stop_Ask: 'Are you sure to stop this action?',
    Node_Logs_Stop_Success: 'Stop it successfully.',
    Node_Logs_Stop: 'Stop',
    Node_Logs_Removing: 'Removing...',
    Node_Logs_Adding: 'Adding...',
    Node_Detail_Add: 'Add',
    Node_Detail_Remove: 'Remove',
    Node_Detail_Refresh: 'Refresh log',
    Node_Detail_Title_Add: 'Add node',
    Node_Detail_Title_Remove: 'Removing node',
    Node_Detail_Title_Log: 'Add node log',
    Node_Detail_Node_Ip: 'Node ip',
    Node_Detail_Node_Password: 'Node password',
    Node_Detail_Remove_Success: 'Remove node successfully.',
    Node_Detail_Add_Success: 'Add node successfully.',
    Node_Detail_Remove_Failed: 'Remove node failed.',
    Node_Detail_Add_Failed: 'Add node failed.',
    Node_Detail_Master_Title: 'Master {{0}} password',
    Node_Detail_Master_Hint: 'Please input the password of master.',
    Node_Detail_Host_Username: 'Host {{0}} username',
    Node_Detail_Host_Password: 'Host {{0}} password',
    Node_Detail_Host_Username_Hint: 'Please input the username of the host.',
    Node_Detail_Host_Password_Hint: 'Please input the password of the host.',
    Node_Detail_Node_Ip_Hint: 'Please input the node ip.',
    Node_Detail_Node_Password_Hint: 'Please input the password of the node.'
  },
  CONFIGURATIONS: {
    API_SERVER: {
      TITLE: 'Api Server',
      HOSTNAME: {
        NAME: 'Hostname',
        PLACEHOLDER: 'eg：192.168.1.10 or reg.yourdomain.com',
        TIPS: 'The IP address or hostname to access admin UI and Kubernetes Cluster apiserver. DO NOT use localhost or 127.0.0.1, because Board needs to be accessed by external clients.'
      },
      API_SERVER_PORT: {
        NAME: 'Api server port',
        PLACEHOLDER: '1 ~ 65535, default: 8088',
        TIPS: 'Port of the Board service deployment'
      },
      KUBE_HTTP_SCHEME: {
        NAME: 'Kube http scheme',
        PLACEHOLDER: 'Kube http scheme',
        TIPS: 'Kubernetes (K8s) deployment uses network protocols that support http and https'
      },
      KUBE_MASTER_IP: {
        NAME: 'Kube master IP',
        PLACEHOLDER: 'eg: 192.168.1.10',
        TIPS: 'Host address of Kubernetes (K8s) deployment, Do not use localhost or 127.0.0.1 because K8s need to be accessed by other nodes'
      },
      KUBE_MASTER_PORT: {
        NAME: 'Kube master port',
        PLACEHOLDER: '',
        TIPS: 'Port of Kubernetes (K8s). By default, http uses port 8080 and https uses port 6443. Please fill in according to the actual situation under non-default conditions.'
      },
      REGISTRY_IP: {
        NAME: 'Registry IP',
        PLACEHOLDER: 'eg: 192.168.1.10',
        TIPS: 'The address of the image registry, do not use localhost or 127.0.0.1, because the image registry needs other nodes to access'
      },
      REGISTRY_PORT: {
        NAME: 'Registry port',
        PLACEHOLDER: '1~65535, default: 5000',
        TIPS: 'The port of the image registry, default: 5000'
      },
      IMAGE_BASELINE_TIME: {
        NAME: 'Image baseline time',
        PLACEHOLDER: 'default: 2016-01-01 09:00:00',
        TIPS: 'The baseline time of images, it\'s 2016-01-01 09:00:00 by default. Board only reads the image after this time'
      },
    },
    GOGITS: {
      TITLE: 'Gogs',
      GOGITS_HOST_IP: {
        NAME: 'Gogs host IP',
        PLACEHOLDER: 'eg: 192.168.1.10',
        TIPS: 'The host address of Gogs, DO NOT use localhost or 127.0.0.1 because it needs to be accessed by other nodes'
      },
      GOGITS_HOST_PORT: {
        NAME: 'Gogs host port',
        PLACEHOLDER: '1~65535, default: 10080',
        TIPS: 'Host port of Gogs, default: 10080'
      },
      GOGITS_SSH_PORT: {
        NAME: 'Gogs SSH port',
        PLACEHOLDER: '1~65535, default: 10022',
        TIPS: 'Request port number for SSH (Secure Shell), default: 10022'
      },
    },
    JENKINS: {
      TITLE: 'Jenkins',
      JENKINS_HOST_IP: {
        NAME: 'Jenkins host IP',
        PLACEHOLDER: 'eg: 192.168.1.10',
        TIPS: 'The address of the Jenkins master for external access. Don\'t use localhost or 127.0.0.1'
      },
      JENKINS_HOST_PORT: {
        NAME: 'Jenkins host port',
        PLACEHOLDER: '1~65535, default: 8888',
        TIPS: 'Port number of the Jenkins master, default: 8888'
      },
      JENKINS_NODE_IP: {
        NAME: 'Jenkins node IP',
        PLACEHOLDER: 'eg: 192.168.1.11',
        TIPS: 'The address of the Jenkins node'
      },
      JENKINS_NODE_SSH_PORT: {
        NAME: 'Jenkins node SSH port',
        PLACEHOLDER: '1~65535, default: 22',
        TIPS: 'SSH (Secure Shell) request port number of the Node node, default: 22'
      },
      JENKINS_NODE_USERNAME: {
        NAME: 'Jenkins node username',
        PLACEHOLDER: 'Jenkins node username',
        TIPS: 'User name for logging in to the Node'
      },
      JENKINS_NODE_PASSWORD: {
        NAME: 'Jenkins node password',
        PLACEHOLDER: '8~20 digits, numbers/letters/#?!@$%^&*-',
        TIPS: 'The password of the Jemkins account, please use this password before production. The password must contain large and small letters, numbers, 8~20 characters, support special characters #?!@$%^&*-',
        PATTERN_ERROR: 'The password must contain large and small letters, numbers, 8~20 characters, support special characters #?!@$%^&*-',
      },
      JENKINS_NODE_VOLUME: {
        NAME: 'Jenkins node volume',
        PLACEHOLDER: 'eg: /data/jenkins_node',
        TIPS: 'The address of the data volume of the Node'
      },
      JENKINS_EXECUTION_MODE: {
        NAME: 'Jenkins execution mode',
        PLACEHOLDER: '',
        TIPS: 'Jenkins cluster operation mode, support single mode (single) or multi mode (multi)'
      },
    },
    KVM: {
      TITLE: 'KVM',
      KVM_REGISTRY_SIZE: {
        NAME: 'Kvm registry size',
        PLACEHOLDER: 'default: 5',
        TIPS: 'Need to be greater than 0, the default is 5'
      },
      KVM_REGISTRY_PORT: {
        NAME: 'Kvm registry port',
        PLACEHOLDER: '1~65535, default: 8890',
        TIPS: '1~65535, default: 8890'
      },
      KVM_REGISTRY_PATH: {
        NAME: 'Kvm registry path',
        PLACEHOLDER: 'eg: /root/kvm_toolkits',
        TIPS: 'KVM tool storage path, default is /root/kvm_toolkits'
      },
    },
    LDAP: {
      TITLE: 'LDAP(Lightweight Directory Access Protocol)',
      LDAP_URL: {
        NAME: 'URL',
        PLACEHOLDER: 'eg: ldaps://ldap.mydomain.com',
        TIPS: 'LDAP service entry address'
      },
      LDAP_BASEDN: {
        NAME: 'Base dn',
        PLACEHOLDER: 'ou=people,dc=mydomain,dc=com',
        TIPS: 'LDAP verifies which source the user came from, ou is people, dc is mydomain, dc is com'
      },
      LDAP_UID: {
        NAME: 'UID type',
        PLACEHOLDER: 'uid/cn/email/...',
        TIPS: 'Attributes used to match users in the search, which can be uid, cn, email, sAMAccountName or other attributes, depending on your LDAP/AD'
      },
      LDAP_SCOPE: {
        NAME: 'Scope',
        PLACEHOLDER: '',
        TIPS: 'Search for the user\'s scope, default: Subtree'
      },
      LDAP_TIMEOUT: {
        NAME: 'Timeout',
        PLACEHOLDER: 'The unit is seconds, default: (the most reasonable)5',
        TIPS: 'The timeout in seconds to connect to the LDAP server. The unit is seconds. The default (most reasonable) is 5 seconds'
      },
    },
    EMAIL: {
      TITLE: 'Email',
      EMAIL_IDENTITY: {
        NAME: 'Identity',
        PLACEHOLDER: 'Identity',
        TIPS: 'The default is empty identity(NULL)'
      },
      EMAIL_SERVER: {
        NAME: 'Server IP',
        PLACEHOLDER: 'eg: smtp.mydomain.com',
        TIPS: 'Email server address'
      },
      EMAIL_SERVER_PORT: {
        NAME: 'Server port',
        PLACEHOLDER: '1~65535, default: 25',
        TIPS: 'Email server port, default: 25'
      },
      EMAIL_USERNAME: {
        NAME: 'Username',
        PLACEHOLDER: 'eg: admin@mydomain.com',
        TIPS: 'Email username, eg: admin@mydomain.com'
      },
      EMAIL_PASSWORD: {
        NAME: 'Password',
        PLACEHOLDER: 'Email password',
        TIPS: 'Used to login to the email and send mail'
      },
      EMAIL_FROM: {
        NAME: 'Email from',
        PLACEHOLDER: 'eg: admin <admin@mydomain.com>',
        TIPS: 'eg: admin <admin@mydomain.com>'
      },
      EMAIL_SSL: {
        NAME: 'Whether SSL',
        PLACEHOLDER: '',
        TIPS: 'Whether to enable SSL to encrypt the network connection, the default is false'
      },
    },
    OTHERS: {
      TITLE: 'Initialization',
      ARCH_TYPE: {
        NAME: 'Architecture',
        PLACEHOLDER: '',
        TIPS: 'The default arch for Board is x86_64, also support mips and arm64v8 now'
      },
      SECURITY_MODE: {
        NAME: 'Security mode',
        PLACEHOLDER: '',
        TIPS: 'Security mode of Board, by default it\'s "Normal". "Security" mode will hide Kibana, Grafana in UI',
      },
      DATABASE_PASSWORD: {
        NAME: 'Database password',
        PLACEHOLDER: '8~20 digits, numbers/letters/#?!@$%^&*-',
        TIPS: 'The password for the root user of mysql db, change this before any production use. The password must contain large and small letters, numbers, 8~20 characters, support special characters #?!@$%^&*-',
        PATTERN_ERROR: 'The password must contain large and small letters, numbers, 8~20 characters, support special characters #?!@$%^&*-',
      },
      DB_MAX_CONNECTION: {
        NAME: 'DB Max Connection',
        PLACEHOLDER: '1~16384, recommended: 1000',
        TIPS: '1~16384, recommended: 1000',
      },
      TOKEN_CACHE_EXPIRE: {
        NAME: 'Token cache expire',
        PLACEHOLDER: 'In seconds, recommended: 1800',
        TIPS: 'The expiration seconds of token stored in cache, recommended: 1800'
      },
      TOKEN_EXPIRE: {
        NAME: 'Token expire',
        PLACEHOLDER: 'In seconds, recommended: 1800',
        TIPS: 'The expiration seconds of token, recommended: 1800'
      },
      ELASETICSEARCH_MEMORY: {
        NAME: 'Elaseticsearch memory',
        PLACEHOLDER: 'In MB, default: 1024',
        TIPS: 'The max memory(MB) of elasticsearch can use, default: 1024'
      },
      TILLER_PORT: {
        NAME: 'Tiller port',
        PLACEHOLDER: '1~65535, default: 31111',
        TIPS: 'The helm tiller node port, default: 31111'
      },
      BOARD_ADMIN_PASSWORD: {
        NAME: 'Board admin password',
        PLACEHOLDER: '8~20 digits, numbers/letters/#?!@$%^&*-',
        TIPS: 'The initial password of Admin only takes effect when Board is started for the first time. This password is shared between the admin account of the Board and Adminserver. To modify it, you need to configure Email parameters, and then forget the password after the Board starts successfully. The password must contain large and small letters, numbers, 8~20 characters, support special characters #?!@$%^&*-',
        PATTERN_ERROR: 'The password must contain large and small letters, numbers, 8~20 characters, support special characters #?!@$%^&*-',
      },
      BOARD_ADMIN_PASSWORD_OLD: {
        NAME: 'Board Admin old password',
        PLACEHOLDER: 'Board Admin old password',
        TIPS: 'Please enter the old password of Board Admin for verification. If the verification is successful, you can modify the new password below.',
      },
      BOARD_ADMIN_PASSWORD_NEW: {
        NAME: 'New password',
        PLACEHOLDER: 'New password',
        TIPS: 'Please enter a new Board Admin password.',
      },
      BOARD_ADMIN_PASSWORD_CONFIRM: {
        NAME: 'Confirm password',
        PLACEHOLDER: 'Confirm password',
        TIPS: 'Please confirm the new password entered above.',
      },
      AUTH_MODE: {
        NAME: 'Auth mode',
        PLACEHOLDER: 'Auth mode',
        TIPS: 'By default the auth mode is \'Database\', the credentials are stored in a local database. Set it to \'LDAP\' if you want to verify a user\'s credentials against an LDAP server. Set it to \'Indata\' if you want to verify a user\'s credentials against InData integration platform.'
      },
      VERIFICATION_URL: {
        NAME: 'Verification url',
        PLACEHOLDER: 'eg: http://verification.mydomain.com',
        TIPS: 'External token verification URL as to integrate authorization with another platform. NOTE: This option is only available when auth_mode is set to \'Indata\'.'
      },
      REDIRECTION_URL: {
        NAME: 'Redirection url',
        PLACEHOLDER: 'eg: http://redirection.mydomain.com',
        TIPS: 'Specify redirection URL when token is invalid or expired the UI will redirect to. NOTE: This option is only available when auth_mode is set to \'Indata\'.'
      },
      AUDIT_DEBUG: {
        NAME: 'Audit debug',
        PLACEHOLDER: '',
        TIPS: 'Record all operations switch of operation audit function, default: false'
      },
      DNS_SUFFIX: {
        NAME: 'DNS suffix',
        PLACEHOLDER: 'eg: .cluster.local',
        TIPS: 'Kubernetes DNS suffix, default: .cluster.local'
      },
    },
  },
  ACCOUNT: {
    SIGN_IN: 'Sign In',
    USERNAME_PLACEHOLDER: 'Username',
    PASSWORD_PLACEHOLDER: 'Password',
    SIGN_UP: 'Sign Up',
    FORGOT_PASSWORD: 'Forgot Password',
    FORGOT_PASSWORD_HINT_MSG: 'Please input username/email',
    FORGOT_PASSWORD_ERROR_MSG: 'Username/email can\'t empty',
    USERNAME: 'Username',
    REALNAME: 'Real name',
    INPUT_USERNAME: 'Input username',
    EMAIL: 'Email',
    INPUT_EMAIL: 'Input email',
    PASSWORD: 'Password',
    INPUT_PASSWORD: 'Input password',
    CONFIRM_PASSWORD: 'Confirm Password',
    INPUT_CONFIRM_PASSWORD: 'Input confirm password',
    COMMENT: 'Comment',
    REQUIRED_ITEMS: 'Required items.',
    BACK: 'Back',
    REGISTER: 'Register',
    USERNAME_IS_REQUIRED: 'Username is required.',
    USERNAME_ALREADY_EXISTS: 'Username already exists.',
    USERNAME_IS_KEY: 'The username is key.',
    USERNAME_ARE_NOT_IDENTICAL: 'The username consists of numeric, lowercase letter and underscore, and the length range is [4,40]',
    EMAIL_IS_REQUIRED: 'Email is required.',
    EMAIL_IS_ILLEGAL: 'Email is illegal.',
    EMAIL_ALREADY_EXISTS: 'Email already exists',
    PASSWORD_IS_REQUIRED: 'Password is required.',
    PASSWORD_FORMAT: 'Password should be at least 8 characters with at least 1 uppercase, 1 lowercase and 1 number.',
    PASSWORDS_ARE_NOT_IDENTICAL: 'Passwords are not identical.',
    ERROR: 'Error',
    INCORRECT_USERNAME_OR_PASSWORD: 'Incorrect username or password.',
    SUCCESS_TO_SIGN_IN: 'Successful sign in',
    FAILED_TO_SIGN_IN: 'Failed to sign in:',
    FAILED_TO_SIGN_UP: 'Failed to sign up:',
    SUCCESS_TO_SIGN_UP: 'User sign up successfully.',
    FAILED_TO_SIGN_OUT: 'Failed to sign out.',
    ACCOUNT_SETTING_SUCCESS: 'Account setting success.',
    ALREADY_SIGNED_IN: 'The user has already signed in other place.',
    SEND_REQUEST: 'Send Request',
    SEND_REQUEST_SUCCESS: 'Send Request Success',
    SEND_REQUEST_SUCCESS_MSG: 'Send Request Success,Please Check Email.',
    RESET_NEW_PASS: 'Set New Password',
    USER_NOT_EXISTS: 'User Not Exist.',
    SEND_REQUEST_ERR: 'Send Request Error',
    SEND_REQUEST_ERR_MSG: 'Failed to send validation, please contact administrator.',
    RESET_PASS_SUCCESS: 'Reset Password Success',
    RESET_PASS_SUCCESS_MSG: 'Reset Password Success,Please Login',
    RESET_PASS_ERR: 'Reset Password Error',
    RESET_PASS_ERR_MSG: 'Reset Password Error,Please Contact Administrator.',
    INVALID_RESET_UUID: 'The Link is Invalid.',
    REGISTER_INFO: 'The system has just been initialized and you need to set up an account & password before using it.',
    FORBIDDEN: 'Forbidden access!',
    FORGOT_PASSWORD_HELPER: 'If you have configured an email, please change the password in the Board. If you have not configured an email, please contact the administrator for assistance.',
    TOKEN_ERROR: 'User status error! Please login again!',
  },
  CONFIGURATIONPAGE: {
    UPLOAD: 'Upload',
    SEARCH: 'search configuration',
    CURRENT: 'current config',
    TEMPORARY: 'temporary config',
    CURRENTTIP: 'Here shows the CURRENT Board config. You can click the button next to get the temporary configuration.',
    TEMPORARYTIP: 'Here shows the temporary config. You can click the button next to get CURRENT Board configuration.',
    HEADER_HELPER: 'Only allow viewing of configuration.To update the configuration, you need to reconfigure after stopping the Board on the "Dashboard" page.',
    SAVETOSERVER: 'Save to server',
    SAVETOLOCAL: 'Save locally',
    SAVECONFIGURATION: {
      TITLE: 'Whether the configuration takes effect?',
      SAVE: 'Save successfully',
      COMMENT: 'The new configuration has been saved successfully. Is the configuration applied in the Board?It takes some time to restart the system for the configuration to take effect. To confirm the effect, enter the administrator account of the host.',
      ACCOUNT: {
        NAME: 'Account',
        PASSWORD: 'Password',
        REQUIRED: 'This field is required'
      },
      HELPER: 'Please make sure you have been stop Board instance.',
      CANCEL: 'Not now',
      APPLY: 'Yes, Do it',
    }
  },
  INITIALIZATION: {
    TITLE: 'Adminserver Initialization',
    BUTTONS: {
      NEXT: 'Next',
      TRANSLATE: '切换为汉语',
      APPLY: 'Apply',
      RESTART_DB: 'Restart DB',
      FAST_MODE: 'Quick mode',
      EDIT_CONFIG: 'Edit cfg',
      START_BOARD: 'Start Board',
      UNINSTALL: 'Uninstall Board',
      APPLY_AND_START_BOARD: 'Apply & Start Board',
      GO_TO_BOARD: 'Go to Board',
      GO_TO_ADMINSERVER: 'Go to Adminserver',
      REINSTALL: 'Re-install Board',
    },
    PAGE_TITLE: {
      WELCOME: 'Welcome to Adminserver',
      UUID: 'Verify UUID',
      DATABASE: 'Config Database',
      SSH: 'Start Database',
      ACCOUNT: 'Initialize Account',
      EDIT_CONFIG_CONFIRM: 'Do you want to edit the configuration?',
      EDIT_CONFIG: 'Edit configuration of Board',
      FINISH: 'Congratulations!',
      UNINSTALL: 'Uninstallation is complete!',
      SSH_ACCOUNT: 'Please input the SSH account',
    },
    PAGE_NAV_TITLE: {
      WELCOME: 'Welcome',
      UUID: 'Verify UUID',
      DATABASE: 'Config Database',
      SSH: 'Start Database',
      ACCOUNT: 'Initialize Account',
    },
    TOOLS: {
      LOADING: 'Loading...',
    },
    ALERTS: {
      INITIALIZATION: 'Initialization error! Please check that the docker or service is working normally.',
      UUID: 'UUID format error! Please re-enter UUID!',
      DATABASE: 'Database initialization failed! Please check that the docker is working normally.',
      PASSWORDS_IDENTICAL: 'Passwords are not identical.',
      MAX_CONNECTION: 'Database\'s max connection error!',
      SSH: 'Network error or account error, please try again!',
      ACCOUNT: 'Network error or docker error, please try again!',
      GET_SYS_STATUS_FAILED: 'Failed to get system status. Please check whether the service is running normally.',
      VALIDATE_UUID_FAILED: 'UUID verification failed! Please check whether the network is normal or the UUID is filled in correctly.',
      GET_CFG_FAILED: 'Failed to get the configuration. Please check whether the service is running normally.',
      GET_TMP_FAILED: 'Failed to get temporary configuration, use current configuration.',
      START_BOARD_FAILED: 'Failed to start Board. Please check whether the configuration is correct or the service is running normally.',
      UNINSTALL_BOARD_FAILED: 'Failed to uninstall Board, please check whether the service is running normally',
      POST_CFG_FAILED: 'Failed to save the configuration. Please check whether the service is running normally.',
      ALREADY_START: 'Board have been successfully started. The next steps cannot be operated.',
    },
    CONTENTS: {
      WELCOME: 'Welcome to Adminserver! Because this is the first time for the initialization process, you need to complete the following configurations to start the system normally.',
      UUID: 'In order to confirm your identity, you need to enter the UUID in the /data/board/secrets folder for the next.',
      DATABASE: 'Requires to configure the database password to initialize the database. This step may take some time.',
      SSH: 'Please enter the account and password of the host machine. Note: The account needs certain permissions to install and run related components. It will take some time to run.The system will not store your account and password.',
      ACCOUNT: 'Initialize the administrator password of Adminserver. The admin account will be shared with the Board.',
      EDIT_CONFIG_CONFIRM: 'It seems that there is already a configured cfg, do you want to re-edit it or start the Board directly?',
      FINISH: 'Board have been successfully started. Initialization needs to wait for a while, after waiting, you can access the following connection to access the Board or Adminserver.',
      UNINSTALL: 'Board component uninstallation is complete!You can now go to the background to uninstall Adminserver and clear related data.',
      CLEAR_DATA: 'Clear all relevant data of Board',
      RESPONSIBILITY: 'I already know the impact of this operation and I am responsible for its consequences.',
    },
    LABELS: {
      UUID: 'UUID',
      MAX_CONNECTION: 'Max Connection',
      PASSWORD: 'Password',
      PASSWORD_CONFIRM: 'Confirm Password',
      ACCOUNT: 'Account',
    },
    HELPER: {
      REQUIRED: 'This field is required!',
      NUMS: 'Requires 10 ~ 16384.',
      MIN_LENGTH_8: 'It must be at least 8 characters!',
      NOT_EDITABLE: 'Unable to edit!',
      FAST_MODE: 'Use default values ​​for some configurations',
    },
    MODAL_TITLE: 'Database Error!',
  },
};
