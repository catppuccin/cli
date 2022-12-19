{

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
  };

  outputs = { self, nixpkgs }: 
    let pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in {
      defaultPackage.x86_64-linux = 
        pkgs.stdenv.mkDerivation {
          name = "ctp";
          src = self;
          buildInputs = with pkgs; [ go gopls ];
          buildPhase = "go mod tidy; go build cmd/ctp";
          installPhase = "mkdir -p $out/bin; install -t $out/bin ctp";
        };
    };
}
