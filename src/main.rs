use anyhow::{Result, Context};
use env_logger::Env;
use log::{debug, error, info, warn};
use std::process::Command as ShellCommand;

mod cli;
use crate::cli::CLI;
use clap::Parser;


struct GitWorktree {
    branch: String,
    commit: String,
    path: String,
}

impl GitWorktree {
    fn new() -> GitWorktree {
        GitWorktree {
            branch: String::from(""),
            commit: String::from(""),
            path: String::from(""),
        }
    }

    fn set_branch(&mut self, branch: &str) {
        self.branch = branch.to_string();
    }

    fn set_commit(&mut self, commit: &str) {
        self.commit = commit.to_string();
    }

    fn set_path(&mut self, path: &str) {
        self.path = path.to_string();
    }
}

fn get_worktree_list() -> anyhow::Result<Vec<GitWorktree>> {
    let raw_output = ShellCommand::new("git")
        .args(["worktree", "list", "--porcelain"])
        .output()
        .context("Cannot obtain the worktree list")?;

    let output = String::from_utf8(raw_output.stdout)
        .context("Cannot parse the worktree list output as UTF-8")?;

    let mut worktrees: Vec<GitWorktree> = Vec::new();

    // iterate over output to parse the worktree list
    let mut current_worktree = GitWorktree::new();
    for line in output.lines() {
        let line_length = line.len();

        if line_length == 0 {
            debug!("Pushing last block to worktree list");
            worktrees.push(current_worktree);
            current_worktree = GitWorktree::new();
            continue;
        }

        let line_prefix = line.split_whitespace()
            .next()
            .context("Cannot parse the worktree list output: {line}")?;

        match line_prefix{
            "worktree" => {
                let path = line.strip_prefix("worktree ")
                    .context("Cannot parse the worktree path: {line}")?;
                debug!("Parsed worktree path: {path}");
                current_worktree.set_path(path);
            },

            "HEAD" => {
                let commit = line.strip_prefix("HEAD ")
                    .context("Cannot parse the worktree commit: {line}")?;
                debug!("Parsed worktree commit: {commit}");
                current_worktree.set_commit(commit);
            },

            "branch" => {
                let branch = line.strip_prefix("branch ")
                    .context("Cannot parse the worktree branch: {line}")?;
                debug!("Parsed worktree branch: {branch}");
                current_worktree.set_branch(branch);
            },

            _ => {
                // getting here means that something gone terribly wrong -> panic
                panic!("Unexpected line in the worktree list output: {line}");
            }
        }
    }

    Ok(worktrees)
}


fn run() -> anyhow::Result<()> {
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

fn main() -> Result<()> {
    run()
        .context("Failed to run the application");
    
    Ok(())
}
