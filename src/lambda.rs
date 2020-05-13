use std::collections::HashMap;
use std::collections::HashSet;
use std::mem;
use std::ptr;

#[derive(Clone)]
pub enum LambdaTerm {
    BoundVariable(usize),
    FreeVariable(String),
    Application(Box<LambdaTerm>, Box<LambdaTerm>),
    Abstraction(Box<LambdaTerm>),
}

use self::LambdaTerm::*;

#[allow(dead_code)] // TODO: remove in the future
impl LambdaTerm {
    pub fn alpha_eq(&self, other: &LambdaTerm) -> bool {
        match self {
            BoundVariable(x) => match other {
                BoundVariable(y) => *x == *y,
                _ => false,
            },
            FreeVariable(x) => match other {
                FreeVariable(y) => *x == *y,
                _ => false,
            },
            Application(x1, x2) => match other {
                Application(y1, y2) => x1.alpha_eq(y1) && x2.alpha_eq(y2),
                _ => false,
            },
            Abstraction(x) => match other {
                Abstraction(y) => x.alpha_eq(y),
                _ => false,
            },
        }
    }

    fn free_variables(&self) -> HashSet<String> {
        match self {
            BoundVariable(_) => HashSet::with_capacity(0),
            FreeVariable(x) => {
                let mut set = HashSet::with_capacity(1);
                set.insert(x.clone());
                set
            }
            Application(x1, x2) => {
                let left_vars = x1.free_variables();
                let right_vars = x2.free_variables();
                left_vars.union(&right_vars).cloned().collect()
            }
            Abstraction(x) => x.free_variables(),
        }
    }

    fn shift_index(&mut self, inc: bool, cutoff: usize) {
        match self {
            BoundVariable(x) => {
                if *x >= cutoff {
                    if inc {
                        *self = BoundVariable(*x + 1);
                    } else {
                        *self = BoundVariable(*x - 1);
                    }
                }
            }
            Application(x1, x2) => {
                x1.shift_index(inc, cutoff);
                x2.shift_index(inc, cutoff);
            }
            Abstraction(x) => x.shift_index(inc, cutoff + 1),
            _ => {}
        }
    }

    fn heighten_index(&mut self) {
        self.shift_index(true, 0);
    }

    fn lower_index(&mut self) {
        self.shift_index(false, 1);
    }

    fn substitute(&mut self, variable: usize, mut other: LambdaTerm) {
        match self {
            BoundVariable(x) => {
                if *x == variable {
                    *self = other
                }
            }
            Application(x1, x2) => {
                // TODO: optimise clone away in obvious cases.
                x1.substitute(variable, other.clone());
                x2.substitute(variable, other);
            }
            Abstraction(x) => {
                other.heighten_index();
                x.substitute(variable + 1, other);
            }
            _ => {}
        }
    }

    pub fn beta_reduce(&mut self) {
        if let Application(abstr, arg_d) = self {
            if let Abstraction(ref mut body_d) = **abstr {
                // Create cheap objects to be destroyed together with the
                // application and abstraction and save the relevant bits.
                let mut arg = Box::new(BoundVariable(0));
                mem::swap(&mut arg, arg_d);
                let mut body = Box::new(BoundVariable(0));
                mem::swap(&mut body, body_d);

                arg.heighten_index();
                body.substitute(0, *arg);
                body.lower_index();
                *self = *body;
            }
        }
    }

    fn contains(&self, variable: usize) -> bool {
        match self {
            BoundVariable(x) => *x == variable,
            FreeVariable(_) => false,
            Application(x, y) => x.contains(variable) || y.contains(variable),
            Abstraction(x) => x.contains(variable),
        }
    }

    pub fn eta_reduce(&mut self) {
        match self {
            Application(x, y) => {
                x.eta_reduce();
                y.eta_reduce();
            }
            Abstraction(x) => {
                x.eta_reduce();

                if let Application(ref mut body_d, ref z) = **x {
                    if let BoundVariable(0) = **z {
                        if !body_d.contains(0) {
                            // Create cheap objects to be destroyed with the
                            // application and save the relevant bits.
                            let mut body = Box::new(BoundVariable(0));
                            mem::swap(&mut body, body_d);

                            body.eta_reduce();
                            body.lower_index();
                            *self = *body;
                        }
                    }
                }
            }
            _ => {}
        }
    }

    pub fn whnf(&mut self) {
        if let Abstraction(_) = self {
            // Already in WHNF.
        } else {
            unsafe {
                let body = ptr::read(self);
                ptr::write(
                    self,
                    Abstraction(Box::new(Application(
                        Box::new(body),
                        Box::new(BoundVariable(0)),
                    ))),
                );
            }
        }
    }

    fn can_reduce(&self) -> bool {
        match self {
            BoundVariable(_) => false,
            FreeVariable(_) => false,
            Application(x, y) => match **x {
                Abstraction(_) => true,
                _ => x.can_reduce() || y.can_reduce(),
            },
            Abstraction(x) => x.can_reduce(),
        }
    }

    fn normal_order_reduce_once(&mut self) {
        match self {
            Application(x, y) => {
                if let Abstraction(_) = **x {
                    self.beta_reduce();
                } else if x.can_reduce() {
                    x.normal_order_reduce_once();
                } else if y.can_reduce() {
                    y.normal_order_reduce_once();
                }
            }
            Abstraction(x) => {
                if x.can_reduce() {
                    x.normal_order_reduce_once();
                }
            }
            _ => {}
        }
    }

    /// WARNING: Might never stop
    pub fn normal_order_reduce(&mut self) {
        while self.can_reduce() {
            self.normal_order_reduce_once();
        }
        self.eta_reduce()
    }

    fn applicative_order_reduce_once(&mut self) {
        match self {
            Application(x, y) => {
                if y.can_reduce() {
                    y.applicative_order_reduce_once();
                } else if let Abstraction(_) = **x {
                    self.beta_reduce();
                } else if x.can_reduce() {
                    x.applicative_order_reduce_once();
                }
            }
            Abstraction(x) => {
                if x.can_reduce() {
                    x.normal_order_reduce_once();
                }
            }
            _ => {}
        }
    }

    /// WARNING: Might never stop
    pub fn applicative_order_reduce(&mut self) {
        while self.can_reduce() {
            self.applicative_order_reduce_once();
        }
        self.eta_reduce()
    }
}

fn variable_to_letter(i: usize) -> String {
    // x, y, z
    if i < 3 {
        ((120 + i) as u8 as char).to_string()
    // u, v, w
    } else if i < 6 {
        ((117 + i - 3) as u8 as char).to_string()
    // a, b, c ... t
    } else if i < 26 {
        ((97 + i - 6) as u8 as char).to_string()
    // x1, x2, x3 ...
    } else {
        String::from("x") + &(i - 25).to_string()
    }
}

struct FormatEnvironment<'a> {
    variable_offset: usize,
    variables: &'a mut HashMap<usize, String>,
    free_variables: &'a HashSet<String>,
    next_variable: usize,
}

fn fmt_aux(
    term: &LambdaTerm,
    env: &mut FormatEnvironment,
    f: &mut std::fmt::Formatter,
) -> Result<(), std::fmt::Error> {
    match term {
        FreeVariable(v) => f.write_str(v),
        BoundVariable(i) => {
            f.write_str(env.variables.get(&(*i + env.variable_offset - 1)).unwrap())
        }
        Application(x, y) => {
            if let Abstraction(_) = **x {
                f.write_str("(")?;
                fmt_aux(x, env, f)?;
                f.write_str(")")?;
            } else {
                fmt_aux(x, env, f)?;
            }
            f.write_str(" ")?;
            if let Abstraction(_) = **x {
                f.write_str("(")?;
                fmt_aux(y, env, f)?;
                f.write_str(")")
            } else {
                fmt_aux(y, env, f)
            }
        }
        Abstraction(x) => {
            f.write_str("λ")?;

            // Seach for variable that can be used
            let old_next_variable = env.next_variable;
            let mut var = variable_to_letter(env.next_variable);
            while env.free_variables.contains(&var) {
                env.next_variable += 1;
                var = variable_to_letter(env.next_variable);
            }

            f.write_str(var.as_str())?;
            env.variables.insert(env.variable_offset, var);

            f.write_str(".")?;

            env.variable_offset += 1;
            fmt_aux(x, env, f)?;
            env.variable_offset -= 1;

            env.variables.remove(&env.variable_offset);
            env.next_variable = old_next_variable;
            Ok(())
        }
    }
}

impl std::fmt::Display for LambdaTerm {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> Result<(), std::fmt::Error> {
        let mut env = FormatEnvironment {
            variable_offset: 0,
            variables: &mut HashMap::<usize, String>::new(),
            free_variables: &mut self.free_variables(),
            next_variable: 0,
        };
        fmt_aux(self, &mut env, f)
    }
}
