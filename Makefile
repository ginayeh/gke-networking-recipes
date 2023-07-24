# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

BOSKOS_RESOURCE_TYPE ?= gke-internal-project
RUN_IN_PROW ?= false
LOCATION ?= us-central1-c
NUM_NODES ?= 3
TEST_TO_RUN ?= .*

all: bin/recipes-test

bin/recipes-test:
	mkdir bin/
	go test -c -o $@ ./test

.PHONY: test
test: bin/recipes-test
	bin/recipes-test \
		--run-in-prow=$(RUN_IN_PROW) \
		--boskos-resource-type=$(BOSKOS_RESOURCE_TYPE) \
		--test-project-id=$(PROJECT_ID) \
		--cluster-name=$(CLUSTER_NAME) \
		--location=$(LOCATION) \
		--num-nodes=$(NUM_NODES) \
		-test.run=$(TEST_TO_RUN) \

.PHONY: clean
clean:
	rm -rf bin/