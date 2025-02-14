use std::io::{BufRead, BufReader};
use std::ptr::write_volatile;
use std::sync::mpsc::{channel, Sender};
use std::thread::JoinHandle;
use std::fs::{read, File};

use batch_channel::{bounded_sync, SyncReceiver, SyncSender};
use clap::Parser;

mod config;
use config::Config;




fn main() -> std::io::Result<()> {
  let config = Config::parse();

  // channel for fileReaderThread -> workers
  let (
    sender_to_worker,
    receiver_worker
  ) = bounded_sync::<String>(config.worker_count as usize);
  let (sender_from_worker, receiver_response) = channel::<Response>();



  let file_reader_thread = produce_file(config.word_list_path, sender_to_worker)?;

  let workers =  start_workers(
    config.worker_count as usize,
    receiver_worker,
    sender_from_worker
  );

  // TODO here get Responses and log it

  file_reader_thread.join().unwrap();
  join_all_workers(workers);

  Ok(())
}

// ! request
// TODO
fn request(body: String, request: Request) -> Response {
  Response{}
}

struct Response;
struct Request {
  content_type: String,
  method: String,
  target: String
}

impl Request {
  pub fn new(
    target: String,
    method: String,
    content_type: String
    // target: String
  ) -> Request {
    Request { content_type, method, target}
  }

  pub fn get_content_type(&self) -> &String {
    &self.content_type
  }

  pub fn get_method(&self) -> &String {
    &self.method
  }

  pub fn get_target(&self) -> &String {
    &self.target
  }
}


// ! workers
fn start_workers(count: usize,
                  receiver: SyncReceiver<String>,
                  sender: Sender<Response>)
  -> Vec<JoinHandle<()>>
{
  let mut workers: Vec< JoinHandle<()> > = Vec::with_capacity(count);

  for i in 0..count {
    let rec_clone = receiver.clone();
    let sender_clone = sender.clone();
    workers[i] = std::thread::spawn(move ||
      {
        worker(rec_clone, sender_clone);
      }
    );
  }

  workers
}

// TODO
fn worker(receiver: SyncReceiver<String>, sender: Sender<Response>) {
  let payload = receiver.recv().unwrap();
}

fn join_all_workers(workers: Vec<JoinHandle<()>>) {
  for worker in workers {
    worker.join().unwrap();
  }
}

// ! file thread
fn read_file(file_descriptor: File, sender: SyncSender<String>) {
  let reader = BufReader::new(file_descriptor);
  for line in reader.lines() {
    let line = line.unwrap();
    sender.send(line).unwrap();
  }
}

fn produce_file(path: String, sender_channel: SyncSender<String>)
  -> std::io::Result<JoinHandle<()>>
{
  let file_descriptor = File::open(path)?;
  let reader_thread = std::thread::spawn(move || {
      read_file(file_descriptor, sender_channel);
    }
  );

  Ok(reader_thread)
}
