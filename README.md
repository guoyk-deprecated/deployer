# deployer

[![BMC Donate](https://img.shields.io/badge/BMC-Donate-orange)](https://www.buymeacoffee.com/vFa5wfRq6)

一个用 Go 编写的命令行工具，用以构建，打包并部署 Docker 镜像到 k8s 环境

## 获取镜像

一般情况下不会使用 Docker 镜像来获取此可执行文件

`guoyk/deployer`

## 前置条件

* `docker` 已安装
* `kubectl` 已安装

## 使用方法

假设集群名为 `test-cluster`, 命名空间名字为 `test-namespace`, 工作负责名为 `test-deployment`, 环境名为 `test`

1. 将 `deployer` 放置在可被 `Jenkins` 调用的位置

2. 准备 `test-cluster.yml` 文件, 放置为 `$HOME/.deployer/preset-test-cluster.yml`, Jenkins 默认 `$HOME` 为 `/var/lib/jenkins`

    文件内容如下

    ```yaml
    registry: 127.0.0.1:5000 #所使用的的镜像仓库基地址
    imagePullSecrets: 
      - test_secret #k8s集群拉取镜像所需的密文名称
    kubeconfig: # kubeconfig 文件内容, 用于部署镜像
      # ...
    dockerconfig: # dockerconfig 文件内容, 用于推送镜像, 参考 $HOME/.docker/config.json
      auth:
        registry.example.com:
          auth: "xxxxxxx"
      # ...
    ```
   
3. 准备 `deployer.yml` 文件, 放置到项目根目录

   文件内容如下

   ```yaml
   default:
     build: # 构建 Bash 脚本, 数组形式
       - echo default
     package: # Dockerfile 打包脚本, 数组形式
       - FROM nginx
       - ADD index.html /usr/share/nginx/html/index.html
   test: # 可以针对不同环境, 覆盖默认的 build 和 package
     build:
       - echo test
   ```

 4. Jenkins 脚本内容

   ```yaml
   deployer --cluster test-cluster --namespace test-namespace --deployment test-deployment --env test
   ```

  将任务名设置为 `test-deployment.test` 来省略 `--deployment test-deployment --env test` 两个选项

  将任务名设置为 `test-cluster.test-namespace.test-deployment.test` 来省略所有选项
  
  其他选项:
  
  ```text
Usage of deployer:
  -cluster string
    	集群, 决定使用哪个集群配置文件, 使用 $HOME/.deployer/preset-[CLUSTER].yml
  -container string
    	容器名称,  默认和 deployment 相同（基于 rancher 习惯）
  -cpu string
    	CPU 资源配置, 格式 '[申请]:[限额]', 如 50:250, 单位 'm' 千分之一核心
  -deployment string
    	工作负载, 要求和 Kubernetes 上的工作负载名称完全一致（已废弃，使用 -workload）
  -env string
    	环境名, 决定 deployer.yml 或者 docker-build.XXX.sh,  Dockerfile.XXX 文件的选择, 和构建后的镜像标签
  -image string
    	镜像基础名称, 默认为 '[命名空间]-[工作负载]'
  -init
    	容器为 init 类型的容器
  -keep-generated
    	保存生成的中间文件, 用于调试
  -keep-image
    	在本地 Docker 保留镜像, 便于重复构建, 或者调试
  -limits-cpu string
    	CPU 资源限制, 单位必须为 'm', 千分之一核心（废弃, 使用 --cpu 参数）
  -limits-mem string
    	MEM 资源限制, 单位必须为 'Mi', 兆字节（废弃, 使用 --mem 参数）
  -mem string
    	内存资源配置, 格式 '[申请]:[限额]', 如 64:256, 单位 'Mi' 兆字节
  -namespace string
    	命名空间, 工作负载在集群中的命名空间
  -only-build
    	只执行 build 步骤, 不进行 package, push, deploy 步骤
  -only-deploy
    	只执行 deploy 步骤, 不进行 build, package, push 步骤
  -registry string
    	镜像仓库, 指定镜像要推往的仓库地址
  -requests-cpu string
    	CPU 资源请求, 单位必须为 'm', 千分之一核心（废弃, 使用 --cpu 参数）
  -requests-mem string
    	MEM 资源请求, 单位必须为 'Mi', 兆字节（废弃, 使用 --mem 参数）
  -skip-deploy
    	只执行 build, package, push 步骤, 不进行 deploy 步骤
  -workload string
    	工作负载名
  -workload-type string
    	工作负载类型，默认为 deployment (default "deployment")
  ```

## deployer.yml

本质上是一个将 构建脚本 (build, 即 bash 脚本) 和 打包脚本 (package, 即 Dockerfile 文件) 合并到一起的 YAML 文件, 如下所示

```yaml
default:
  build:
    - mvn -Pdev -DskipTests=true clean package
  package:
    - FROM common-jre8
    - WORKDIR /app
    - ADD target/web.jar /app/web.jar
    - CMD ["java", "-jar", "web.jar"]
```

可以为不同环境指定不同的 构建脚本 和 打包脚本, 如果没有指定, 则是用默认脚本, 如下所示

```yaml
default:
  package:
    - FROM common-jre8
    - WORKDIR /app
    - ADD target/web.jar /app/web.jar
    - CMD ["java", "-jar", "web.jar"]

dev:
  build:
    - mvn -Pdev -DskipTests=true clean package
test:
  build:
    - mvn -Ptest -DskipTests=true clean package
```

该文件为 `dev` 和 `test` 指定了不同的 构建脚本, 但是他们共享同一个 打包脚本

## 许可证

Guo Y.K., MIT License
