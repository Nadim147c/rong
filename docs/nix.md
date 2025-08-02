# NixOS or Home Manager

As `Rong` is very experimental, this hasn't been added to nixpkgs. You have create a
custom package.

I'm currently working on creating modules for NixOS and Home Manager. For now, you
can use this following package snippet to create an overlay.

::: info
If you know how to create NixOS modules, consider contributing to the [GitHub
repository](https://github.com/Nadim147c/rong).
:::

```nix
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
    installShellCompletion --cmd $out/bin/rong \
        --bash <(echo "$bashComp") \
        --fish <(echo "$fishComp") \
        --zsh <(echo "$zshComp")
  '';

  meta = {
    description = "A Material You color generator";
    homepage = "https://github.com/Nadim147c/rong";
    license = lib.licenses.gpl3Only;
    mainProgram = "rong";
  };
}
```

::: info
Update the **hash** and **revision** to latest commit for up-to-date features.
:::
