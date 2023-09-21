# 后端

## docker compose

启动： docker compose up

停止： docker compose down

如果报  xx  it is not a shared mount.  

可以采用 `docker-compose-v1 up` 

原因： [Docker Compose V2 with Docker Desktop on Windows / WSL2 gives "/mnt/c not a shared mount" error · Issue #8558 ](https://github.com/docker/compose/issues/8558)

# 前端

## 安装 npm 和 nodejs
```shell
sudo apt update && sudo apt install npm
sudo npm install -g n
sudo n stable
```

## 快速启动

进入fe 目录底下，执行 

```shell
npm install && npm run dev
```

就可以在浏览器里面打开 localhost:3000 来查看