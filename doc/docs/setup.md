* Installer golang

```shell
sudo apt install golang
```

* Install and configure mysql
```shell
sudo apt install mysql-server
```

Configure DB:
```shell
sudo mysql
create database xxxxx
use xxxxx
```

Executing an sql source (.sql file with DROP CREATE instructions)
``` sql
source /path/to/sqlFile
```

If driver is needed
```shell
go get -u github.com/go-sql-driver/mysql
```

* Build and display this documentation locally

Install Python
```shell
sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt update
sudo apt install python3.13
```

Create a venv (do it in the doc folder)
```shell
cd doc
python3.13 -m venv venv
source venv/bin/activate
pip install mkdocs-material
mkdocs serve
```