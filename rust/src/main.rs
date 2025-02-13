use std::thread::JoinHandle;
use std::sync::mpsc;

use clap::Parser;

mod config;
use config::Config;

fn main() -> Result<(), std::io::Error> {
  let config = Config::parse();

  // let file_descriptor = File::open(config.word_list_path)?;
  // let reader = BufReader::new(file_descriptor);

  // for line in reader.lines() {
  //   let line = line?;
  // }
  // 

  // создание канала
  let (sender, receiver) = mpsc::channel();
  let sender_clone = sender.clone();

  // поток, который отправляет payloads
  let reader = std::thread::spawn(move ||
    {
      sender_clone.send("hah").unwrap();
    }
  );

  let mut workers: Vec< JoinHandle<()> > = Vec::with_capacity(
    config.worker_count as usize
  );
  for i in 0..config.worker_count {
    let i = i as usize;
    // workers
    workers[i] = std::thread::spawn(move ||
      {
  
      }
    );
  }

  reader.join().unwrap();
  for worker in workers {
    worker.join().unwrap();
  }

  Ok(())
}
