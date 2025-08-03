{
  lib,
  buildGoModule,
  installShellFiles,
  fetchFromGitHub,
  stdenv,
}:
buildGoModule {
  pname = "rong";
  version = "0-unstable-2025-07-27";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "9752a110a88d79242b77143474216ede75204a48";
    hash = "sha256-CFrnMc1sUMEsBnMcmxszqMIea87A2pbZXsa6V3ackmI=";
  };

  vendorHash = "sha256-gT5iAYcUif2PQO6lVJRfUjddeAJc5ZrHg5hmkLkZeME=";

  ldflags = ["-s" "-w"];

  nativeBuildInputs = [installShellFiles];

  postInstall = lib.optionalString (stdenv.buildPlatform.canExecute stdenv.hostPlatform) ''
    installShellCompletion --cmd rong \
        --bash <($out/bin/rong _carapace bash) \
        --fish <($out/bin/rong _carapace fish) \
        --zsh <($out/bin/rong _carapace zsh)
  '';

  meta = {
    description = "A Material You color generator";
    homepage = "https://github.com/Nadim147c/rong";
    license = lib.licenses.gpl3Only;
    mainProgram = "rong";
  };
}
