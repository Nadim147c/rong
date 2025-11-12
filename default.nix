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
  version = "3.0.1";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "v${version}";
    hash = "sha256-0/b0QL680rmdRI1IBv2yRoYTuHp8ep3HftwSRq7xBgQ=";
  };

  vendorHash = "sha256-RAEEDhFDqXZri3iJaBrEQ4NdY9rYP3o+V6LF1cLvH6Y=";

  ldflags = ["-s" "-w" "-X" "main.Version=${version}"];

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
