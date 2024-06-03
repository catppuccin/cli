{
  go,
  gopls,
  goreleaser,
  callPackage,
  gcc
}:
let
  mainPkg = callPackage ./default.nix { };
in
mainPkg.overrideAttrs (oa: {
  nativeBuildInputs = [
    go
    gopls
    goreleaser
    gcc
  ] ++ (oa.nativeBuildInputs or [ ]);
})
