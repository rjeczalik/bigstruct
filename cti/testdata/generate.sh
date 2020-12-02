#!/bin/bash -x

set -euo pipefail

docker_cat() {
	mkdir -p "docker/${1}/tag/${2}/$(dirname ${3})"
	docker run --rm -ti --entrypoint /bin/sh "${1}:${2}" -c "cat ${3}" > "docker/${1}/tag/${2}/${3}"
}

docker_cat_scylla() {
	docker_cat scylladb/scylla "${1}" /etc/scylla/scylla.yaml
	docker_cat scylladb/scylla "${1}" /etc/scylla.d/cpuset.conf
	docker_cat scylladb/scylla "${1}" /etc/scylla.d/io.conf
	docker_cat scylladb/scylla "${1}" /etc/sysconfig/scylla-server
}

main() {
	docker_cat_scylla 4.0.11
	docker_cat_scylla 4.1.9
	docker_cat_scylla 4.2.1

	find docker -type f -name io.conf -exec sed -i -s 's/# SEASTAR_IO/SEASTAR_IO/g' {} \;
	find docker -type f -name cpuset.conf -exec sed -i -s 's/# CPUSET/CPUSET/g' {} \;
}

main 1>&2
