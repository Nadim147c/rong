{
  buildGoModule,
  fetchFromGitHub,
  ffmpeg,
  installShellFiles,
  lib,
  makeWrapper,
  stdenv,
}:
buildGoModule rec {
  pname = "rong";
  version = "5.0.1";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "v${version}";
    hash = "sha256-iFLxFwgMfuoSClvkTegRhnfBRAtcaPgXCINGTDJhAZ4=";
  };

  vendorHash = "sha256-TMFfw5s/Y8wTHzlg6El0ksji/ryAjA/GF8vsHNzqrSE=";

  ldflags = [
    "-s"
    "-w"
    "-X"
    "main.Version=${version}"
  ];

  nativeBuildInputs = [
    installShellFiles
    makeWrapper
  ];
  propagatedBuildInputs = [ ffmpeg ];

  postInstall = ''
    wrapProgram $out/bin/rong \
        --prefix PATH : ${lib.makeBinPath [ ffmpeg ]}
  ''
  ++ lib.optionalString (stdenv.buildPlatform.canExecute stdenv.hostPlatform) ''
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
}
