use dagger_sdk::{Container, Directory, Query};

/// Base Go container with module and build caches.
pub fn go_builder(client: &Query, source: Directory) -> Container {
    let mod_cache = client.cache_volume("forge-sdk-go-mod");
    let build_cache = client.cache_volume("forge-sdk-go-build");

    client
        .container()
        .from("golang:1.22-bookworm")
        .with_mounted_directory("/build", source)
        .with_workdir("/build")
        .with_mounted_cache("/go/pkg/mod", mod_cache)
        .with_mounted_cache("/root/.cache/go-build", build_cache)
}
