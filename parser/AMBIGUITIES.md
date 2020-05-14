# Ambiguities

Modern ES262 spec has a few ambiguities parser has to deal with.

We're taking try and fail approach meaning that we will try
to parse expression as one variant and if we fail then try another.

This means that we're wasting all the parsed state
from the first try and start from the scratch.

Even though some cases can be successfully parsed as a superset
for the sake of other stuff that needs to be done
we just follow try-and-fail approach for them as well.

All the ambiguities parser has to deal with:

#### 1) Arrow function's parameters or sequence expression?

**Code:**
```js
const test = (foo, bar)
```

Until we see an arrow token we're not sure.

**Possible superset:** _SequenceExpression_

**Superset obstacles:**
- Need for _ObjectLiteral_ and _ObjectBinding_ superset.
- Need for _ArrayLiteral_ and _ArrayBinding_ superset.

**Example:**
```js
... ({a: b = c, f}, d = 0, ...e)
```

This could be parsed as:
```json
SequenceSuperset {
  List [
    ObjectSuperset {
      Properties [
        ObjectProperty {
          Id "a"
          Binder IdentifierBinder { "b" }
          Initializer Id { "c" }
        }
        ObjectProperty {
          Id "f"
        }
      ]
    }
    AssignmentExpression {
      Left Id { "d" }
      Right NumberLiteral { "0" }
    }
    SpreadExpression {
      Expression Id { "e" }
    }
  ]
}
```

Which can be boiled down to _SequenceExpression_ or _FunctionParameters_.
We can further extend _SequenceSuperset_ to have a `NotASequece` flag that could be determined during parsing phase
because some syntax of _FunctionParameters" is illegal for _SequenceExpression_.
