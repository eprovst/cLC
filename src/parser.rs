use std::collections::HashMap;

use super::lambda::{LambdaTerm, LambdaTerm::*};

struct ParsingEnvironment<'a> {
    bound_variables: &'a mut Vec<String>,
    constants: &'a HashMap<String, LambdaTerm>,
}

fn skip_space(i: &str) -> &str {
    i.trim_start()
}

fn split_at(sep: char, i: &str) -> Result<(&str, &str), &'static str> {
    let mut iter = i.splitn(2, sep);

    if let (Some(l), Some(r)) = (iter.next(), iter.next()) {
        Ok((l, r))
    } else {
        Err("Can't split at whitespace.")
    }
}

// TODO: should gobble whitespace...
fn split_space(i: &str) -> Result<(&str, &str), &'static str> {
    let mut iter = i.splitn(2, |c: char| c.is_whitespace());

    if let (Some(l), Some(r)) = (iter.next(), iter.next()) {
        Ok((l, r))
    } else {
        Err("Can't split at whitespace.")
    }
}

// TODO
fn parse_braces(i: &str) -> Result<(&str, &str), &'static str> {
    Ok((i, ""))
}

fn parse_identifier<'a>(i: &'a str) -> Result<(&'a str, &'a str), &'static str> {
    let (i, rem) = split_space(i)?;
    let i = trim_braces_and_space(i);

    if i.is_empty() {
        Err("Expected variable.")
    } else if i.starts_with('\\') || i.starts_with('λ') {
        Err("Variable can't start with lambda.")
    } else if i.chars().any(|c| c.is_whitespace() || c == '(' || c == ')') {
        Err("Variable can't contain whitespace or braces.")
    } else {
        Ok((i, rem))
    }
}

fn parse_variable<'a>(
    i: &'a str,
    env: &ParsingEnvironment,
) -> Result<(LambdaTerm, &'a str), &'static str> {
    let (i, rem) = parse_identifier(i)?;

    match env.bound_variables.iter().rposition(|x| *x == i) {
        Some(idx) => Ok((BoundVariable(idx), rem)),
        None => match env.constants.get(&i.to_string()) {
            Some(cons) => Ok((cons.clone(), rem)),
            None => Ok((FreeVariable(i.to_string()), rem)),
        },
    }
}

// TODO
fn parse_abstraction<'a>(
    i: &'a str,
    env: &mut ParsingEnvironment,
) -> Result<(LambdaTerm, &'a str), &'static str> {
    if i.starts_with('\\') || i.starts_with('λ') {
        // TODO: in future, use strip_prefix
        let mut iter = i.chars();
        iter.next();
        let i = iter.as_str();

        let (var, rem) = parse_identifier(i)?;
        env.bound_variables.push(var.to_string());
        // TODO: parse remainder
        env.bound_variables.pop();

        Ok((BoundVariable(0), rem))
    } else {
        Err("Abstraction should start with λ or \\.")
    }
}

fn parse_application<'a>(
    i: &'a str,
    env: &mut ParsingEnvironment,
) -> Result<(LambdaTerm, &'a str), &'static str> {
    let (i, s) = if i.starts_with('(') {
        parse_braces(i)?
    } else {
        split_space(i)?
    };

    let fst = parse_full_term(i, env)?;
    let (snd, rem) = parse_term(s, env)?;

    Ok((Application(Box::new(fst), Box::new(snd)), rem))
}

// TODO
fn parse_term<'a>(
    i: &'a str,
    env: &mut ParsingEnvironment,
) -> Result<(LambdaTerm, &'a str), &'static str> {
    if i.starts_with('(') {
        let (i, s) = parse_braces(i)?;
        let res = parse_full_term(i, env)?;
        Ok((res, s))
    } else {
        // TODO
        Ok((BoundVariable(0), ""))
    }
}

fn trim_braces_and_space(i: &str) -> &str {
    let mut i = i.trim();
    while i.starts_with('(') && i.ends_with(')') {
        i = i.trim_start_matches('(').trim_end_matches(')');
        i = i.trim();
    }
    i
}

fn parse_full_term<'a>(
    i: &'a str,
    env: &mut ParsingEnvironment,
) -> Result<LambdaTerm, &'static str> {
    let i = trim_braces_and_space(i);
    let (res, rem) = parse_term(i, env)?;

    if rem.is_empty() {
        Ok(res)
    } else {
        Err("Superfluous characters.")
    }
}
