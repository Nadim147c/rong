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
  version = "4.1.0";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "v${version}";
    hash = "sha256-CKtpAW4v60q9fHoj4fuGR0s38B/INNzEcLBQZ1TrCak=";
  };

  vendorHash = "sha256-+uD3aZtKDmC4ExpctwYIgiFycztEKHkAY4Rof1d4WlE=";

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

  postInstall = lib.optionalString (stdenv.buildPlatform.canExecute stdenv.hostPlatform) /* bash */ ''
    installShellCompletion --cmd rong \
      --bash <($out/bin/rong _carapace bash) \
      --fish <($out/bin/rong _carapace fish) \
      --zsh <($out/bin/rong _carapace zsh)

    wrapProgram $out/bin/rong \
      --prefix PATH : ${lib.makeBinPath [ ffmpeg ]}
  '';

  meta = {
    description = "A Material You and Base16 color generator";
    homepage = "https://github.com/Nadim147c/rong";
    license = lib.licenses.gpl3Only;
    mainProgram = "rong";
  };
}
