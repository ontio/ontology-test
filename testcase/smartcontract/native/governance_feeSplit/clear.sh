kill -9 $(cat pid)
rm -rf ./Chain/ ./Log/
echo 'passwordtest' | /app/gopath/src/github.com/ontio/ontology/ontology --testmode --gasprice=0 &
echo $! > pid