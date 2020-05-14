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

**Possible superset:** _`SequenceSuperset`_

**Superset dependencies:**
- Need for _`ObjectLiteral`_ and _`ObjectBinding`_ superset: _`ObjectSuperset`_.
- Need for _`ArrayLiteral`_ and _`ArrayBinding`_ superset: _`ArraySuperset`_.

**Example:**
```js
... ({a: b = c, f}, d = 0, ...e)
```

This could be parsed as:
```
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

Which can be boiled down to _`SequenceExpression`_ or _`FunctionParameters`_.
We can further extend _`SequenceSuperset`_ to have a `NotASequece` flag that could be determined during parsing phase
because some syntax of _`FunctionParameters`_ is illegal for _`SequenceExpression`_.
