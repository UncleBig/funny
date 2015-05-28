##插件说明

- 插件执行需要参数为告警阀值（w参数），紧急参数（c参数），post地址（p参数），执行样例：./CheckDisk  -w 60 -c 90 -p https://127.0.0.1:8088/api/platform

- 告警参数选取
    - check_cpu       cpu使用率
    - check_mem    内存使用率
    - check_load     1,5,15的负载
    - check_disk        磁盘使用率
    - check_net         每秒种接收发送数据的比特值