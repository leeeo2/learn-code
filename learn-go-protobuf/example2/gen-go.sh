
ARGS="\
--proto_path=. \
--go_out=. \
--go_opt=paths=source_relative "

for dir in $(find . -name '*.proto' | xargs -I{} dirname {} | sort | uniq); do
    echo "building $dir/*.proto"
    protoc $ARGS $dir/*.proto
done