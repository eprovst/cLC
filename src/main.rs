mod lambda;

use lambda::LambdaTerm::*;

fn main() {
    let mut lam = Application(
        Box::new(Abstraction(Box::new(BoundVariable(0)))),
        Box::new(Abstraction(Box::new(FreeVariable(String::from("a"))))),
    );
    println!("{}", lam);
    lam.beta_reduce();
    println!("{}", lam);
}
