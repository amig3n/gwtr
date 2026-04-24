use anyhow::{Result, Context};

mod git;

mod cli;

mod app;
use crate::app::run;

fn main() -> Result<()> {
    run()
        .context("Failed to run the application");
    
    Ok(())
}
