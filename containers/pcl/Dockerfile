ARG TAG
ARG FROM
FROM ${FROM:-thordrive/dev-gcc:${TAG}}

ARG VCPKG_COMMIT="origin/master"
RUN cd "${VCPKG_ROOT}" \
	&& git checkout "${VCPKG_COMMIT}"

RUN apt-get update \
	&& apt-get install --yes --no-install-recommends \
		libxmu-dev \
		libxi-dev \
		libgl-dev \
		libgl1-mesa-dev \
	&& rm -rf /var/lib/apt/lists/*
	
RUN sudo -u thor ${VCPKG_ROOT}/vcpkg --binarysource=clear install pcl[vtk] \
	&& rm -rf ${VCPKG_ROOT}/buildtrees \
	&& rm -rf ${VCPKG_ROOT}/packages \
	&& rm -rf ${VCPKG_DOWNLOADS}/*
