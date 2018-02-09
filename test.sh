echo "mode: $COVERMODE" > $COVERAGE_FILE

for PKG in $(go list ./...); do
  go test -v -coverprofile=$TMP_FILE -covermode=$COVERMODE $PKG
  if [ -f $TMP_FILE ]; then
    cat $TMP_FILE | tail -n +2 >> $COVERAGE_FILE
    rm $TMP_FILE
  fi
done
