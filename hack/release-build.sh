#!/bin/sh
# Builds the files of a release into dist/: a static binary per platform, install.sh, a
# CycloneDX SBOM and sha256sums.txt. The release workflow runs it; run it yourself to see
# what a release would contain.
#
#   hack/release-build.sh VERSION      (VERSION like 0.1.0)
set -eu

version=${1:?usage: hack/release-build.sh VERSION}
echo "$version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' || { echo "version must look like 0.1.0, not $version" >&2; exit 2; }
name=jenkins-upgrade-check
module=$(go list -m)

rm -rf dist
mkdir dist
for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do
    os=${target%/*} arch=${target#*/} ext=''
    [ "$os" != windows ] || ext=.exe
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -trimpath -ldflags "-s -w -X main.Version=$version" \
        -o "dist/${name}_${os}_$arch$ext" ./cmd/$name
done
cp install.sh dist/

# The SBOM lists what is compiled in: the Go standard library (its version decides which
# Go vulnerabilities apply) and any third-party module.
go_version=$(go env GOVERSION | sed 's/^go//')
# A random serial number per SBOM, as CycloneDX asks; actions/attest refuses an SBOM without one.
uuid() {
    if [ -r /proc/sys/kernel/random/uuid ]; then cat /proc/sys/kernel/random/uuid
    else uuidgen | tr 'A-F' 'a-f'
    fi
}
go list -deps -f '{{with .Module}}{{if .Version}}{{.Path}} {{.Version}}{{end}}{{end}}' ./cmd/$name | sort -u |
    jq -R -s --arg name "$name" --arg module "$module" --arg version "$version" --arg go "$go_version" \
        --arg time "$(git log -1 --format=%cI)" --arg serial "urn:uuid:$(uuid)" '
    def comp($type; $name; $version; $purl; $license):
      {type: $type, "bom-ref": $purl, name: $name, version: $version, purl: $purl}
      + (if $license then {licenses: [{license: {id: $license}}]} else {} end);
    ("pkg:golang/\($module)@v\($version)") as $self
    | ([comp("library"; "stdlib"; $go; "pkg:golang/stdlib@\($go)"; "BSD-3-Clause")]
       + [split("\n")[] | select(length > 0) | split(" ") | comp("library"; .[0]; .[1]; "pkg:golang/\(.[0])@\(.[1])"; null)]) as $deps
    | {bomFormat: "CycloneDX", specVersion: "1.6", serialNumber: $serial, version: 1,
       metadata: {timestamp: $time,
         component: (comp("application"; $name; $version; $self; "Apache-2.0")
           + {supplier: {name: "keelinfra", url: ["https://keelinfra.io"]}})},
       components: $deps,
       dependencies: [{ref: $self, dependsOn: [$deps[]["bom-ref"]]}]}' \
    > "dist/$name.cdx.json"
jq -e '.bomFormat == "CycloneDX" and .specVersion == "1.6" and (.serialNumber | test("^urn:uuid:[0-9a-f-]{36}$"))' \
    "dist/$name.cdx.json" > /dev/null || { echo "dist/$name.cdx.json is not a CycloneDX SBOM actions/attest accepts" >&2; exit 1; }

if command -v sha256sum > /dev/null; then sha256='sha256sum'; else sha256='shasum -a 256'; fi
(cd dist && $sha256 "${name}"_* install.sh "$name.cdx.json" > sha256sums.txt)
ls -l dist
cat dist/sha256sums.txt
