ARG TAG
ARG FROM
FROM ${FROM:-thordrive/dev-vtk:${TAG}}

RUN sudo -u thor ${VCPKG_ROOT}/vcpkg --binarysource=clear install pcl[vtk] \
	&& rm -rf ${VCPKG_ROOT}/buildtrees \
	&& rm -rf ${VCPKG_ROOT}/packages \
	&& rm -rf ${VCPKG_DOWNLOADS}/*
