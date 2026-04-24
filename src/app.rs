use anyhow::{Result, Context};
use clap::Parser;
use log::{debug, info};

use crate::git::get_worktree_list;
use crate::cli::CLI;
use std::process::Command as ShellCommand;



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

    match arguments.command_args {
        Some(args) => {
            debug!("Detected arguments passage {args:?}");
            if args.len() == 1 {
                debug!("Single argument provided, checking it's a worktree reference");
                // check if the argument is a string or parsable as int
                match args[0].parse::<i32>() {
                    Ok(index) => {
                        // if this is an int -> get the worktree by the index in the database, and switch
                        debug!("Arugment is int, assuming it's worktree index");
                    },
                    Err(_) => {
                        // if this is a string -> get the worktree by the branch name, and switch cwd to it
                        debug!("Argument is not an int, assuming it's a branch name");
                    }
                }

                return Ok(());
            } else {
                debug!("Invoking passthrough to git worktree");
                // create final arguments list for the passthrough command
                let mut passthrough_args = vec!["worktree".to_string()];
                passthrough_args.extend(args);

                let status = ShellCommand::new("git")
                    .args(passthrough_args)
                    .status()
                    .context("Error while executing the passthrough command")?;

                if status.success() {
                    debug!("Passthrough command executed successfully, updating the worktree database");
                    //TODO update the worktree database
                } else {
                    debug!("Passthrough command failed with status code: {}", status.code().unwrap_or(-1));
                }

                return Ok(());
            }
        },
        None => {
           debug!("No command arguments provided, listing worktrees");

           let worktrees = get_worktree_list()?;

           if worktrees.is_empty() {
               println!("No worktrees found.");
               return Ok(());
           }

           //FIXME should be replaced with table output
           worktrees.iter().for_each(|worktree| {
               println!("Worktree path: {}", worktree.path);
               println!("Worktree branch: {}", worktree.branch);
               println!("Worktree commit: {}", worktree.commit);
               println!();
           });
        },
    }
   Ok(())
}
