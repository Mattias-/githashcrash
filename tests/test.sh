#!/bin/bash
set -ex

pattern='^00.*'

# Correctness
cp -r /repo /tmp/repo0 && cd $_
recreate_cmd=$(githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd" &> /dev/null
hash_0=$(git rev-parse HEAD)
[[ $hash_0 =~ $pattern ]] || exit 1

# Correctness with seed
cp -r /repo /tmp/repo1 && cd $_
recreate_cmd=$(GITHASHCRASH_SEED=11 githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd" &> /dev/null
hash_1=$(git rev-parse HEAD)
[[ $hash_1 =~ $pattern ]] || exit 1

# Consistency
cp -r /repo /tmp/repo2 && cd $_
recreate_cmd2=$(GITHASHCRASH_SEED=11 githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd2" &> /dev/null
hash_2=$(git rev-parse HEAD)
diff <(echo $hash_1) <(echo $hash_2)

# Recreate in another copy of repo
cp -r /repo /tmp/repo3 && cd $_
bash -c "$recreate_cmd" &> /dev/null
hash_3=$(git rev-parse HEAD)
diff <(echo $hash_1) <(echo $hash_3)

# Different prefix yield different hash
cp -r /repo /tmp/repo4 && cd $_
recreate_cmd4=$(GITHASHCRASH_SEED=22 githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd4" &> /dev/null
hash_4=$(git rev-parse HEAD)
[[ $hash_4 =~ $pattern ]] || exit 1
if diff <(echo $hash_1) <(echo $hash_4); then
    echo "Expected different hashes"
    exit 1
fi

# Pass object as argument
cd /repo
obj=$(git cat-file -p HEAD)
recreate_cmd=$(GITHASHCRASH_SEED=33 githashcrash "$pattern" "$obj" | tail -1)
cp -r /repo /tmp/repo5
cd /tmp/repo5
bash -c "$recreate_cmd" &> /dev/null
hash_5=$(git rev-parse HEAD)
[[ $hash_5 =~ $pattern ]] || exit 1

echo "Tests passed!"
