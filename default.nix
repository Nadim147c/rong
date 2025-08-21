{
  buildGoModule,
  fetchFromGitHub,
  ffmpeg,
  installShellFiles,
  lib,
  stdenv,
  makeWrapper,
}:
buildGoModule {
  pname = "rong";
  version = "0-unstable-2025-07-27";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "96cf9e6ee50d6a4c06eaf1cb2aad028f5c251ab0";
    hash = "sha256-QuqN5xbSrhrI/Nw3D/cdRAlXQbpm/gQsiaLUPpzN2NE=";
  };

  vendorHash = "sha256-85QIPgacYq0QYP3WCuNGWYxbOSMRHct0pELL6lur1kU=";

  ldflags = ["-s" "-w"];

  nativeBuildInputs = [installShellFiles makeWrapper];
  propagatedBuildInputs = [ffmpeg];

  postInstall = lib.optionalString (stdenv.buildPlatform.canExecute stdenv.hostPlatform) ''
    installShellCompletion --cmd rong \
      --bash <($out/bin/rong _carapace bash) \
      --fish <($out/bin/rong _carapace fish) \
      --zsh <($out/bin/rong _carapace zsh)

    wrapProgram $out/bin/rong \
      --prefix PATH : ${lib.makeBinPath [ffmpeg]}
  '';

  meta = {
    description = "A Material You color generator";
    homepage = "https://github.com/Nadim147c/rong";
    license = lib.licenses.gpl3Only;
    mainProgram = "rong";
  };
}
