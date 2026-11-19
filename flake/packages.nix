_: {
  perSystem = {
    lib,
    pkgs,
    self',
    ...
  }: {
    packages =
      {
        default = pkgs.buildGoModule {
          pname = "qezta-ssh";
          version = "0.0.0";
          src = ../.;
          vendorHash = "sha256-pvL5gCP0bnLk6r4Hs7ryC8VA4DqsLvc9B5bnysHKOOw=";
          # Static binary: the image then needs no libc in its closure.
          env.CGO_ENABLED = "0";
          meta.mainProgram = "ssh";
        };
      }
      // lib.optionalAttrs pkgs.stdenv.hostPlatform.isLinux {
        # OCI image for k8s. The tag defaults to the image's content hash, so a
        # consumer pinning `imageName:imageTag` rolls out on every real change.
        image = pkgs.dockerTools.buildLayeredImage {
          name = "qezta-ssh";
          config = {
            Entrypoint = [(lib.getExe self'.packages.default)];
            # Distroless-style nonroot UID; the server never needs privileges
            # because it binds an unprivileged port (23234).
            User = "65532:65532";
            Env = ["SSH_HOST_KEY_PATH=/data/ssh_host_ed25519_key"];
            ExposedPorts."23234/tcp" = {};
          };
        };
      };
  };
}
