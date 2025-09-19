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
  version = "0.0.7";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "v${version}";
    hash = "sha256-D8yHtwSrdk86n8t55QWwmzUPgMdrTM+3TWzYmyEgfQ0=";
  };

  vendorHash = "sha256-ZDPYUYG43+ryUVTxvspHVbNuwky8k6xp62id8rVkoSw=";

  ldflags = ["-s" "-w" "-X" "main.Version=v${version}"];

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
    description = "A Material You and Base16 color generator";
    homepage = "https://github.com/Nadim147c/rong";
    license = lib.licenses.gpl3Only;
    mainProgram = "rong";
  };
}
