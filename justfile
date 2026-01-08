name := "rong"
bin := "build" / name
tag := `git describe --long --tags --abbrev=7 | sed 's/^v//;s/\([^-]*-g\)/r\1/;s/-/./g'`
version := env("VERSION", tag)
user-local := `echo ~/.local`
prefix := absolute_path(env("PREFIX", user-local))

build-install:
    @just build
    @just install

install:
    install -Dm755 {{ bin }} "{{ prefix }}/bin/{{ name }}"
    install -Dm644 <({{ bin }} _carapace bash) "{{ prefix }}/share/bash-completion/completions/{{ name }}.bash"
    install -Dm644 <({{ bin }} _carapace fish) "{{ prefix }}/share/fish/vendor_completions.d/{{ name }}.fish"
    install -Dm644 <({{ bin }} _carapace zsh ) "{{ prefix }}/share/zsh/site-functions/_{{ name }}"

build:
    # Build for default os and arch
    @just compile "" ""

compile os arch:
    env \
      CGO_ENABLED=0 \
      GOOS="{{ os }}" \
      GOARCH="{{ arch }}" \
      go build -trimpath -ldflags '-s -w -X main.Version={{ version }}' -o {{ bin }}

version:
    @echo {{ version }}

lint:
    golangci-lint run ./...

test:
    gotestsum --format pkgname-and-test-fails

generate:
    go-enum --values --names --marshal --no-iota \
      --output-suffix _generated \
      -f ./internal/config/enums/enums.go

docs:
    bun i
    bun run dev
