FROM mcr.microsoft.com/devcontainers/go:1.22

# Remove the broken yarn apt repository (GPG key expired: NO_PUBKEY 62D54FD4003F6525)
RUN rm -f /etc/apt/sources.list.d/yarn.list \
    && apt-get update \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
