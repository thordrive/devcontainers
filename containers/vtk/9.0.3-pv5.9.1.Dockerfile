FROM thordrive/dev-cpp:ubuntu-21.04

RUN apt-get update && apt-get install -y --no-install-recommends \
		libxmu-dev \
		libxi-dev \
		libgl-dev \
		libgl1-mesa-dev \
        && rm -rf /var/lib/apt/lists/*

RUN sudo -u thor ${VCPKG_ROOT}/vcpkg --binarysource=clear install vtk \
	&& rm -rf ${VCPKG_ROOT}/buildtrees \
	&& rm -rf ${VCPKG_ROOT}/packages \
	&& rm -rf ${VCPKG_DOWNLOADS}/*
