#!/bin/bash -x

set -euo pipefail

rm -vf ~/.config/bigstruct/storage.db

alias bs=bigstruct

bigstruct pak import bigpack/testdata/bigpack-scylla

bigstruct index set --name base --property 4 \
	-v base=4

bigstruct index set --name service --property cql \
	-v base=4 -v service=cql

bigstruct index set --name instance --property i3.xlarge \
	-v base=4 -v service=cql -v instance=i3.xlarge

bigstruct index set --name account --property 2 \
	-v base=4 -v service=cql -v instance=i3.xlarge -v account=2

bigstruct index set --name cluster --property 3 \
	-v base=4 -v service=cql -v instance=i3.xlarge -v account=2 -v cluster=3

bigstruct index set --name dc --property 4 \
	-v base=4 -v service=cql -v instance=i3.xlarge -v account=2 -v cluster=3 -v dc=4

bigstruct index set --name server --property 5 \
	-v base=4 -v service=cql -v instance=i3.xlarge -v account=2 -v cluster=3 -v dc=4 -v server=5

bigstruct set --index cluster=3 \
	--value /etc/scylla/scylla.yaml/batch_size_fail_threshold_in_kb=250 \
	--value /etc/scylla/scylla.yaml/batch_size_warn_threshold_in_kb

bigstruct set --index dc=4 \
	--value /etc/scylla/scylla.yaml/commitlog_segment_size_in_mb=128 \
	--value /etc/scylla/scylla.yaml/commitlog_sync_period_in_ms=500  \
	--value /etc/scylla/scylla.yaml/commitlog_total_space_in_mb=1024 \
	--value /etc/scylla/scylla.yaml/commitlog_sync=fixed

bigstruct set --index account=2 \
	--value /etc/scylla/scylla.yaml/rpc_port=19160 \
	--value /etc/scylla/scylla.yaml/api_port=8080

bigstruct set --index server=5 \
	--value /etc/scylla/scylla.yaml/rpc_port=29160 \
	--value /etc/scylla/scylla.yaml/api_port=9090 \
	--value /etc/scylla/scylla.yaml/endpoint_snitch=ComplexSnitch

bigstruct get --index server=5 /etc/scylla/scylla.yaml

bigstruct schema list --namespace version=4.2.1
