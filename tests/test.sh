#!/bin/bash
set -ex

githashcrash=$1
repo=$2
testsdir=$(mktemp -d)
pattern='^00.*'

cleanup() {
    rm -rf "$testsdir"
}

trap cleanup EXIT

# Correctness
cp -r "$repo" "$testsdir/repo0" && cd "$_"
recreate_cmd=$($githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd" &>/dev/null
hash_0=$(git rev-parse HEAD)
[[ $hash_0 =~ $pattern ]] || exit 1

# Correctness with seed
cp -r "$repo" "$testsdir/repo1" && cd "$_"
recreate_cmd=$(GITHASHCRASH_SEED=11 GITHASHCRASH_THREADS=1 $githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd" &>/dev/null
hash_1=$(git rev-parse HEAD)
[[ $hash_1 =~ $pattern ]] || exit 1

# Consistency
cp -r "$repo" "$testsdir/repo2" && cd "$_"
recreate_cmd2=$(GITHASHCRASH_SEED=11 GITHASHCRASH_THREADS=1 $githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd2" &>/dev/null
hash_2=$(git rev-parse HEAD)
diff <(echo "$hash_1") <(echo "$hash_2")

# Recreate in another copy of repo
cp -r "$repo" "$testsdir/repo3" && cd "$_"
bash -c "$recreate_cmd" &>/dev/null
hash_3=$(git rev-parse HEAD)
diff <(echo "$hash_1") <(echo "$hash_3")

# Different seed yield different hash
cp -r "$repo" "$testsdir/repo4" && cd "$_"
recreate_cmd4=$(GITHASHCRASH_SEED=22 GITHASHCRASH_THREADS=1 $githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd4" &>/dev/null
hash_4=$(git rev-parse HEAD)
[[ $hash_4 =~ $pattern ]] || exit 1
if diff <(echo "$hash_1") <(echo "$hash_4"); then
    echo "Expected different hashes"
    exit 1
fi

# Pass object as argument
cp -r "$repo" "$testsdir/repo5.1" && cd "$_"
obj=$(git cat-file -p HEAD)
recreate_cmd=$(GITHASHCRASH_SEED=33 $githashcrash "$pattern" "$obj" | tail -1)
cp -r "$repo" "$testsdir/repo5.2" && cd "$_"
bash -c "$recreate_cmd" &>/dev/null
hash_5=$(git rev-parse HEAD)
[[ $hash_5 =~ $pattern ]] || exit 1

# Replacement
repl_1="REPLACEME"
cp -r "$repo" "$testsdir/repo6" && cd "$_"
git commit --allow-empty -m "Test commit" -m "Message message" -m "hello world $repl_1 hello world" -m "More message"
git show | grep "$repl_1"
recreate_cmd=$($githashcrash "$pattern" | tail -1)
bash -c "$recreate_cmd" &>/dev/null
hash_6=$(git rev-parse HEAD)
[[ $hash_6 =~ $pattern ]] || exit 1
! git show | grep "$repl_1"


echo "Tests passed!"
