{
  buildGoModule,
  fetchFromGitHub,
  ffmpeg,
  installShellFiles,
  lib,
  stdenv,
  makeWrapper,
}:
buildGoModule (finalAttr: {
  pname = "rong";
  version = "0.0.4";

  src = fetchFromGitHub {
    owner = "Nadim147c";
    repo = "rong";
    rev = "v${finalAttr.version}";
    hash = "sha256-Z3KMBrCFlscPBEfob9HJWAEIJTvxalRVFX1OYs8u8D4=";
  };

  vendorHash = "sha256-YYKn8RsqtoqEIlC+dyl8s6OsUVH1eZYZfNoYLJxGe4c=";

  ldflags = ["-s" "-w" "-X" "main.Version=v${finalAttr.version}"];

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
})
