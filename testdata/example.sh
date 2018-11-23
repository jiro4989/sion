#!/bin/bash

set -eu

sion help
sion config init

# 連続してコマンドを投げたい場合
sion session open &

sion command localhost ls
sion command localhost git clone https://github.com/jiro4989/arth --before 'which arth'
sion command server1 cp hoge.conf /tmp/hoge.conf -o root -g root
sion command server1 mkdir /var/log/test -o root -g develop
sion command server1 groupadd developers
sion command server1 adduser developer1 -g developers
sion command server1 replace hoge.conf 'hoge' '#hoge'
sion command server1 rm hoge.conf
sion command server1 yum httpd
sion command web yum jboss
sion command web yum netutil
sion command db yum mysqlclients
sion command db yum go
sion command db yum oracle

sion session run

# 都度SSHしてコマンドを投げる場合はsession openしない
sion command db yum oracle
