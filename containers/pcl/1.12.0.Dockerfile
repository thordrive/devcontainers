FROM thordev/devcontainers/cpp:ubuntu-21.04

RUN sudo -u thor ${VCPKG_ROOT}/vcpkg install pcl:x64-linux-release \
	&& rm -rf ${VCPKG_ROOT}/buildtrees \
	&& rm -rf ${VCPKG_ROOT}/packages \
	&& rm -rf ${VCPKG_DOWNLOADS}/*
