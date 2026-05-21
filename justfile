set windows-shell := ["powershell.exe", "-NoProfile", "-Command"]

binary_name := "SISR"
main_pkg := "./cmd/sisr"
dist_dir := "./dist/"
target_goos := env_var_or_default("GOOS", if os_family() == "windows" { "windows" } else { "linux" })
exe_ext := if target_goos == "windows" { ".exe" } else { "" }
mkdir_p := if os_family() == "windows" { "New-Item -ItemType Directory -Force" } else { "mkdir -p" }
rm_rf := if os_family() == "windows" { "Remove-Item -Recurse -Force -ErrorAction 0" } else { "rm -rf" }

version := env_var_or_default("VERSION", `git describe --tags --always`)
buildType := env_var_or_default("BUILD_TYPE", "Debug")
commit := `git rev-parse --short HEAD`
build_time := if os_family() == "windows" {
	`Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'`
} else {
	`date -u '+%Y-%m-%dT%H:%M:%SZ'`
}
ldflags := "-s -w -X github.com/Alia5/SISR/meta.Version=" + version + " -X github.com/Alia5/SISR/meta.Commit=" + commit + " -X github.com/Alia5/SISR/meta.Date=" + build_time
build_path := join(dist_dir, binary_name + exe_ext)

default:
	just --list

[working-directory: 'dist']
run *args:
	go run ../{{ main_pkg }} {{ args }}

[working-directory: 'dist']
run-built:
	{{ if os_family() == "windows" { "& './" + binary_name + exe_ext + "'" } else { "'./" + binary_name + exe_ext + "'" } }}

win-resource:
	go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest -64 -o cmd/sisr/resource.syso versioninfo.json

[arg("type", long="type", help="Build type (Debug/Release)")]
build-sdl type="Debug":
	{{
		if os_family() == "windows" {
			"if (!(Test-Path deps/SDL/build)) { cmake -S deps/SDL -B deps/SDL/build }"
		} else {
			"[ -d deps/SDL/build ] || cmake -S deps/SDL -B deps/SDL/build -DCMAKE_BUILD_TYPE=" + type
		}
	}}
	cmake --build deps/SDL/build --config {{ type }}
	{{ mkdir_p }} deps/SDL/build/lib
	{{ 
		if os_family() == "windows" {
			"Copy-Item -Path deps/SDL/build/" + type + "/* -Destination deps/SDL/build/lib -Recurse -Force"
		} else {
			"find deps/SDL/build -maxdepth 1 \\( -name '*.so*' -o -name '*.a' \\) -exec cp -a {} deps/SDL/build/lib/ \\;"
		} 
	}}
	{{ mkdir_p }} {{ dist_dir }}
	{{ 
		if target_goos == "windows" {
			"Copy-Item -Path deps/SDL/build/" + type + "/*.dll -Destination " + dist_dir + " -Force"
		} else {
			"find deps/SDL/build -maxdepth 1 -name '*.so*' -exec cp -a {} " + dist_dir + " \\;"
		} 
	}}

[arg("type", long="type", help="Build type (Debug/Release)")]
build-polyhook2 type="Debug":
	{{
		if os_family() == "windows" {
			"if (!(Test-Path deps/PolyHook2/build)) { cmake -S deps/PolyHook2 -B deps/PolyHook2/build }"
		} else {
			"[ -d deps/PolyHook2/build ] || cmake -S deps/PolyHook2 -B deps/PolyHook2/build -DCMAKE_BUILD_TYPE=" + type
		}
	}}
	cmake --build deps/PolyHook2/build --config {{ type }}
	{{ mkdir_p }} deps/PolyHook2/build/lib
	{{ 
		if os_family() == "windows" {
			"Copy-Item -Path deps/PolyHook2/build/" + type + "/* -Destination deps/PolyHook2/build/lib -Recurse -Force"
		} else {
			"find deps/PolyHook2/build -maxdepth 1 \\( -name '*.so*' -o -name '*.a' \\) -exec cp -a {} deps/PolyHook2/build/lib/ \\;"
		} 
	}}

[arg("type", long="type", help="Build type (Debug/Release)")]
build-deps type="Debug": (build-sdl type) (build-polyhook2 type)

build-sisr type="Debug": (build-deps type)
    {{ 
        if target_goos == "windows" {
            "just win-resource"
        } else { 
            "echo Skipping win-resource for non-windows target" 
        } 
    }}
    {{ mkdir_p }} {{ dist_dir }}
    go build -ldflags "{{ ldflags }}" -o {{ build_path }} {{ main_pkg }}

[arg("type", pattern="Debug|Release")]
build type="Debug": (build-sisr type)

clean-sdl:
    {{ rm_rf }} deps/SDL/build

clean-polyhook2:
    {{ rm_rf }} deps/PolyHook2/build

clean-deps: clean-sdl clean-polyhook2

clean: clean-deps
    -{{ rm_rf }} {{ dist_dir }}
    go clean

test:
    go test -v ./...

fmt:
    go fmt ./...

lint:
    golangci-lint run ./...
