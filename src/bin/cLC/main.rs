#![allow(non_snake_case)]
use lambdacalc::parser::parse_lambda_term;
use std::collections::HashMap;

fn main() {
    let mut lam = parse_lambda_term("(\\a.c a b) c", &HashMap::new()).unwrap();
    println!("{}", lam);
    lam.normal_order_reduce();
    println!("{}", lam);
}
