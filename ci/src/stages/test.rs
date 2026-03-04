use dagger_sdk::{Directory, Query};
use eyre::WrapErr;

use crate::containers::go_builder;

/// Run go tests on all packages.
pub async fn run(client: &Query, source: Directory) -> eyre::Result<String> {
    let output = go_builder(client, source)
        .with_exec(vec!["go", "test", "-v", "./..."])
        .with_exec(vec!["sh", "-c", "echo 'test: all tests passed'"])
        .stdout()
        .await
        .wrap_err("test failed")?;

    Ok(output)
}
