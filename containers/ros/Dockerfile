ARG TAG
ARG FROM
FROM ${FROM:-ros:${TAG}}

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
        nano \
        vim-tiny \
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
	&& rm -rf /var/lib/apt/lists/*

ARG USERNAME=thor
ARG USER_UID=1000
ARG USER_GID=${USER_UID}

RUN groupadd ${USERNAME} -g ${USER_GID} \
	&& useradd -ms /bin/bash ${USERNAME} -u ${USER_UID} -g ${USER_GID} \
	&& printf "${USERNAME}:${USERNAME}" | chpasswd \
	&& adduser ${USERNAME} sudo \
	&& echo "${USERNAME} ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
	&& chmod 0440 /etc/sudoers.d/${USERNAME}

WORKDIR /home/thor


# Install GeographicLib
RUN curl --insecure -fsSL -o GeographicLib-1.49.tar.gz "https://sourceforge.net/projects/geographiclib/files/distrib/GeographicLib-1.49.tar.gz/download" \
    && tar -xzf GeographicLib-1.49.tar.gz \
    && mkdir GeographicLib-1.49/build \
    && cd GeographicLib-1.49/build \
    && cmake .. \
    && make install \
    && cd ../.. \
    && rm -rf GeographicLib-1.49 GeographicLib-1.49.tar.gz
