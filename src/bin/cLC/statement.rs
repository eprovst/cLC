use self::Statement::*;
use crossbeam_channel::{bounded, select, Receiver, TryRecvError};
use lambdacalc::lambda::LambdaTerm;
use lambdacalc::parser::{is_valid_identifier, parse_lambda_term};
use std::collections::HashMap;
use std::thread;

pub enum Statement {
    NoOp,
    Exit,
    Help,
    Info,
    //TODO: Load(String),
    Reduce(LambdaTerm),
    Weak(LambdaTerm),
    Let(String, LambdaTerm),
    WeakLet(String, LambdaTerm),
    Free(Vec<String>),
    Match(LambdaTerm, Vec<String>),
}

pub fn parse_statement(
    inpt: &str,
    env: &HashMap<String, LambdaTerm>,
) -> Result<Statement, &'static str> {
    // Remove comments and trim
    let inpt = inpt.splitn(2, "--").next().unwrap().trim();

    if inpt.is_empty() {
        Ok(NoOp)
    } else if inpt == "exit" {
        Ok(Exit)
    } else if inpt == "help" {
        Ok(Help)
    } else if inpt == "info" {
        Ok(Info)
    } else if inpt.starts_with("let") {
        let inpt = &inpt[3..]; // TODO: replace with strip_prefix
        match inpt.splitn(2, '=').collect::<Vec<&str>>()[..] {
            [var, term] => {
                let var = var.trim();
                if is_valid_identifier(var) {
                    match parse_lambda_term(term, env) {
                        Ok(term) => Ok(Let(var.to_string(), term)),
                        Err(error) => Err(error),
                    }
                } else {
                    Err("Not a valid variable.")
                }
            }
            _ => Err("No right hand side in let."),
        }
    } else if inpt.starts_with("wlet") {
        let inpt = &inpt[4..]; // TODO: replace with strip_prefix
        match inpt.splitn(2, '=').collect::<Vec<&str>>()[..] {
            [var, term] => {
                let var = var.trim();
                if is_valid_identifier(var) {
                    match parse_lambda_term(term, env) {
                        Ok(term) => Ok(WeakLet(var.to_string(), term)),
                        Err(error) => Err(error),
                    }
                } else {
                    Err("Not a valid variable.")
                }
            }
            _ => Err("No right hand side in wlet."),
        }
    } else if inpt.starts_with("free") {
        let inpt = &inpt[4..]; // TODO: replace with strip_prefix
        let vars: Vec<String> = inpt.split_whitespace().map(|s| s.to_string()).collect();
        if vars.is_empty() {
            Err("No globals given.")
        } else if vars.iter().any(|v| !is_valid_identifier(v)) {
            Err("Invalid identifier in list.")
        } else {
            Ok(Free(vars))
        }
    } else if inpt.starts_with("match") {
        let inpt = &inpt[5..]; // TODO: replace with strip_prefix
        match inpt.splitn(2, "with").collect::<Vec<&str>>()[..] {
            [term, vars] => {
                let vars: Vec<String> = vars.split_whitespace().map(|s| s.to_string()).collect();
                if vars.is_empty() {
                    Err("No globals given.")
                } else if vars.iter().any(|v| !is_valid_identifier(v)) {
                    Err("Invalid identifier in list.")
                } else {
                    match parse_lambda_term(term, env) {
                        Ok(term) => Ok(Match(term, vars)),
                        Err(error) => Err(error),
                    }
                }
            }
            _ => Err("No with in match."),
        }
    } else if inpt.starts_with("weak") {
        let inpt = &inpt[4..]; // TODO: replace with strip_prefix
        match parse_lambda_term(inpt, env) {
            Ok(term) => Ok(Weak(term)),
            Err(error) => Err(error),
        }
    } else {
        match parse_lambda_term(inpt, env) {
            Ok(term) => Ok(Reduce(term)),
            Err(error) => Err(error),
        }
    }
}

impl Statement {
    pub fn execute(&self, ctrlc_channel: &Receiver<()>, env: &mut HashMap<String, LambdaTerm>) {
        match self {
            NoOp => {}
            Exit => std::process::exit(0),
            Help => print_help(),
            Info => print_info(),
            Reduce(lam) => {
                println!("\n{} =\n", lam);
                print!("    ");
                match reduce(&lam, &ctrlc_channel) {
                    Ok(res) => println!("{}", res),
                    Err(error) => println!("{}", error),
                }
                println!();
            }
            Weak(lam) => {
                println!("\n{} =\n", lam);
                let mut res = lam.clone();
                res.whnf();
                println!("    {}", res);
                println!();
            }
            Let(var, lam) => match reduce(&lam, &ctrlc_channel) {
                Ok(res) => {
                    env.insert(var.to_string(), res);
                }
                Err(error) => println!("{}", error),
            },
            WeakLet(var, lam) => {
                let mut lam = lam.clone();
                lam.whnf();
                env.insert(var.to_string(), lam);
            }
            Free(vars) => {
                for v in vars {
                    env.remove(v);
                }
            }
            Match(lam, vars) => {
                println!("\n{} =", lam);
                match reduce(&lam, &ctrlc_channel) {
                    Ok(lam) => {
                        println!("{} =\n", lam);
                        print!("    ");
                        for var in vars {
                            if let Ok(term) = parse_lambda_term(var, env) {
                                match reduce(&term, &ctrlc_channel) {
                                    Ok(term) => {
                                        if lam.alpha_eq(&term) {
                                            println!("{}", var);
                                            break;
                                        }
                                    }
                                    Err(error) => {
                                        println!("{}", error);
                                        return;
                                    }
                                }
                            }
                        }
                        println!();
                    }
                    Err(error) => println!("{}", error),
                }
            }
        }
    }
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

load <path/to/file>
→ Execute the given file in the current environment.

help
→ Shows help for the cLC.

info
→ Shows information about the cLC.

exit
→ Closes the cLC.
"#
    );
}
