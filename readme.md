# Aplication test with golang

## Prerequisites

- [Docker](https://www.docker.com/) installed on your machine.
- [Docker Compose](https://docs.docker.com/compose/) installed.

## Getting Started

clone repository

```text
$ git clone https://github.com/yourusername/your-repo.git
```

running docker container

```text
$ docker-compose up
```

waiting container go running

<pre>
belajar-app-1  | wait-for-it.sh: timeout occurred after waiting 15 seconds for mysql:3306
mysql          | 2024-10-25T03:16:56.994005Z 1 [System] [MY-013577] [InnoDB] InnoDB initialization has ended.
mysql          | 2024-10-25T03:16:57.514258Z 0 [Warning] [MY-010068] [Server] CA certificate ca.pem is self signed.
mysql          | 2024-10-25T03:16:57.514872Z 0 [System] [MY-013602] [Server] Channel mysql_main configured to support TLS. Encrypted connections are now supported for this channel.
mysql          | 2024-10-25T03:16:57.535017Z 0 [Warning] [MY-011810] [Server] Insecure configuration for --pid-file: Location '/var/run/mysqld' in the path is accessible to all OS users. Consider choosing a different directory.
mysql          | 2024-10-25T03:16:57.563221Z 0 [System] [MY-011323] [Server] X Plugin ready for connections. Bind-address: '::' port: 33060, socket: /var/run/mysqld/mysqlx.sock
mysql          | 2024-10-25T03:16:57.563366Z 0 [System] [MY-010931] [Server] /usr/sbin/mysqld: ready for connections. Version: '9.1.0'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server - GPL.
belajar-app-1  | 2024/10/25 03:17:18 Starting the server...
belajar-app-1  | 2024/10/25 03:17:18 Migration executed successfully
belajar-app-1  | 2024/10/25 03:17:18 Migration executed successfully
belajar-app-1  | 2024/10/25 03:17:18 Migration executed successfully
belajar-app-1  | 2024/10/25 03:17:18 Migration executed successfully
belajar-app-1  | 2024/10/25 03:17:18 Server is running on :8080
</pre>
