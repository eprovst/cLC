#![allow(non_snake_case)]
use lambdacalc::parser::parse_lambda_term;
use std::collections::HashMap;

use rustyline::error::ReadlineError;
use rustyline::Editor;

fn main() {
    let mut rl = Editor::<()>::new();

    loop {
        let readline = rl.readline("(cLC) ");
        match readline {
            Ok(line) => {
                rl.add_history_entry(line.as_str());
                match parse_lambda_term(&line.as_str(), &HashMap::new()) {
                    Ok(mut lam) => {
                        println!("\n{} =\n", lam);
                        lam.normal_order_reduce();
                        println!("    {}\n", lam);
                    }
                    Err(e) => println!("Error: {}", e),
                }
            }
            Err(ReadlineError::Interrupted) => {
                // TODO: stop computation if any is running, else break
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
