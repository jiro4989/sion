#!/bin/bash

set -eu

sion localhost git clone https://github.com/jiro4989/arth --before 'which arth'
sion server1 cp hoge.conf /tmp/hoge.conf -o root -g root
sion server1 mkdir /var/log/test -o root -g develop
sion server1 groupadd developers
sion server1 adduser developer1 -g developers
sion server1 replace hoge.conf 'hoge' '#hoge'
sion server1 rm hoge.conf
sion server1 yum httpd
sion web yum jboss
sion web yum netutil
sion db yum mysqlclients
sion db yum go
sion db yum oracle
