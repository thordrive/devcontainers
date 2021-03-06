ARG TAG
ARG FROM
FROM ${FROM:-gcc:${TAG}}

RUN apt-get update \
	&& apt-get install --yes --no-install-recommends \
		apt-utils \
		openssh-client \
		gnupg2 \
		dirmngr \
		iproute2 \
		procps \
		lsof \
		htop \
		net-tools \
		psmisc \
		curl \
		wget \
		rsync \
		ca-certificates \
		unzip \
		zip \
		git \
		git-lfs \
		nano \
		vim \
		less \
		jq \
		lsb-release \
		apt-transport-https \
		dialog \
		libc6 \
		libgcc1 \
		libkrb5-3 \
		libgssapi-krb5-2 \
		libicu[0-9][0-9] \
		liblttng-ust0 \
		libstdc++6 \
		zlib1g \
		locales \
		sudo \
		ncdu \
		man-db \
		strace \
		manpages \
		manpages-dev \
		init-system-helpers \
		bash-completion \
		zsh \
		python3 \
		pip \
	&& apt-get clean --yes \
	&& rm -rf /var/lib/apt/lists/* \
	&& echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
	&& locale-gen



ARG USERNAME=thor
ARG USER_UID=${UID:-1000}
ARG USER_GID=${GID:-1000}
RUN groupadd ${USER_GID} --gid ${USER_GID} \
	&& useradd --create-home ${USERNAME} --shell /bin/bash --uid ${USER_UID} --gid ${USER_GID} \
	&& echo "${USERNAME}:${USERNAME}" | chpasswd \
	&& echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > "/etc/sudoers.d/${USERNAME}"

WORKDIR /home/thor



RUN apt-get update \
	&& apt-get install --no-install-recommends --yes \
		gdb \
		build-essential \
		pkg-config \
		ninja-build \
		gcovr \
		clang-tidy \
	&& apt-get clean --yes \
	&& rm -rf /var/lib/apt/lists/*



# Install CMake
ARG CMAKE_VERSION
RUN TMP_DIR=$(mktemp -d -t cmake-${CMAKE_VERSION}-XXXX) \
	&& echo ${TMP_DIR} \
	&& cd ${TMP_DIR} \
	&& curl -sSL "https://github.com/Kitware/CMake/releases/download/v${CMAKE_VERSION}/cmake-${CMAKE_VERSION}-linux-x86_64.tar.gz" -O \
	&& curl -sSL "https://github.com/Kitware/CMake/releases/download/v${CMAKE_VERSION}/cmake-${CMAKE_VERSION}-SHA-256.txt" -o checksum \
	&& sha256sum --check --ignore-missing "checksum" \
	&& tar -xzf "cmake-${CMAKE_VERSION}-linux-x86_64.tar.gz" \
		--exclude bin/cmake-gui \
		--exclude doc/cmake \
		--exclude "share/cmake-*/Help" \
	&& cd */ \
	&& mkdir -p /opt/cmake \
	&& cp -rf * /opt/cmake/. \
	&& ln -s /opt/cmake/bin/cmake /usr/local/bin/cmake \
	&& rm -rf ${TMP_DIR}



# Install vcpkg
ENV VCPKG_ROOT=/opt/vcpkg \
	VCPKG_DOWNLOADS=/opt/vcpkg-downloads
RUN mkdir -p "${VCPKG_ROOT}" "${VCPKG_DOWNLOADS}" \
	&& git clone https://github.com/microsoft/vcpkg "${VCPKG_ROOT}" \
	&& cd "${VCPKG_ROOT}" \
	&& ./bootstrap-vcpkg.sh -disableMetrics \
	&& chown -R "${USER_UID}:${USER_GID}" "${VCPKG_ROOT}" "${VCPKG_DOWNLOADS}" \
	&& ln -s /opt/vcpkg/vcpkg /usr/local/bin/vcpkg \
	&& VCPKG_FORCE_SYSTEM_BINARIES=1 su "${USERNAME}" -c "${VCPKG_ROOT}/vcpkg integrate bash"

