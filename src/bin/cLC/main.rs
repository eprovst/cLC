#![allow(non_snake_case)]
use lambdacalc::lambda::LambdaTerm;
use lambdacalc::parser::parse_lambda_term;
use std::collections::HashMap;

use crossbeam_channel::{bounded, select, Receiver, TryRecvError};

use std::thread;

use rustyline::error::ReadlineError;
use rustyline::Editor;

fn main() {
    let mut rl = Editor::<()>::new();

    let ctrlc_channel = match ctrlc_channel() {
        Ok(chan) => chan,
        Err(_) => bounded(0).1,
    };

    print_info();
    loop {
        let readline = rl.readline("(cLC) ");
        match readline {
            Ok(line) => {
                rl.add_history_entry(line.as_str());
                let line = line.trim();

                if line == "help" {
                    print_help();
                } else if line == "info" {
                    print_info();
                } else if line == "exit" {
                    break;
                } else {
                    match parse_lambda_term(&line, &HashMap::new()) {
                        Ok(lam) => {
                            println!("\n{} =\n", lam);
                            print!("    ");
                            match reduce(&lam, &ctrlc_channel) {
                                Ok(res) => println!("{}", res),
                                Err(error) => println!("{}", error),
                            }
                            println!();
                        }
                        Err(e) => println!("\nError: {}\n", e),
                    }
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

fn reduce(term: &LambdaTerm, ctrlc_channel: &Receiver<()>) -> Result<LambdaTerm, &'static str> {
    let (res_ntx, res_rx) = bounded(2);
    let res_atx = res_ntx.clone();

    let (abt_ntx, abt_nrx) = bounded(1);
    let (abt_atx, abt_arx) = bounded(1);

    let mut n_term = term.clone();
    let mut a_term = term.clone();

    thread::spawn(move || {
        while n_term.can_reduce() {
            match abt_nrx.try_recv() {
                Err(TryRecvError::Empty) => n_term.normal_order_reduce_once(),
                Ok(_) | Err(TryRecvError::Disconnected) => return,
            }
        }
        n_term.eta_reduce();
        let _ = res_ntx.send(n_term);
    });

    thread::spawn(move || {
        while a_term.can_reduce() {
            match abt_arx.try_recv() {
                Err(TryRecvError::Empty) => a_term.normal_order_reduce_once(),
                Ok(_) | Err(TryRecvError::Disconnected) => return,
            }
        }
        a_term.eta_reduce();
        let _ = res_atx.send(a_term);
    });

    let res = select! {
        recv(res_rx) -> res => Ok(res.unwrap()),
        recv(ctrlc_channel) -> _ => Err("Computation interrupted."),
    };

    let _ = abt_ntx.send(());
    let _ = abt_atx.send(());
    res
}

fn print_info() {
    println!(
        r#"               _      ___
           __   \    /
          /     /\  (
          \__  /  \  \___

commandline Lambda Calculator v1.4.0
------------------------------------

Copyright (c) 2017-2020 Evert Provoost.
Some rights reserved.
"#
    );
}

// TODO: Implement all these...
fn print_help() {
    println!(
        r#"
Help:
-----
For full details: visit the project's wiki.

Available commands:

<lambda expression>
→ Normal order and applicative order expansion are tried for
  the expression, if there's a result it will be shown.

let <new global> = <lambda expression>
→ If the expansion can be fully reduced sets the global equal
  to that reduced form.

free <global1> <global2> <...>
→ Unbinds the global(s) and thus makes it a free variable.

match <lambda expression> with <global1> <global2> <...>
→ Tries to fully expand the expression and then shows the first
  listed global variable which is equivalent to that reduction.

weak <lambda expression>
→ Transforms the expression to a weak head normal form, then
  shows the result. Useful for expressions which wouldn't terminate
  reducing otherwise.

wlet <new global> = <lambda expression>
→ Equivalent to let but only transforms the expression to a weak
  head normal form.

<command> -- <comment>
→ Everything after -- is ignored.

help
→ Shows help for the cLC.

info
→ Shows information about the cLC.

exit
→ Closes the cLC.
"#
    );
}
