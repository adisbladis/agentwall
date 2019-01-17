with import <nixpkgs> {};

mkShell {
  buildInputs = [
    gnome3.zenity
    vgo2nix
    golint
    gocode
    go
  ];

  shellHook = ''
    unset GOPATH
  '';
}
