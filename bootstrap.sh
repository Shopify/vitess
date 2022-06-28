#!/bin/bash
# shellcheck disable=SC2164

# Copyright 2019 The Vitess Authors.
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

### This file is executed by 'make tools'. You do not need to execute it directly.

source ./dev.env

# Outline of this file.
# 0. Initialization and helper methods.
# 1. Installation of dependencies.

BUILD_JAVA=${BUILD_JAVA:-1}
BUILD_CONSUL=${BUILD_CONSUL:-1}
BUILD_CHROME=${BUILD_CHROME:-1}

VITESS_RESOURCES_DOWNLOAD_BASE_URL="https://github.com/vitessio/vitess-resources/releases/download"
VITESS_RESOURCES_RELEASE="v2.0"
VITESS_RESOURCES_DOWNLOAD_URL="${VITESS_RESOURCES_DOWNLOAD_BASE_URL}/${VITESS_RESOURCES_RELEASE}"
#
# 0. Initialization and helper methods.
#

[[ "$(dirname "$0")" = "." ]] || fail "bootstrap.sh must be run from its current directory"

# install_dep is a helper function to generalize the download and installation of dependencies.
#
# If the installation is successful, it puts the installed version string into
# the $dist/.installed_version file. If the version has not changed, bootstrap
# will skip future installations.
install_dep() {
  if [[ $# != 4 ]]; then
    fail "install_dep function requires exactly 4 parameters (and not $#). Parameters: $*"
  fi
  local name="$1"
  local version="$2"
  local dist="$3"
  local install_func="$4"

  version_file="$dist/.installed_version"
  if [[ -f "$version_file" && "$(cat "$version_file")" == "$version" ]]; then
    echo "skipping $name install. remove $dist to force re-install."
    return
  fi

  echo "<<< Installing $name $version >>>"

  # shellcheck disable=SC2064
  trap "fail '$name build failed'; exit 1" ERR

  # Cleanup any existing data and re-create the directory.
  rm -rf "$dist"
  mkdir -p "$dist"

  # Change $CWD to $dist before calling "install_func".
  pushd "$dist" >/dev/null
  # -E (same as "set -o errtrace") makes sure that "install_func" inherits the
  # trap. If here's an error, the trap will be called which will exit this
  # script.
  set -E
  $install_func "$version" "$dist"
  set +E
  popd >/dev/null

  trap - ERR

  echo "$version" > "$version_file"
}


#
# 1. Installation of dependencies.
#

# We should not use the arch command, since it is not reliably
# available on macOS or some linuxes:
# https://www.gnu.org/software/coreutils/manual/html_node/arch-invocation.html
get_arch() {
  uname -m
}

# Install protoc.
install_protoc() {
  local version="$1"
  local dist="$2"

  case $(uname) in
    Linux)  local platform=linux;;
    Darwin) local platform=osx;;
    *) echo "ERROR: unsupported platform for protoc"; exit 1;;
  esac

  case $(get_arch) in
      aarch64)  local target=aarch_64;;
      x86_64)  local target=x86_64;;
      arm64) case "$platform" in
          osx) local use_homebrew=1;;
          *) echo "ERROR: unsupported architecture for protoc"; exit 1;;
      esac;;
      *)   echo "ERROR: unsupported architecture for protoc"; exit 1;;
  esac

  # TODO (ajm188): remove this branch after protoc includes signed protoc binaries for M1 Macs.
  if [[ "$use_homebrew" -eq 1 ]]; then
    cat >&2 <<WARNING
WARN: Protobuf does not have a protoc binary for arm64 macos.
Checking for homebrew installation. Altertatively, you may install the x86-64
version if you have Rosetta installed; your mileage may vary.

See https://github.com/protocolbuffers/protobuf/issues/9397.
WARNING

    if [[ -z "$(command -v brew)" ]]; then
      echo "Could not find \`brew\` command. Please install homebrew and retry." >&2;
      exit 1;
    fi

    brew install protobuf;
    protobuf_base="$(brew list protobuf | grep -E 'LICENSE$' | sed 's:/LICENSE$::')"
    if [[ -z "$protobuf_base" ]]; then
      echo "Could not find \`protobuf\` directory after installing. Please verify the output of \`brew info protobuf\`" >&2;
      exit 1;
    fi

    ln -snf "${protobuf_base}/bin" "${dist}/bin"
    ln -snf "${protobuf_base}/include" "${dist}/include"
  else
    # This is how we'd download directly from source:
    # wget https://github.com/protocolbuffers/protobuf/releases/download/v$version/protoc-$version-$platform-${target}.zip
    $VTROOT/tools/wget-retry "${VITESS_RESOURCES_DOWNLOAD_URL}/protoc-$version-$platform-${target}.zip"
    unzip "protoc-$version-$platform-${target}.zip"
  fi

  ln -snf "$dist/bin/protoc" "$VTROOT/bin/protoc"
}


# Install Zookeeper.
install_zookeeper() {
  local version="$1"
  local dist="$2"
  zk="zookeeper-$version"
  vtzk="vt-zookeeper-$version"
  # This is how we'd download directly from source:
  # wget "https://dlcdn.apache.org/zookeeper/$zk/apache-$zk.tar.gz"
  $VTROOT/tools/wget-retry "${VITESS_RESOURCES_DOWNLOAD_URL}/apache-${zk}.tar.gz"
  tar -xzf "$dist/apache-$zk.tar.gz"
  mv $dist/apache-$zk $dist/$vtzk
  mvn -f $dist/$vtzk/zookeeper-contrib/zookeeper-contrib-fatjar/pom.xml clean install -P fatjar -DskipTests
  mkdir -p $dist/$vtzk/lib
  cp "$dist/$vtzk/zookeeper-contrib/zookeeper-contrib-fatjar/target/$zk-fatjar.jar" "$dist/$vtzk/lib/$zk-fatjar.jar"
  rm -rf "$zk.tar.gz"
}


# Download and install etcd, link etcd binary into our root.
install_etcd() {
  local version="$1"
  local dist="$2"

  case $(uname) in
    Linux)  local platform=linux; local ext=tar.gz;;
    Darwin) local platform=darwin; local ext=zip;;
    *)   echo "ERROR: unsupported platform for etcd"; exit 1;;
  esac

  case $(get_arch) in
      aarch64)  local target=arm64;;
      x86_64)  local target=amd64;;
      arm64)  local target=arm64;;
      *)   echo "ERROR: unsupported architecture for etcd"; exit 1;;
  esac

  file="etcd-${version}-${platform}-${target}.${ext}"

  # This is how we'd download directly from source:
  # download_url=https://github.com/etcd-io/etcd/releases/download
  # wget "$download_url/$version/$file"
  $VTROOT/tools/wget-retry "${VITESS_RESOURCES_DOWNLOAD_URL}/${file}"
  if [ "$ext" = "tar.gz" ]; then
    tar xzf "$file"
  else
    unzip "$file"
  fi
  rm "$file"
  ln -snf "$dist/etcd-${version}-${platform}-${target}/etcd" "$VTROOT/bin/etcd"
  ln -snf "$dist/etcd-${version}-${platform}-${target}/etcdctl" "$VTROOT/bin/etcdctl"
}


# Download and install k3s, link k3s binary into our root
install_k3s() {
  local version="$1"
  local dist="$2"
  case $(uname) in
    Linux)  local platform=linux;;
    *)   echo "WARNING: unsupported platform. K3s only supports running on Linux, the k8s topology will not be available for local examples."; return;;
  esac

  case $(get_arch) in
      aarch64)  local target="-arm64";;
      x86_64) local target="";;
      arm64)  local target="-arm64";;
      *)   echo "WARNING: unsupported architecture, the k8s topology will not be available for local examples."; return;;
  esac

  file="k3s${target}"

  local dest="$dist/k3s${target}-${version}-${platform}"
  # This is how we'd download directly from source:
  # download_url=https://github.com/rancher/k3s/releases/download
  # wget -O  $dest "$download_url/$version/$file"
  $VTROOT/tools/wget-retry -O $dest "${VITESS_RESOURCES_DOWNLOAD_URL}/$file-$version"
  chmod +x $dest
  ln -snf  $dest "$VTROOT/bin/k3s"
}


# Download and install consul, link consul binary into our root.
install_consul() {
  local version="$1"
  local dist="$2"

  case $(uname) in
    Linux)  local platform=linux;;
    Darwin) local platform=darwin;;
    *)   echo "ERROR: unsupported platform for consul"; exit 1;;
  esac

  case $(get_arch) in
    aarch64)  local target=arm64;;
    x86_64)  local target=amd64;;
    arm64)  local target=arm64;;
    *)   echo "ERROR: unsupported architecture for consul"; exit 1;;
  esac

  # This is how we'd download directly from source:
  # download_url=https://releases.hashicorp.com/consul
  # wget "${download_url}/${version}/consul_${version}_${platform}_${target}.zip"
  $VTROOT/tools/wget-retry "${VITESS_RESOURCES_DOWNLOAD_URL}/consul_${version}_${platform}_${target}.zip"
  unzip "consul_${version}_${platform}_${target}.zip"
  ln -snf "$dist/consul" "$VTROOT/bin/consul"
}


# Download chromedriver
install_chromedriver() {
  local version="$1"
  local dist="$2"

  case $(uname) in
    Linux)  local platform=linux;;
    *)   echo "Platform not supported for vtctl-web tests. Skipping chromedriver install."; return;;
  esac

  if [ "$(arch)" == "aarch64" ] ; then
      os=$(cat /etc/*release | grep "^ID=" | cut -d '=' -f 2)
      case $os in
        ubuntu|debian)
          sudo apt-get update -y && sudo apt install -y --no-install-recommends unzip libglib2.0-0 libnss3 libx11-6
        ;;
        centos|fedora)
          sudo yum update -y && yum install -y libX11 unzip wget
        ;;
      esac
      echo "For Arm64, using prebuilt binary from electron (https://github.com/electron/electron/) of version 76.0.3809.126"
      $VTROOT/tools/wget-retry https://github.com/electron/electron/releases/download/v6.0.3/chromedriver-v6.0.3-linux-arm64.zip
      unzip -o -q chromedriver-v6.0.3-linux-arm64.zip -d "$dist"
      rm chromedriver-v6.0.3-linux-arm64.zip
  else
      $VTROOT/tools/wget-retry "https://chromedriver.storage.googleapis.com/$version/chromedriver_linux64.zip"
      unzip -o -q chromedriver_linux64.zip -d "$dist"
      rm chromedriver_linux64.zip
  fi
}

install_all() {
  echo "##local system details..."
  echo "##platform: $(uname) target:$(get_arch) OS: $os"
  # protoc
  protoc_ver=3.19.4
  install_dep "protoc" "$protoc_ver" "$VTROOT/dist/vt-protoc-$protoc_ver" install_protoc

  # zk
  zk_ver=${ZK_VERSION:-3.8.0}
  if [ "$BUILD_JAVA" == 1 ] ; then
    install_dep "Zookeeper" "$zk_ver" "$VTROOT/dist" install_zookeeper
  fi

  # etcd
  command -v etcd && echo "etcd already installed" || install_dep "etcd" "v3.5.3" "$VTROOT/dist/etcd" install_etcd

  # k3s
  command -v  k3s || install_dep "k3s" "v1.0.0" "$VTROOT/dist/k3s" install_k3s

  # consul
  if [ "$BUILD_CONSUL" == 1 ] ; then
    install_dep "Consul" "1.11.4" "$VTROOT/dist/consul" install_consul
  fi

  # chromedriver
  if [ "$BUILD_CHROME" == 1 ] ; then
    install_dep "chromedriver" "90.0.4430.24" "$VTROOT/dist/chromedriver" install_chromedriver
  fi

  echo
  echo "bootstrap finished - run 'make build' to compile"
}

install_all
