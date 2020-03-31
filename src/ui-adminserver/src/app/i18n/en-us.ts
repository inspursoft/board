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
    Node_List_Status_Unknown: 'Unknown',
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
    Node_Detail_Title_Log: 'Node log',
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
      HOSTNAME: {
        NAME: 'Hostname',
        TIPS: 'The IP address or hostname to access admin UI and Kubernetes Cluster apiserver. DO NOT use localhost or 127.0.0.1, because Board needs to be accessed by external clients.'
      },
      API_SERVER_PORT: {
        NAME: 'Api server port',
        TIPS: 'Port of the Board service deployment'
      },
      KUBE_HTTP_SCHEME: {
        NAME: 'Kube http scheme',
        TIPS: 'Kubernetes (K8s) deployment uses network protocols that support http and https'
      },
      KUBE_MASTER_IP: {
        NAME: 'Kube master IP',
        TIPS: 'Host address of Kubernetes (K8s) deployment, *Do not* use localhost or 127.0.0.1 because K8s need to be accessed by other nodes'
      },
      KUBE_MASTER_PORT: {
        NAME: 'Kube master port',
        TIPS: 'The port number deployed by Kubernetes (K8s). If the network protocol is http, port 8080 is used; if https is used, port 6443 is used.'
      },
      REGISTRY_IP: {
        NAME: 'Registry IP',
        TIPS: 'The IP of the image registry'
      },
      REGISTRY_PORT: {
        NAME: 'Registry port',
        TIPS: 'The port of the image registry'
      },
      IMAGE_BASELINE_TIME: {
        NAME: 'Image baseline time',
        TIPS: 'The baseline time of images, it\'s 2016-01-01 09:00:00 by default'
      },
    },
    GOGITS: {
      GOGITS_HOST_IP: {
        NAME: 'Gogits host IP',
        TIPS: 'The host address of the service deployment, *DO NOT* use localhost or 127.0.0.1 because it needs to be accessed by other nodes'
      },
      GOGITS_HOST_PORT: {
        NAME: 'Gogits host port',
        TIPS: 'Host port of Gogits'
      },
      GOGITS_SSH_PORT: {
        NAME: 'Gogits SSH port',
        TIPS: 'Request port number for SSH (Secure Shell)'
      },
    },
    JENKINS: {
      JENKINS_HOST_IP: {
        NAME: 'Jenkins host IP',
        TIPS: 'The address of the Jenkins master for external access'
      },
      JENKINS_HOST_PORT: {
        NAME: 'Jenkins host port',
        TIPS: 'Port number of the Jenkins master'
      },
      JENKINS_NODE_IP: {
        NAME: 'Jenkins node IP',
        TIPS: 'The address of the Jenkins node'
      },
      JENKINS_NODE_SSH_PORT: {
        NAME: 'Jenkins node SSH port',
        TIPS: 'SSH (Secure Shell) request port number of the Node node'
      },
      JENKINS_NODE_USERNAME: {
        NAME: 'Jenkins node username',
        TIPS: 'User name for logging in to the Node'
      },
      JENKINS_NODE_PASSWORD: {
        NAME: 'Jenkins node password',
        TIPS: 'Password for logging in to the Node host'
      },
      JENKINS_NODE_VOLUME: {
        NAME: 'Jenkins node volume',
        TIPS: 'The address of the data volume of the Node'
      },
      JENKINS_EXECUTION_MODE: {
        NAME: 'Jenkins execution mode',
        TIPS: 'Jenkins cluster operation mode, support single node (single) or multi-node (multi)'
      },
    },
    KVM: {
      KVM_REGISTRY_SIZE: {
        NAME: 'Kvm registry size',
        TIPS: ''
      },
      KVM_REGISTRY_PORT: {
        NAME: 'Kvm registry port',
        TIPS: ''
      },
      KVM_REGISTRY_PATH: {
        NAME: 'Kvm registry path',
        TIPS: ''
      },
    },
    LDAP: {
      LDAP_URL: {
        NAME: 'URL',
        TIPS: 'LDAP service entry address'
      },
      LDAP_BASEDN: {
        NAME: 'Base dn',
        TIPS: 'LDAP verifies which source the user came from, ou is people, dc is mydomain, dc is com'
      },
      LDAP_UID: {
        NAME: 'UID type',
        TIPS: 'Attributes used to match users in the search, which can be uid, cn, email, sAMAccountName or other attributes, depending on your LDAP/AD'
      },
      LDAP_SCOPE: {
        NAME: 'Scope',
        TIPS: 'Search for the user\'s scope (LDAP_SCOPE_BASE, LDAP_SCOPE_ONELEVEL, LDAP_SCOPE_SUBTREE)'
      },
      LDAP_TIMEOUT: {
        NAME: 'Timeout',
        TIPS: 'The timeout in seconds to connect to the LDAP server. The default (most reasonable) is 5 seconds'
      },
    },
    EMAIL: {
      EMAIL_IDENTITY: {
        NAME: 'Identity',
        TIPS: 'Identity left blank to act as username'
      },
      EMAIL_SERVER: {
        NAME: 'Server IP',
        TIPS: 'Email server address'
      },
      EMAIL_SERVER_PORT: {
        NAME: 'Server port',
        TIPS: 'Email server port'
      },
      EMAIL_USERNAME: {
        NAME: 'Username',
        TIPS: 'Email username, such as admin@mydomain.com'
      },
      EMAIL_PASSWORD: {
        NAME: 'Password',
        TIPS: ''
      },
      EMAIL_FROM: {
        NAME: 'Email from',
        TIPS: ''
      },
      EMAIL_SSL: {
        NAME: 'Whether to enable SSL',
        TIPS: ''
      },
    },
    OTHERS: {
      ARCH_TYPE: {
        NAME: 'Architecture',
        TIPS: 'The default arch for Board is x86_64, also support mips now'
      },
      DATABASE_PASSWORD: {
        NAME: 'Database password',
        TIPS: 'The password for the root user of mysql db, change this before any production use.'
      },
      TOKEN_CACHE_EXPIRE: {
        NAME: 'Token cache expire',
        TIPS: 'The expiration seconds of token stored in cache.'
      },
      TOKEN_EXPIRE: {
        NAME: 'Token expire',
        TIPS: 'The expiration seconds of token.'
      },
      ELASETICSEARCH_MEMORY: {
        NAME: 'Elaseticsearch memory',
        TIPS: 'The max memory(MB) of elasticsearch can use.'
      },
      TILLER_PORT: {
        NAME: 'Tiller port',
        TIPS: 'The helm tiller node port'
      },
      BOARD_ADMIN_PASSWORD: {
        NAME: 'Board admin password',
        TIPS: 'The initial password of Board admin, only works for the first time when Board starts. Change the admin password from UI after launching Board.'
      },
      BOARD_ADMIN_PASSWORD_OLD: {
        NAME: 'Board Admin old password',
        TIPS: 'Please enter the old password of Board Admin for verification. If the verification is successful, you can modify the new password below.',
      },
      BOARD_ADMIN_PASSWORD_NEW: {
        NAME: 'New password',
        TIPS: 'Please enter a new Board Admin password.',
      },
      BOARD_ADMIN_PASSWORD_CONFIRM: {
        NAME: 'Confirm password',
        TIPS: 'Please confirm the new password entered above.',
      },
      AUTH_MODE: {
        NAME: 'Auth mode',
        TIPS: 'By default the auth mode is db_auth, i.e. the credentials are stored in a local database. Set it to ldap_auth if you want to verify a user\'s credentials against an LDAP server. Set it to indata_auth if you want to verify a user\'s credentials against InData integration platform.'
      },
      VERIFICATION_URL: {
        NAME: 'Verification url',
        TIPS: 'External token verification URL as to integrate authorization with another platform.*NOTE:*: This option is only available when auth_mode is set to \'indata_auth\'.'
      },
      REDIRECTION_URL: {
        NAME: 'Redirection url',
        TIPS: 'Specify redirection URL when token is invalid or expired the UI will redirect to.< strong > NOTE:< /strong> This option is only available when auth_mode is set to \'indata_auth\'.'
      },
      AUDIT_DEBUG: {
        NAME: 'Audit debug',
        TIPS: 'Record all operations switch of operation audit function '
      },
      DNS_SUFFIX: {
        NAME: 'DNS suffix',
        TIPS: 'Kubernetes DNS suffix'
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
  },
  CONFIGURATIONPAGE: {
    UPLOAD: 'Upload',
    SEARCH: 'search configuration',
    CURRENT: 'current config',
    TEMPORARY: 'temporary config',
    CURRENTTIP: 'Here shows the CURRENT Board config. You can click the button next to get the temporary configuration.',
    TEMPORARYTIP: 'Here shows the temporary config. You can click the button next to get CURRENT Board configuration.',
    SAVETOSERVER: 'Save to server',
    SAVETOLOCAL: 'Save locally',
    SAVECONFIGURATION: {
      TITLE: 'Whether the configuration takes effect?',
      SAVE: 'Save successfully',
      COMMENT: 'The new configuration has been saved successfully. Is the configuration applied in the board?It takes some time to restart the system for the configuration to take effect. To confirm the effect, enter the administrator account of the host.',
      ACCOUNT: {
        NAME: 'Account',
        PASSWORD: 'Password',
        REQUIRED: 'This field is required'
      },
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
    },
    PAGE_TITLE: {
      WELCOME: 'Welcome to Adminserver',
      UUID: 'Verify UUID',
      DATABASE: 'Config Database',
      SSH: 'Start Database',
      ACCOUNT: 'Initialize Account',
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
      INITIALIZATION: 'Initialization error! Please check that the docker is working normally.',
      UUID: 'Incorrect UUID! Please re-enter UUID!',
      DATABASE: 'Database initialization failed! Please check that the docker is working normally.',
      PASSWORDS_IDENTICAL: 'Passwords are not identical.',
      MAX_CONNECTION: 'Database\'s max connection error!',
      SSH: 'Network error or account error, please try again!',
      ACCOUNT: 'Network error or docker error, please try again!',
    },
    CONTENTS: {
      WELCOME: 'Welcome to Adminserver! Because this is the first time for the initialization process, you need to complete the following configurations to start the system normally.',
      UUID: 'In order to confirm your identity, you need to enter the UUID in the /data/board/secrets folder for the next configuration.',
      DATABASE: 'Requires to configure the database password to initialize the database. This step may take some time.',
      SSH: 'Please enter the account and password of the host machine. Note: The account needs certain permissions to install and run related components. It will take some time to run.The system will not store your account and password.',
      ACCOUNT: 'Initialize the administrator password of adminserver. The admin account will be shared with the board.',
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
    },
    MODAL_TITLE: 'Database Error!',
  },
};
