#!/bin/bash

set -eu

sion help
sion config init

# 対象サーバとコネクションを貼る(サーバ)
# 貼ってる間はロックファイルを生成する
sion connection open web &

# 対象サーバにコマンドをなげる(クライアント)
# クライアント・サーバ方式でコネクションを張りっぱなしにしてコマンドを投げ続ける
sion command ls
sion command git clone https://github.com/jiro4989/arth --before 'which arth'
sion command cp hoge.conf /tmp/hoge.conf -o root -g root
sion command mkdir /var/log/test -o root -g develop
sion command groupadd developers
sion command adduser developer1 -g developers
sion command replace hoge.conf 'hoge' '#hoge'
sion command rm hoge.conf
sion command yum httpd

# サーバを停止、ロックファイルを削除
sion connection close web

# いちいちコネクションを貼るのがめんどくさくて、都度コネクションを貼るので十分な
# 場合はこっち
sion command standalone db yum oracle
