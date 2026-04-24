use std::{fs, fmt, env};
use log::{info,debug};
use std::path::Path;
use anyhow::{Result, Context};

const REPO_PATH: &str = ".git";

enum RepoState {
    Repo,
    Worktree
}

struct ScanResult {
    path: String,
    state: RepoState
}


/// Scan current directory and all upwards to detect git repository
fn scan_directory() -> anyhow::Result<ScanResult> {
    // read the directory content
    loop {
        let cwd = env::current_dir()?;
        let git_path = cwd.join(REPO_PATH);
        debug!("Looking for .git in: {}", cwd.display());

        // check if .git exists
        if git_path.exists() {
            if git_path.is_dir() {
                info!("Found .git directory at: {}", git_path.display());
                return Ok(ScanResult {
                    path: git_path.to_string_lossy().to_string(),
                    state: RepoState::Repo
                });
                
            }
            else if git_path.is_file() {
                info!("Found .git file at: {}", git_path.display());
                return Ok(ScanResult {
                    path: git_path.to_string_lossy().to_string(),
                    state: RepoState::Worktree
                });
            }
        }
        
        let parent_directory = cwd.parent()
            .context("Not possible to go up from root directory")?;

        env::set_current_dir(parent_directory)
            .context("Reached root directory without finding .git")?;
    }
}

/// Reads the .git file to find the actual git directory
fn get_git_dir(reposcan: ScanResult) -> anyhow::Result<String> {
    debug!("Parsing repository scan results");
    match reposcan.state {
        RepoState::Repo => {
            debug!("Obtained git path: {}", reposcan.path);
            Ok(reposcan.path)
        },
        RepoState::Worktree => {
            // open .git file
            debug!("Attempting to read .git file: {}", reposcan.path);
            let gitfile = fs::read_to_string(reposcan.path)
                .context("Cannot read .git file")?;

            let parsed_path = gitfile.strip_prefix("gitdir: ")
                    .context("failed obtaining path from .git file")?
                    .to_string();

            Ok(parsed_path)
        }
    }
}


fn run() -> Result<()> {
    env_logger::init();
    debug!("Looking for git repository files");
    let dirscan = scan_directory()?;
    let git_path = get_git_dir(dirscan)?;
    info!("obtained git path: {}", git_path);
    Ok(())
}

fn main() {
    match run() {
        Ok(()) => info!("Detected git repository"),
        Err(e) => eprintln!("Error: {}", e),
    }
}
