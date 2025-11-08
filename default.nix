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
  version = "3.0.0";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "v${version}";
    hash = "sha256-zK1cvq2VVQJ8gqaY8sWMQS0YxMmeRuQLoV68wKVHYKQ=";
  };

  vendorHash = "sha256-cK5oIVOiEis759z3cMenkLJGKpQ6iynqaoFpchwYbYg=";

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
