In general scheme follows:

```
Parser -> ?Optimizer -> Transpiler -> Generator -> Bundler
```

### Parser
Parser produces AST from source code.
Additionally, it allocates all symbols.
_`Symbol`_ - is a mediator to bind _`Identifier`_ to it's _`Ref`_ .
During parse stage we don't try to resolve 
identifiers because of all the ambiguities present when parsing javascript.

Instead, we reference all symbols after parse was complete.
This is done w/o traversing over AST as all symbols
stored separately in _`Module`_'s _`SymbolsScope`_.

### ?Optimizer
Question stands for an optional step.
Optimizations available:
**TODO**

### Transpiler
Transpiler will modify AST to conform target ES version.

### Generator
Produces code from AST and generates source map.

### Bundler
Bundles all modules into chunks.
