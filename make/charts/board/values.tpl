image:
  registry: "10.111.24.131:5000"
  tag: "7.4-150-g11a0d2d"
localtime:
  path: /etc/localtime
apiserver:
  name: "apiserver"
  replicaCount: 1
  image:
    repository: $registry/board_apiserver
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 8088
  resources:
    limits:
     cpu: 100m
     memory: 128Mi
    requests:
     cpu: 100m
     memory: 128Mi
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
chartmuseum:
  name: "chartmuseum"
  replicaCount: 1
  image:
    repository: $registry/board_chartmuseum
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 8080
  resources: {}
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
  persistence:
    enabled: true
    existingClaim: ""
    ## database data Persistent Volume Storage Class
    ## If defined, storageClassName: <storageClass>
    ## If set to "-", storageClassName: "", which disables dynamic provisioning
    ## If undefined (the default) or set to null, no storageClassName spec is
    ##   set, choosing the default provisioner.  
    # storageClass: "-"
    accessMode: ReadWriteOnce
    size: 8Gi
    volumeName: chartmuseum
db:
  name: "db"
  replicaCount: 1
  image:
    repository: $registry/board_db
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 3306
  resources: {}
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
  persistence:
    enabled: true
    existingClaim: ""
    ## database data Persistent Volume Storage Class
    ## If defined, storageClassName: <storageClass>
    ## If set to "-", storageClassName: "", which disables dynamic provisioning
    ## If undefined (the default) or set to null, no storageClassName spec is
    ##   set, choosing the default provisioner.  
    # storageClass: "-"
    accessMode: ReadWriteOnce
    size: 8Gi
    volumeName: db
elasticsearch:
  name: "elasticsearch"
  replicaCount: 1
  image:
    repository: $registry/board_elasticsearch
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: NodePort
    port: 9200
  resources: {}
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
  persistence:
    enabled: true
    existingClaim: ""
    ## database data Persistent Volume Storage Class
    ## If defined, storageClassName: <storageClass>
    ## If set to "-", storageClassName: "", which disables dynamic provisioning
    ## If undefined (the default) or set to null, no storageClassName spec is
    ##   set, choosing the default provisioner.  
    # storageClass: "-"
    accessMode: ReadWriteOnce
    size: 8Gi
    volumeName: elasticsearch
grafana:
  name: "grafana"
  replicaCount: 1
  image:
    repository: $registry/board_grafana
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 3000
  resources: {}
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
  persistence:
    enabled: true
    lib:
      existingClaim: ""
      ## database data Persistent Volume Storage Class
      ## If defined, storageClassName: <storageClass>
      ## If set to "-", storageClassName: "", which disables dynamic provisioning
      ## If undefined (the default) or set to null, no storageClassName spec is
      ##   set, choosing the default provisioner.  
      # storageClass: "-"
      accessMode: ReadWriteOnce
      size: 8Gi
      volumeName: grafana-lib
    log:
      existingClaim: ""
      ## database data Persistent Volume Storage Class
      ## If defined, storageClassName: <storageClass>
      ## If set to "-", storageClassName: "", which disables dynamic provisioning
      ## If undefined (the default) or set to null, no storageClassName spec is
      ##   set, choosing the default provisioner.  
      # storageClass: "-"
      accessMode: ReadWriteOnce
      size: 8Gi
      volumeName: grafana-log
kibana:
  name: "kibana"
  replicaCount: 1
  image:
    repository: $registry/board_kibana
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 5601
  resources: {}
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
prometheus:
  name: "prometheus"
  replicaCount: 1
  image:
    repository: $registry/board_prometheus
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 9090
  resources: {}
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
  persistence:
    enabled: true
    existingClaim: ""
    ## database data Persistent Volume Storage Class
    ## If defined, storageClassName: <storageClass>
    ## If set to "-", storageClassName: "", which disables dynamic provisioning
    ## If undefined (the default) or set to null, no storageClassName spec is
    ##   set, choosing the default provisioner.  
    # storageClass: "-"
    accessMode: ReadWriteOnce
    size: 8Gi
    volumeName: prometheus
proxy:
  name: "proxy"
  replicaCount: 1
  image:
    repository: $registry/board_proxy
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: NodePort
    httpPort: 80
    httpsPort: 443
    proxyPort: 8080
  resources: {}
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
tokenserver:
  name: "tokenserver"
  replicaCount: 1
  image:
    repository: $registry/board_tokenserver
    tag: "$tag"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    port: 4000
  resources: {}
  nodeSelector: {}
  tolerations: {}
  affinity: {}
  restartPolicy: Always
native-elasticsearch:
  image: $registry/elasticsearch/elasticsearch
  extraInitContainers:
  - command:
    - chmod
    - -R
    - "777"
    - /usr/share/elasticsearch/data
    image: $registry/elasticsearch/elasticsearch:7.9.3
    imagePullPolicy: IfNotPresent
    name: chmod
    resources: {}
    securityContext:
      privileged: true
      runAsUser: 0
    volumeMounts:
    - mountPath: /usr/share/elasticsearch/data
      name: elasticsearch-master
native-fluentd-elasticsearch:
  image:
    repository: $registry/fluentd_elasticsearch/fluentd
  hostLogDir:
    dockerContainers: $dockercontainers
native-kibana:
  image: $registry/kibana/kibana