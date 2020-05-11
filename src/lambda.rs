use std::collections::HashMap;

#[derive(Clone)]
pub enum LambdaTerm {
    BoundVariable(usize),
    FreeVariable(String),
    Application(Box<LambdaTerm>, Box<LambdaTerm>),
    Abstraction(Box<LambdaTerm>),
}

use self::LambdaTerm::*;

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
        // TODO: can we get rid of the clones...
        if let Application(x, y) = self {
            if let Abstraction(ref mut z) = **x {
                y.heighten_index();
                z.substitute(0, *y.clone());
                *self = *z.clone();
                self.lower_index();
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
        // TODO: get rid of clone
        match self {
            Application(x, y) => {
                x.eta_reduce();
                y.eta_reduce();
            }
            Abstraction(x) => {
                x.eta_reduce();

                if let Application(ref mut y, ref z) = **x {
                    if let BoundVariable(0) = **z {
                        if !y.contains(0) {
                            y.eta_reduce();
                            y.lower_index();
                            *self = *y.clone();
                        }
                    }
                }
            }
            _ => {}
        }
    }

    pub fn whnf(&mut self) {
        // TODO: get rid of clone
        match self {
            Abstraction(_) => {}
            _ => {
                *self = Abstraction(Box::new(Application(
                    Box::new(self.clone()),
                    Box::new(BoundVariable(0)),
                )))
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

impl LambdaTerm {
    fn sub_fmt(
        &self,
        variable_offset: usize,
        variables: &mut HashMap<usize, String>,
        next_variable: &mut usize,
        f: &mut std::fmt::Formatter,
    ) -> Result<(), std::fmt::Error> {
        match self {
            FreeVariable(v) => f.write_str(v),
            BoundVariable(i) => f.write_str(variable_to_letter(*i).as_str()),
            Application(x, y) => {
                if let Abstraction(_) = **x {
                    f.write_str("(");
                    x.sub_fmt(variable_offset, variables, next_variable, f);
                    f.write_str(")");
                } else {
                    x.sub_fmt(variable_offset, variables, next_variable, f);
                }
                f.write_str(" ");
                y.sub_fmt(variable_offset, variables, next_variable, f)
            }
            Abstraction(x) => {
                f.write_str("λ");
                let var = variable_to_letter(*next_variable);
                f.write_str(var.as_str());
                variables.insert(variable_offset, var);
                *next_variable += 1;
                f.write_str(".");
                x.sub_fmt(variable_offset + 1, variables, next_variable, f)
            }
        }
    }
}

impl std::fmt::Display for LambdaTerm {
    // Format highest level:
    fn fmt(&self, f: &mut std::fmt::Formatter) -> Result<(), std::fmt::Error> {
        self.sub_fmt(0, &mut HashMap::<usize, String>::new(), &mut 0, f)
    }
}
