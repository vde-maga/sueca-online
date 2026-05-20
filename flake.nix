{
  description = "Ambiente de desenvolvimento para Sueca Online (Go + React)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          # Pacotes necessários para a Fase 1 (Backend)
          buildInputs = with pkgs; [
            go          # Compilador e ferramentas Go
            gopls       # Language Server para o seu editor (VSCode/Neovim)
            golangci-lint # Linter avançado para Go
            gotools     # Ferramentas extra como `goimports`
            tree
          ];

          # Variáveis de ambiente úteis
          shellHook = ''
            echo "🃏 Ambiente de desenvolvimento Sueca Online ativado!"
            echo "Versão do Go: $(go version)"
          '';
        };
      }
    );
}