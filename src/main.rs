mod git;
mod cli;
mod app;

use anyhow::{Result, Context};
use crate::app::run;

fn main() -> Result<()> {
    run()
        .context("Failed to run the application");
    
    Ok(())
}
