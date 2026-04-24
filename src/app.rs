use anyhow::{Result, Context};
use clap::Parser;

use crate::git::get_worktree_list;
use crate::cli::CLI;

pub fn run() -> anyhow::Result<()> {
    let arguments = CLI::parse();

    let log_level = match arguments.verbosity {
        0 => log::LevelFilter::Off,
        1 => log::LevelFilter::Info,
        2 => log::LevelFilter::Debug,
        _ => log::LevelFilter::Trace,
    };

    env_logger::Builder::new()
        .filter_level(log_level)
        .init();

   let worktrees = get_worktree_list()?;

   if worktrees.is_empty() {
       println!("No worktrees found.");
       return Ok(());
   }

   worktrees.iter().for_each(|worktree| {
       println!("Worktree path: {}", worktree.path);
       println!("Worktree branch: {}", worktree.branch);
       println!("Worktree commit: {}", worktree.commit);
       println!();
   });

   Ok(())
}
