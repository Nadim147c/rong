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
  version = "5.0.0";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "v${version}";
    hash = "sha256-yNW9I8+5qxMZ27NOH0/nD3/pb95gA6sJ5fW5kqkPmAs=";
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
