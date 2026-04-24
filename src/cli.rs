use clap::Parser;

/// GWTR - simple tool that enchanes the usage of git worktree system
#[derive(Debug, Parser)]
#[command(version,about,long_about = None)]
pub struct CLI {

    ///Set the application logs verbosity level, by default no logs are displayed
    #[arg(short,long, default_value_t = 0)]
    pub verbosity: u8,

    ///If arguments are provided, they will be passed through to git worktree command
    pub command_args: Option<Vec<String>>,
}
