# AliCloud-tools

> 本工具为辅助使用阿里云API操作ECS以及策略组的小工具,如用于非法用途,与本人无关,出现任何问题,后果自负.

# 使用方式

## 查看所有地域信息
```
./AliCloud-Tools -a <AccessKey> -s <SecretKey> --regions

```
## 查看所有实例信息
```
./AliCloud-Tools -a <AccessKey> -s <SecretKey> ecs --list

```


## 查看所有正在运行的实例信息
```
./AliCloud-Tools -a <AccessKey> -s <SecretKey> ecs --list --runner

```

## 查看指定实例的信息
```
./AliCloud-Tools -a <AccessKey> -s <SecretKey> ecs --eid <InstanceId>

```

## 执行命令
```
./AliCloud-Tools -a <AccessKey> -s <SecretKey> ecs exec -I <InstanceId[,InstanceId,InstanceId,...]> -c "touch /tmp/123123aaaa.txt"
```

## 查看安全组策略
```
./AliCloud-Tools -a <AccessKey> -s <SecretKey> sg --eid <SecruityGroupId>

```

## 增加安全组策略
```
./AliCloud-Tools -a <AccessKey> -s <SecretKey> --eid <SecruityGroupId> --action add --protocol tcp --port 0/0 --ip 0.0.0.0/0

```

## 删除安全组策略
```
./AliCloud-Tools -a <AccessKey> -s <SecretKey> --eid <SecruityGroupId> --action del --protocol tcp --port 0/0 --ip 0.0.0.0/0

```
