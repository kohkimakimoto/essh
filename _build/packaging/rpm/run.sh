#!/usr/bin/env bash
set -eu

indent() {
    local n="${1:-4}"
    local p=""
    for i in `seq 1 $n`; do
        p="$p "
    done;

    local c="s/^/$p/"
    case $(uname) in
      Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
      *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
    esac
}

echo "--> Running packaging script in '$DOCKER_IMAGE' container..."
echo "    PRODUCT_NAME: $PRODUCT_NAME"
echo "    PRODUCT_VERSION: $PRODUCT_VERSION"
echo "    COMMIT_HASH: $COMMIT_HASH"

echo "    Copying files..."

repo_dir=$(pwd)
platform=el${RHEL_VERSION}

cp -pr _build/packaging/rpm/SPECS $HOME/rpmbuild/
cp -pr _build/packaging/rpm/SOURCES $HOME/rpmbuild/
cp -pr _build/dist/${PRODUCT_NAME}_linux_amd64.zip $HOME/rpmbuild/SOURCES/${PRODUCT_NAME}_linux_amd64.zip

echo "    Building rpm..."
cd $HOME
rpmbuild \
    --define "_product_name ${PRODUCT_NAME}" \
    --define "_product_version ${PRODUCT_VERSION}" \
    --define "_rhel_version ${RHEL_VERSION}" \
    -ba rpmbuild/SPECS/${PRODUCT_NAME}.spec \
    | indent

echo "    Copying RPMs back to shared folder..."
cd $repo_dir

mkdir -p _build/dist/${platform}
cp -pr $HOME/rpmbuild/RPMS _build/dist/${platform}
cp -pr $HOME/rpmbuild/SRPMS _build/dist/${platform}
