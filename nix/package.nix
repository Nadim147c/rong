{
  bash,
  buildGoModule,
  fetchFromGitHub,
  ffmpeg,
  installShellFiles,
  lib,
  makeWrapper,
  stdenv,
}:
buildGoModule (finalAttrs: {
  pname = "rong";
  version = "5.2.0";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "v${finalAttrs.version}";
    hash = "sha256-3wfaB3zqLCK0rmN2VBFnBkyJuv3vnOCXUptcAcQohd4=";
  };

  vendorHash = "sha256-kTK8R1qugeMpy4GsM5eI0TKGY1YWSy/TCcsZ8JpWQJY=";

  ldflags = [
    "-s"
    "-w"
    "-X"
    "main.Version=${finalAttrs.version}"
  ];

  nativeBuildInputs = [
    installShellFiles
    makeWrapper
  ];

  propagatedBuildInputs = [
    ffmpeg
    bash
  ];

  postInstall = ''
    wrapProgram $out/bin/rong \
        --prefix PATH : ${lib.makeBinPath finalAttrs.propagatedBuildInputs}
  ''
  + lib.optionalString (stdenv.buildPlatform.canExecute stdenv.hostPlatform) ''
    installShellCompletion --cmd rong \
      --bash <($out/bin/rong _carapace bash) \
      --fish <($out/bin/rong _carapace fish) \
      --zsh <($out/bin/rong _carapace zsh)
  '';

  meta = {
    description = "Material You and Base16 color generator";
    homepage = "https://github.com/Nadim147c/rong";
    license = lib.licenses.gpl3Only;
    mainProgram = "rong";
  };
})
