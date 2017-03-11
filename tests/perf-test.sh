#!/bin/bash
set -ex

cp -r /repo /tmp/repo
cd /tmp/repo
recreate_cmd=$(githashcrash '^000000.*' | tail -1)
bash -c "$recreate_cmd"
[[ $(git rev-parse HEAD) =~ ^000000.* ]] || exit 1
git show -s
git cat-file -p HEAD
