#![allow(non_snake_case)]
use crossbeam_channel::{bounded, Receiver};
use rustyline::error::ReadlineError;
use rustyline::Editor;
use statement::{parse_statement, Statement};
use std::collections::HashMap;

mod statement;

fn main() {
    let mut rl = Editor::<()>::new();

    let ctrlc_channel = match ctrlc_channel() {
        Ok(chan) => chan,
        Err(_) => bounded(0).1,
    };

    let mut env = HashMap::new();

    Statement::Info.execute(&ctrlc_channel, &mut env);
    loop {
        let readline = rl.readline("(cLC) ");
        match readline {
            Ok(line) => {
                rl.add_history_entry(line.as_str());
                match parse_statement(&line, &env) {
                    Ok(stmt) => stmt.execute(&ctrlc_channel, &mut env),
                    Err(error) => println!("{}", error),
                }
            }
            Err(ReadlineError::Interrupted) => {
                break;
            }
            Err(ReadlineError::Eof) => {
                break;
            }
            Err(err) => {
                println!("Error: {:?}", err);
                break;
            }
        }
    }
}

fn ctrlc_channel() -> Result<Receiver<()>, ctrlc::Error> {
    let (sender, receiver) = bounded(100);
    ctrlc::set_handler(move || {
        let _ = sender.send(());
    })?;

    Ok(receiver)
}
