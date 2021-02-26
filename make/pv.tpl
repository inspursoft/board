apiVersion: v1
items:
- apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: board-prometheus
  spec:
    accessModes:
    - ReadWriteOnce
    capacity:
      storage: 8Gi
    nfs:
      path: __nfs_path__/board-prometheus
      server: __nfs_server__
    persistentVolumeReclaimPolicy: Retain
    volumeMode: Filesystem
- apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: chartmuseum
  spec:
    accessModes:
    - ReadWriteOnce
    capacity:
      storage: 8Gi
    nfs:
      path: __nfs_path__/chartmuseum
      server: __nfs_server__
    persistentVolumeReclaimPolicy: Retain
    volumeMode: Filesystem
- apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: db
  spec:
    accessModes:
    - ReadWriteOnce
    capacity:
      storage: 8Gi
    nfs:
      path: __nfs_path__/db
      server: __nfs_server__
    persistentVolumeReclaimPolicy: Retain
    volumeMode: Filesystem
- apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: elasticsearch
  spec:
    accessModes:
    - ReadWriteOnce
    capacity:
      storage: 8Gi
    nfs:
      path: __nfs_path__/elasticsearch
      server: __nfs_server__
    persistentVolumeReclaimPolicy: Retain
    volumeMode: Filesystem
- apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: elasticsearch-log
  spec:
    accessModes:
    - ReadWriteOnce
    capacity:
      storage: 30Gi
    nfs:
      path: __nfs_path__/elasticsearch-log
      server: __nfs_server__
    persistentVolumeReclaimPolicy: Delete
    volumeMode: Filesystem
- apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: grafana-lib
  spec:
    accessModes:
    - ReadWriteOnce
    capacity:
      storage: 8Gi
    nfs:
      path: __nfs_path__/grafana-lib
      server: __nfs_server__
    persistentVolumeReclaimPolicy: Retain
    volumeMode: Filesystem
- apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: grafana-log
  spec:
    accessModes:
    - ReadWriteOnce
    capacity:
      storage: 8Gi
    nfs:
      path: __nfs_path__/grafana-log
      server: __nfs_server__
    persistentVolumeReclaimPolicy: Retain
    volumeMode: Filesystem
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""
