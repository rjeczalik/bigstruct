#!/bin/bash -x

set -euo pipefail

rm -vf ~/.config/bigstruct/storage.db

alias cti=bigstruct
alias bs=bigstruct

bigstruct namespace set --name global  --property false --priority 0
bigstruct namespace set --name version --property true  --priority 10
bigstruct namespace set --name account --property true  --priority 20
bigstruct namespace set --name cluster --property true  --priority 30
bigstruct namespace set --name dc      --property true  --priority 40
bigstruct namespace set --name server  --property true  --priority 50

bigstruct value set --namespace global --import isr/testdata/docker/scylladb/scylla/tag/4.2.1/ --schema

bigstruct value set --namespace cluster/3 \
	--value /etc/scylla/scylla.yaml/batch_size_fail_threshold_in_kb=250 \
	--value /etc/scylla/scylla.yaml/batch_size_warn_threshold_in_kb=25  \
	--value /etc/scylla/scylla.yaml/api_ui_dir

bigstruct value set --namespace dc/4 \
	--value /etc/scylla/scylla.yaml/commitlog_segment_size_in_mb=128 \
	--value /etc/scylla/scylla.yaml/commitlog_sync_period_in_ms=500  \
	--value /etc/scylla/scylla.yaml/commitlog_total_space_in_mb=1024 \
	--value /etc/scylla/scylla.yaml/commitlog_sync=fixed

bigstruct value set --namespace account/2 \
	--value /etc/scylla/scylla.yaml/rpc_port=19160 \
	--value /etc/scylla/scylla.yaml/api_port=8080

bigstruct value set --namespace server/5 \
	--value /etc/scylla/scylla.yaml/rpc_port=29160 \
	--value /etc/scylla/scylla.yaml/api_port=9090 \
	--value /etc/scylla/scylla.yaml/endpoint_snitch=ComplexSnitch

bigstruct value set --namespace version/1 \
	--value /etc/scylla/scylla.yaml/listen_address=127.0.0.1 \
	--value /etc/scylla/scylla.yaml/rpc_address=127.0.0.1    \
	--value /etc/scylla/scylla.yaml/api_doc_dir

bigstruct value set --namespace version/1 \
	--value /etc/scylla/scylla.yaml/listen_address=localhost

bigstruct index set --name cluster --property 3 \
	-v version=1 -v account=2 -v cluster=3 -v dc=4 -v server=5 \
	-x global

bigstruct get --index cluster/3 /etc/scylla/scylla.yaml

bigstruct schema list --namespace global
