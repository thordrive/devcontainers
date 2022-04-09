FROM thordrive/dev-vtk:9.0.3-pv5.9.1

RUN apt-get update && apt-get install -y --no-install-recommends \
		bison \
		gperf \
		python3-distutils \
		autoconf \
		automake \
		libtool \
		libtool-bin \
		'^libxcb.*-dev' \
		libx11-xcb-dev \
		libglu1-mesa-dev \
		libxrender-dev \
		libxi-dev \
		libxkbcommon-dev \
		libxkbcommon-x11-dev \
	&& rm -rf /var/lib/apt/lists/*

RUN sudo -u thor ${VCPKG_ROOT}/vcpkg --binarysource=clear install qtbase \
	&& rm -rf ${VCPKG_ROOT}/buildtrees \
	&& rm -rf ${VCPKG_ROOT}/packages \
	&& rm -rf ${VCPKG_DOWNLOADS}/*

RUN sudo -u thor ${VCPKG_ROOT}/vcpkg --binarysource=clear install qtquick3d \
	&& rm -rf ${VCPKG_ROOT}/buildtrees \
	&& rm -rf ${VCPKG_ROOT}/packages \
	&& rm -rf ${VCPKG_DOWNLOADS}/*
