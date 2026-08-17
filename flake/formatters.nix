{inputs, ...}: {
  imports = [inputs.treefmt-nix.flakeModule];

  perSystem.treefmt = {
    projectRootFile = "flake.nix";
    settings.global = {
      excludes = [
        ".github/*"
      ];
    };

    flakeCheck = false;

    programs = {
      ## Nix
      alejandra.enable = true;
      deadnix.enable = true;
      statix.enable = true;
      ## Go
      gofumpt.enable = true;
      ## Markdown
      mdformat.enable = true;
    };
  };
}
