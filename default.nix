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
    rev = "d959883b5cb23d44f83df65ff049e6426fc962fc";
    hash = "sha256-qHJ5KllX/Kb6DAAFMNpPUO/UC0zoAhXLDdxQnnX8vwQ=";
  };

  vendorHash = "sha256-BGjNmaFrugJt7aYIrXfr5UB76haSZBTxfRGPTElPjPQ=";

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
