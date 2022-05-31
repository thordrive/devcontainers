ARG TAG
ARG FROM
FROM ${FROM:-thordrive/dev-gcc:${TAG}}

RUN apt-get update \
	&& apt-get install --yes --no-install-recommends \
		libxmu-dev \
		libxi-dev \
		libgl-dev \
		libgl1-mesa-dev \
	&& rm -rf /var/lib/apt/lists/*

COPY --chown=1000:1000 "pcl1.11.1.vcpkg.json" "/tmp/vcpkg/vcpkg.json"

RUN cd /tmp/vcpkg \
	&& sudo -u thor vcpkg install \
	&& mv vcpkg_installed /opt/vcpkg/installed \
	&& rm -rf \
		/tmp/vcpkg \
		${VCPKG_ROOT}/buildtrees \
		${VCPKG_ROOT}/packages \
		${VCPKG_DOWNLOADS}/*
