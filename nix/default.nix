{ lib, buildGoModule }:
let
  version = "unstable";
in
buildGoModule {
  pname = "catppuccin-cli";
  inherit version;

  src = lib.fileset.toSource {
    root = ../.;
    fileset = lib.fileset.intersection (lib.fileset.fromSource (lib.sources.cleanSource ../.)) (
      lib.fileset.unions [
        ../go.mod
        ../go.sum
        ../main.go
        ../commands
        ../query
        ../shared
      ]
    );
  };

  vendorHash = "sha256-kwTPdktn1p4E/PvJ5kLDZZ0cZdx3u6gvFaPvtaeO8nw=";

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];

  postInstall = ''
    mv $out/bin/cli $out/bin/ctp
  '';

  meta.mainPackage = "ctp";
}
