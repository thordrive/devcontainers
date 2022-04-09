FROM thordrive/dev-vtk:9.0.3-pv5.9.1

RUN sudo -u thor ${VCPKG_ROOT}/vcpkg --binarysource=clear install pcl[vtk] \
	&& rm -rf ${VCPKG_ROOT}/buildtrees \
	&& rm -rf ${VCPKG_ROOT}/packages \
	&& rm -rf ${VCPKG_DOWNLOADS}/*
