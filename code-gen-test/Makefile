domain=oam.fp.net
group=steward
version=v1
apiDir=pkg/apis/${group}/${version}
init:
	mkdir -p ${apiDir}
	mkdir -p pkg/gen/${group}/${version}
	# echo "// +k8s:deepcopy-gen=package\n// +groupName=${group}.${domain}\npackage ${version}" > ${apiDir}/doc.go
	# echo "package v1" > ${apiDir}/types.go
	# echo "package v1\nconst (\nGroupName = \"${group}.${domain}\")" > pkg/apis/${group}/register.go


gen:
	./hack/gen.sh