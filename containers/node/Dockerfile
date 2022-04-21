ARG TAG
ARG FROM
FROM ${FROM:-node:${TAG}}

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
	&& apt-get clean --yes \
	&& rm -rf /var/lib/apt/lists/* \
	&& echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
	&& locale-gen



ARG USERNAME=thor
RUN usermod -l thor -d /home/thor -m "node" \
	&& groupmod -n ${USERNAME} node \
	&& echo "${USERNAME}:${USERNAME}" | chpasswd \
	&& echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > "/etc/sudoers.d/${USERNAME}"

WORKDIR /home/thor
