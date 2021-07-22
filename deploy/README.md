## 部署

主要存放一些部署相关的配置文件和脚本

### 部署Go应用

> see: https://eddycjy.com/posts/kubernetes/2020-05-03-deployment/


## 监控

包含对机器(node或者container)的监控、应用的监控、数据库的监控等。

使用 `docker-compose` 可以在本地一键部署，配置如下：

```yaml

```

## 配置etcd

create a namespace

`$ kubectl create namespace etcd`

create the service 

`$ kubectl apply -f etcd-service.yaml -n etcd`

create the cluster(statefulSet)

`$ cat etcd.yml.tmpl | etcd_config.bash | kubectl apply -n etcd -f -`

Verify the cluster's health

`$ kubectl exec -it etcd-0 -n etcd etcdctl cluster-health`

The cluster is exposed through minikube's IP

```bash
$ IP=$(minikube ip)
$ PORT=$(kubectl get services -o jsonpath="{.spec.ports[].nodePort}" etcd-client -n etcd)
$ etcdctl --endpoints http://${IP}:${PORT} get foo
```

Destroy the services

```bash
$ kubectl delete services,statefulsets --all -n etcd
```

给Etcd集群做个Web UI

```
$ kubectl apply -f etcd-ui-configmap.yaml
$ kubectl apply -f etcd-ui.yaml
```

> Ref: https://github.com/kevinyan815/LearningKubernetes/tree/master/e3w

> 参考： 
> https://mp.weixin.qq.com/s/AkIvkW22dvqcdFXkiTpv8Q
> https://github.com/kevinyan815/LearningKubernetes/tree/master/etcd