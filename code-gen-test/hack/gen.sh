RootDir="github.com/smarkm/k8s-crd/code-gen-test"
	CODEGEN="$(find /go/pkg/mod/k8s.io/ |grep generate-groups.sh| tail -1)"
	echo ">>> using codegen: ${CODEGEN}"
	# ensure we can execute the codegen script
	"${CODEGEN}" all \
    "${RootDir}/pkg/gen/steward" \
    "${RootDir}/pkg/apis" steward:v1 \
	--go-header-file "$(pwd)"/hack/boilerplate.go.txt