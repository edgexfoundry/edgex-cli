//
// Copyright (c) 2020 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

@Library("edgex-global-pipelines@82bb14f08269580ddcce05c9216dcef8a02c9b4b") _

edgeXBuildGoApp (
    project: 'edgex-cli',
    arch: ['amd64'],
    testScript: 'make test',
    buildScript: 'make build-all',
    artifactTypes: ['archive'],
    artifactRoot: './archives/bin'
) 