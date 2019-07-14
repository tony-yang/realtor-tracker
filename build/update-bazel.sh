set -o errexit
set -o nounset
set -o pipefail

if [[ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]]; then
  echo "Updating protos..." >&2
elif ! command -v bazel &>/dev/null; then
  echo "Install bazel >&2"
  exit 1
elif ! command -v kazel &>/dev/null; then
  echo "Install kazel >&2"
  exit 1
else
  (
    set -o xtrace
    bazel run //build:update-bazel
  )
  exit 0
fi

cd "${BUILD_WORKSPACE_DIRECTORY}"

bazel run :gazelle -- fix -external=vendored
kazel --cfg-path=./build/.kazelcfg.json
