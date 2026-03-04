use dagger_sdk::{Directory, Query};
use eyre::WrapErr;

use crate::containers::go_builder;

/// Run go vet on all packages.
pub async fn run(client: &Query, source: Directory) -> eyre::Result<String> {
    let output = go_builder(client, source)
        .with_exec(vec!["go", "vet", "./..."])
        .with_exec(vec!["sh", "-c", "echo 'check: go vet passed'"])
        .stdout()
        .await
        .wrap_err("check failed")?;

    Ok(output)
}
