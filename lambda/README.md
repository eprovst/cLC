# lambda

This page gives an overview of the library on which the cLC is build.

All functions in this library work upon `lambda.Term`s which is an interface covering a few other datatypes.
More detail is provided later on. However for most applications it's easier to use the serialization format instead
of the internal representation.

## Serialization
The library has a serialization format to provide an easier way to start using the library (the internal representation is rather strict).

It's basically De Bruijn index notation using `l` instead of `λ` and with free variables prefixed by `'` and separated by a space or brackets. For example for the lambda expression (λ*x*.(λ*y*.*x*) *x*) *v* *w* you use: `(l(l2)1)'v 'w`.

There are two functions defined for this format:

- `Deserialize(string) (Term, error)`

Turns the De Bruijn index notation into the internal representation.


- `MustDeserialize(string) Term`

Same as above but panics on error.


## Functions defined on `Term`s

- `AlphaEquivalent(Term) bool`

Checks if self and the parameter are syntactically equal.


- `String() string`

Returns the object as a string.


- `Serialize() string`

Returns the object in the serialized format.


- `Reduce() (Term, error)`

Reduces the object using the fastest method available (currently `AorReduce()`).


- `AorReduce() (Term, error)`

Reduces the object using applicative order reduction.


- `NorReduce() (Term, error)`

Reduces the object using normal order reduction.


- `WHNF() Abst`

Transforms the object to a weak head normal form (usually by applying the inverse of η-reduction).


- `EtaReduce() Term`

If there is eta-reduction is possible returns the object after this reduction, else return the object itself.


- (Only defined on `Abst`) `BetaReduce(sub Term) Term`

Substitutes the locally bound variable of the Abst by `sub`.


## Variables

- `var MaxReductions = 10000`

Defines the maximum amount of reductions that are tried by the library before returning an error (useful when, for example, you try to expand the Y combinator).

If set to a negative value it will keep on expanding indefinitely.


## Data types
If you need simple input of lambda expressions, look at the serialization format.

`lambda` contains three datatypes and one interface to represent lambda terms:

- `Term` (interface) is implemented by all three data types.
- `Appl` (`[2]Term`) represents an application.
- `Abst` (`[1]Term`) represents a lambda abstraction.
- `Var` (`uint`) is the De Bruijn index of a variable minus one.
- `Free` (`string`) is the name of a free variable.

Thus as an example, the lambda term (λ*x*.(λ*y*.*x*) *x*) *v* *w* becomes:

```
Appl{
    Appl{
        Abst{
            Appl{
                Abst{
                    Var(1),
                },
                Var(0),
            },
        },
        Free("v"),
    },
    Free("w"),
}
```

As you can see this entire structure is a `Term`.