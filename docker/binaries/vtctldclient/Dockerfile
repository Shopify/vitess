# Copyright 2024 The Vitess Authors.
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

ARG VT_BASE_VER=latest
ARG DEBIAN_VER=stable-slim

FROM vitess/lite:${VT_BASE_VER} AS lite

FROM debian:${DEBIAN_VER}

RUN apt-get update && \
   apt-get upgrade -qq && \
   apt-get install jq curl -qq --no-install-recommends && \
   apt-get autoremove && \
   apt-get clean && \
   rm -rf /var/lib/apt/lists/*

COPY --from=lite /vt/bin/vtctldclient /usr/bin/

# add vitess user/group and add permissions
RUN groupadd -r --gid 2000 vitess && \
   useradd -r -g vitess --uid 1000 vitess

CMD ["/usr/bin/vtctldclient"]
