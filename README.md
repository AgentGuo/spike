# FaaS

## 1. How to get started

step1: Start MySQL, for example, start MySQL with Docker, and then create the `faas` database

```shell
# Download the mysql:8.0.31 image
docker pull mysql:8.0.31

# Start the MySQL instance
docker run --name faas-mysql -e MYSQL_ROOT_PASSWORD=faaspassword -p 3306:3306 -d mysql:8.0.31

# After the MySQL instance starts, enter the MySQL container
docker exec -it mysql-faas mysql -uroot -pfaaspassword -e "CREATE DATABASE faas;"
```

step2: build faas

```shell
go build -o build/faas-server github.com/AgentGuo/faas/cmd/server
```

step3: run faas

```shell
cd build
./faas-server -f faas.yaml
```