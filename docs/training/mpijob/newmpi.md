# arean 支持 kubeflow1.2.0


## 背景

kubeflow1.2.0的版本默认mpi-operato的版本是v1alpha2。而arena只支持到v1alpha1。
由于这两个版本的mpijob的CRD变化很大，所以arena需要做一些改动来支持新的mpi-operator


## 编译及使用

### 编译
```
export GOPATH=/root/arenabuild
mkdir -p ${GOPATH}/src/github.com/kubeflow
cd  ${GOPATH}/src/github.com/kubeflow
git clone  https://github.com/qqsun8819/arena.git
cd arena
make
``` 


### 安装

在已经安装了arena的机器上
```
cp bin/arena /usr/local/bin/arena
cp charts/mpijob/templates/mpijob.yaml ~/charts/mpijob/templates/mpijob.yaml
```

如果没有安装arena，先[安装arena](https://arena-docs.readthedocs.io/en/latest/installation/binary/)，
再执行上述步骤

### 使用

提交
```
arena submit mpi \
 --name=mpi-dist2 \
 --gpus=1 \
 --workers=2 \
 --image=172.19.1.16:5000/qqsun8819/mpi-operator:tensorflow-benchmarks2  \
 --env=GIT_SYNC_BRANCH=master \
 --sync-mode=git \
 --sync-source=https://code.aliyun.com/qqsun8819/benchmarks.git \
 "python code/benchmarks/scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py --model resnet101 --batch_size 64 --variable_update horovod --train_dir=/training_logs --summary_verbosity=3 --save_summaries_steps=10"
```

查看
```
arena get mpi-dist2 --type mpijob
```

删除
```
arena delete mpi-dist2
```


### 高级用法

提交读取minio上的样本的任务
```

arena submit mpi \
  --name=mpi-dist3 \
  --gpus=1 \
  --workers=2 \
  --image=172.19.1.16:5000/qqsun8819/mpi-operator:tensorflow-benchmarks2  \
  --env=GIT_SYNC_BRANCH=master,S3_USE_HTTPS=0,S3_VERIFY_SSL=0,S3_ENDPOINT=http://10.43.4.158:9000,AWS_ACCESS_KEY_ID=49a1af5f-154b-426e-9059-cff696b408c4,AWS_SECRET_ACCESS_KEY=115aa465-a102-40bf-80bc-9b1026ccfce0 \
  --sync-mode=git \
  --sync-source=https://code.aliyun.com/qqsun8819/benchmarks.git \
  "python code/benchmarks/scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py --model resnet101 --batch_size 64 --variable_update horovod --train_dir=/training_logs --summary_verbosity=3 --save_summaries_steps=10

```