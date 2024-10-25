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

## Api

List APi

<pre>
Authorization Basic Auth 
username = admin
password = admin

Category
http://localhost:8080/categories        method:get
http://localhost:8080/categories/{id}   method:get
http://localhost:8080/categories        method:post     body:[ name]
http://localhost:8080/categories/{id}   method:put      body:[ name]
http://localhost:8080/categories/{id}   method:delete

item
http://localhost:8080/items             method:get      param:[search, sort, order, page, limit]
http://localhost:8080/items/{id}        method:get      
http://localhost:8080/items             method:post     body:[category_id, price, name, description]
http://localhost:8080/items/{id}        method:put      body:[category_id, price, name, description]
http://localhost:8080/items/{id}        method:delete
</pre>

## arsitektur

Modular

<pre>
/app
└── app_core
    ├── dto
    ├── handler
    ├── middleware
    ├── migration
    ├── model
    ├── route
    └── utils
└── item_module
    ├── dto
    ├── handler
    ├── model
    ├── route
</pre>

## Module

app_core adalah module utama yang bersi core prject rules, helper, dan startup project. kegunaannya adalah agar saat membuat project / aplikasi baru maka tidak perlu membuat ulang dari awal
item_module adalah sample module yang akan dikembangkan di aplikasi

## Component

DTO (Data Transfer Object)
Fungsi: DTO berfungsi sebagai validator dan untuk mentransfer data dari request ke controller. DTO biasanya berisi struktur data yang diharapkan dari permintaan, dan sering kali digunakan untuk memvalidasi data yang masuk. File DTO juga dapat digunakan untuk dokumentasi otomatis dengan Swagger.

Handler / Controller
Fungsi: Controller bertanggung jawab untuk memproses data yang diterima dari request. Controller mengatur logika bisnis aplikasi, termasuk interaksi dengan model, pengolahan data, dan menentukan respons yang dikirim kembali ke klien.

Middleware
Fungsi: Middleware berfungsi untuk menangani berbagai fungsi perantara sebelum mencapai controller. Ini termasuk otentikasi, logging, dan pengaturan CORS, yang membantu dalam pengelolaan alur permintaan.

Migration
Fungsi: Migration digunakan untuk mengelola perubahan skema basis data. File migration berisi skrip SQL yang digunakan untuk membuat atau memperbarui tabel dan struktur database. Ini membantu dalam mengatur versi database seiring perkembangan aplikasi.

Model
Fungsi: Model adalah representasi dari data yang disimpan dalam basis data. Model mendefinisikan struktur data dan menyediakan metode untuk berinteraksi dengan database, seperti melakukan operasi CRUD (Create, Read, Update, Delete).

Route
Fungsi: Route mendefinisikan endpoint HTTP yang tersedia dalam aplikasi. Ini menghubungkan URL yang diminta oleh klien dengan controller yang sesuai, sehingga ketika sebuah request diterima, rute akan menentukan handler mana yang harus dipanggil.

Utils
Fungsi: Utils berisi fungsi-fungsi utilitas yang dapat digunakan di seluruh aplikasi. Ini termasuk fungsi untuk koneksi database, format tanggal, dan fungsi umum lainnya yang tidak terkait langsung dengan logika bisnis tetapi berguna dalam berbagai konteks.
