use std::str::FromStr;

use clap::Parser;

#[derive(Parser, Debug)]
// #[command(version, about, long_about = None)]
pub struct Config {
  #[arg(short, long)]
  pub target: String,

  #[arg(short, long)]
  pub word_list_path: String,

  #[arg(short, long)]
  // content of body
  pub payload: String,

  #[arg(short, long, default_value_t = String::from_str("").unwrap())]
  pub content_type: String,

  #[arg(short, long, default_value_t = String::from_str("POST").unwrap())]
  pub method: String,

  #[arg(long, default_value_t = 3)]
  pub worker_count: u32

  // TODO filter
}
