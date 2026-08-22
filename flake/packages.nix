_: {
  perSystem = {pkgs, ...}: {
    packages.default = pkgs.buildGoModule {
      pname = "qezta-ssh";
      version = "0.0.0";
      src = ../.;
      vendorHash = "sha256-pvL5gCP0bnLk6r4Hs7ryC8VA4DqsLvc9B5bnysHKOOw=";
      meta.mainProgram = "ssh";
    };
  };
}
