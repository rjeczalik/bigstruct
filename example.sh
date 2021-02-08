#!/bin/bash -x

set -euo pipefail

rm -vf ~/.config/bigstruct/storage.db

alias bs=bigstruct

bigstruct namespace set --name global  --property false --priority 0
bigstruct namespace set --name version --property true  --priority 10
bigstruct namespace set --name account --property true  --priority 20
bigstruct namespace set --name cluster --property true  --priority 30
bigstruct namespace set --name dc      --property true  --priority 40
bigstruct namespace set --name server  --property true  --priority 50

bigstruct index set --name global --property "" \
	-v global \
	-x global

bigstruct index set --name version --property 4.2.1 \
	-v version=4.2.1 \
	-x version=4.2.1

bigstruct index set --name account --property 2 \
	-v version=4.2.1 -v account=2 \
	-x version=4.2.1

bigstruct index set --name cluster --property 3 \
	-v version=4.2.1 -v account=2 -v cluster=3 \
	-x version=4.2.1

bigstruct index set --name dc --property 4 \
	-v version=4.2.1 -v account=2 -v cluster=3 -v dc=4 \
	-x version=4.2.1

bigstruct index set --name server --property 5 \
	-v version=4.2.1 -v account=2 -v cluster=3 -v dc=4 -v server=5 \
	-x version=4.2.1

bigstruct set --index version=4.2.1 --import isr/testdata/docker/scylladb/scylla/tag/4.2.1/ --schema

bigstruct set --index version=4.2.1 --import isr/testdata/docker/scylladb/scylla/tag/4.2.1/

bigstruct set --index cluster=3 \
	--value /etc/scylla/scylla.yaml/batch_size_fail_threshold_in_kb=250 \
	--value /etc/scylla/scylla.yaml/batch_size_warn_threshold_in_kb=25  \
	--value /etc/scylla/scylla.yaml/api_ui_dir

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

bigstruct set --index version=4.2.1 \
	--value /etc/scylla/scylla.yaml/listen_address=127.0.0.1 \
	--value /etc/scylla/scylla.yaml/rpc_address=127.0.0.1    \
	--value /etc/scylla/scylla.yaml/api_doc_dir

bigstruct set --index version=4.2.1 \
	--value /etc/scylla/scylla.yaml/listen_address=localhost

bigstruct get --index server=5 /etc/scylla/scylla.yaml

bigstruct schema list --namespace version=4.2.1
