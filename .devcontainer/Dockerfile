FROM mcr.microsoft.com/devcontainers/go:1.22@sha256:4e314a4eac046a0547f0c4cc5708b9caf51cb0fd6a5cfcd84e3a02f1a0dcdcab

# Remove the broken yarn apt repository (GPG key expired: NO_PUBKEY 62D54FD4003F6525)
RUN rm -f /etc/apt/sources.list.d/yarn.list \
    && apt-get update \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
