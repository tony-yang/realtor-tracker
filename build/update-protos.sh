set -o errexit
set -o nounset
set -o pipefail

if [[ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]]; then
  echo "Updating protos..." >&2
elif ! command -v bazel &>/dev/null; then
  echo "Install bazel >&2"
  exit 1
else
  (
    set -o xtrace
    bazel run //build:update-protos
  )
  exit 0
fi

protoc=$1
plugin=$2
grpc=$3
dest=$BUILD_WORKSPACE_DIRECTORY

genproto() {
  dir=$(dirname "$1")
  base=$(basename "$1")
  out=$dest/$dir/${base%.proto}.pb.go
  rm -f "$out"
  echo -e "\nprotoc=$protoc dir=$dir base=$base out=$out dest=$dest\n"
  "$protoc" "--plugin=$plugin" "--proto_path=${dir}" "--proto_path=${dest}" "--go_out=$grpc:$dest/$dir" "$1"
}

echo -n "Generating protos: " >&2
echo
for p in $(find . -name '*.proto' -not -path './vendor/*'); do
  echo -n "$p "
  echo
  genproto "$p"
done
echo -n "Update proto done" >&2
echo
