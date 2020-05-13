// TODO: detect bracket borders without space...

use std::collections::HashMap;

use super::lambda::{LambdaTerm, LambdaTerm::*};

struct ParsingEnvironment<'a> {
    bound_variables: &'a mut Vec<String>,
    constants: &'a HashMap<String, LambdaTerm>,
}

fn trim_braces_and_space(i: &str) -> &str {
    let i = i.trim();

    if i.starts_with('(') {
        if let Ok((inner, "")) = parse_braces_r(i) {
            trim_braces_and_space(inner)
        } else {
            i
        }
    } else {
        i
    }
}

fn split_at(sep: char, i: &str) -> Result<(&str, &str), &'static str> {
    let i = i.trim();

    let mut iter = i.splitn(2, sep);

    if let (Some(l), Some(r)) = (iter.next(), iter.next()) {
        Ok((l, r))
    } else {
        Err("Missing character.")
    }
}

fn split_blob(i: &str) -> (&str, &str) {
    let i = i.trim();

    if let Some(idx) = i.find(|c: char| c.is_whitespace() || c == '(') {
        let (l, r) = i.split_at(idx);
        (l, r.trim_start())
    } else {
        (i, "")
    }
}

fn split_blob_r(i: &str) -> (&str, &str) {
    let i = i.trim();

    if let Some(idx) = i.rfind(|c: char| c.is_whitespace() || c == ')') {
        let (l, r) = i.split_at(idx + 1);
        (l.trim_end(), r)
    } else {
        (i, "")
    }
}

fn parse_braces_r(i: &str) -> Result<(&str, &str), &'static str> {
    let i = i.trim();

    if i.ends_with(')') {
        let mut iter = i.char_indices().rev();
        iter.next();

        let mut opens = 1;
        while opens > 0 {
            match iter.next() {
                Some((_, ')')) => opens += 1,
                Some((_, '(')) => opens -= 1,
                None => return Err("Unbalanced brackets."),
                _ => {}
            }
        }

        match iter.next() {
            Some((idx, _)) => {
                let (l, r) = i.split_at(idx + 1);
                Ok((l, &r[1..r.len() - 1]))
            }
            None => Ok((&i[1..i.len() - 1], "")),
        }
    } else {
        Err("No bracket found.")
    }
}

fn parse_identifier<'a>(i: &'a str) -> Result<(&'a str, &'a str), &'static str> {
    let i = i.trim();

    let (i, rem) = split_blob(i);

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
    let i = trim_braces_and_space(i);

    let (i, rem) = parse_identifier(i)?;

    match env.bound_variables.iter().rev().position(|x| *x == i) {
        Some(idx) => Ok((BoundVariable(idx), rem)),
        None => match env.constants.get(&i.to_string()) {
            Some(cons) => Ok((cons.clone(), rem)),
            None => Ok((FreeVariable(i.to_string()), rem)),
        },
    }
}

fn parse_abstraction<'a>(
    i: &'a str,
    env: &mut ParsingEnvironment,
) -> Result<(LambdaTerm, &'a str), &'static str> {
    let i = trim_braces_and_space(i);

    if i.starts_with('\\') || i.starts_with('λ') {
        // TODO: in future, use strip_prefix
        let mut iter = i.chars();
        iter.next();
        let i = iter.as_str();

        let (param, body) = split_at('.', i)?;

        // TODO: allow multiple parameters
        let param = param.trim_start();
        let (var, rem) = parse_identifier(param)?;

        if rem.is_empty() {
            env.bound_variables.push(var.to_string());
            let res = parse_term(body, env);
            env.bound_variables.pop();

            let (body, rem) = res?;
            Ok((Abstraction(Box::new(body)), rem))
        } else {
            Err("Superfluous characters in parameter list.")
        }
    } else {
        Err("Abstraction should start with λ or \\.")
    }
}

fn parse_application<'a>(
    i: &'a str,
    env: &mut ParsingEnvironment,
) -> Result<(LambdaTerm, &'a str), &'static str> {
    let i = trim_braces_and_space(i);

    let (i, s) = if i.ends_with(')') {
        parse_braces_r(i)?
    } else {
        split_blob_r(i)
    };

    let s = s.trim_start();

    if !s.is_empty() {
        let fst = parse_full_term(i, env)?;
        let (snd, rem) = parse_term(s, env)?;

        Ok((Application(Box::new(fst), Box::new(snd)), rem))
    } else {
        Err("No argument in application.")
    }
}

fn parse_term<'a>(
    i: &'a str,
    env: &mut ParsingEnvironment,
) -> Result<(LambdaTerm, &'a str), &'static str> {
    parse_abstraction(i, env)
        .or_else(|_| parse_application(i, env))
        .or_else(|_| parse_variable(i, env))
}

fn parse_full_term<'a>(
    i: &'a str,
    env: &mut ParsingEnvironment,
) -> Result<LambdaTerm, &'static str> {
    let (res, rem) = parse_term(i, env)?;

    if rem.trim().is_empty() {
        Ok(res)
    } else {
        Err("Superfluous characters.")
    }
}

#[allow(clippy::implicit_hasher)]
pub fn parse_lambda_term(
    i: &str,
    constants: &HashMap<String, LambdaTerm>,
) -> Result<LambdaTerm, &'static str> {
    let mut env = ParsingEnvironment {
        bound_variables: &mut Vec::new(),
        constants,
    };
    parse_full_term(i, &mut env)
}
