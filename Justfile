hm_release := env("HM_RELEASE", "master")
nix_output := "result/share/doc/home-manager/options.json"
parsed_output := "site/src/data/options-master.json"

fetch:
    nix build github:nix-community/home-manager/{{hm_release}}#docs-json --no-write-lock-file

parse:
    go run tools/nix-parser/main.go {{nix_output}} {{parsed_output}}

update: fetch parse
    @echo "✓ Options updates for {{hm_release}}"

[working-directory: './site']
dev:
    npm run dev

[working-directory: './site']
build:
    npm run build

[working-directory: './site']
preview:
    npm run preview
