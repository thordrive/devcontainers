FROM mcr.microsoft.com/vscode/devcontainers/cpp:0-ubuntu-21.04

RUN apt-get update && apt-get install -y --no-install-recommends \
		gcc-11 \
		g++-11 \
	&& rm -rf /var/lib/apt/lists/* \
	&& update-alternatives --remove-all cpp \
	&& update-alternatives \
		--install /usr/bin/gcc gcc /usr/bin/gcc-11 110 \
		--slave /usr/bin/g++ g++ /usr/bin/g++-11 \
		--slave /usr/bin/gcov gcov /usr/bin/gcov-11 \
		--slave /usr/bin/gcc-ar gcc-ar /usr/bin/gcc-ar-11 \
		--slave /usr/bin/gcc-ranlib gcc-ranlib /usr/bin/gcc-ranlib-11 \
		--slave /usr/bin/cpp cpp /usr/bin/cpp-11

RUN usermod -l thor -d /home/thor -m "vscode" \
	&& groupmod -n thor vscode \
	&& rm /etc/sudoers.d/vscode \
	&& echo "thor ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/thor
