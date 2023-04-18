#!/bin/bash

go install golang.org/x/tools/cmd/godoc@latest

git config --global user.name "Frank Jahn"
git config --global user.email "fjahn78@gmail.com"
git config --global alias.ci commit
git config --global alias.cm "commit -m"
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.bd "branch --delete"
git config --global alias.bdf "branch --delete --force"
git config --global alias.unstage "reset HEAD --"
git config --global alias.p push
git config --global alias.pp pull
