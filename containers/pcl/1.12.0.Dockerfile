FROM thordrive/dev-cpp:ubuntu-21.04

RUN sudo -u thor ${VCPKG_ROOT}/vcpkg install pcl \
	&& rm -rf ${VCPKG_ROOT}/buildtrees \
	&& rm -rf ${VCPKG_ROOT}/packages \
	&& rm -rf ${VCPKG_DOWNLOADS}/*
