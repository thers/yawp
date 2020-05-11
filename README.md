# Yet Another Web Packer

Supports ES2020, Flow, JSX

In order to simplify parser we always assume strict mode.

---
### Parser problems left to solve
- [ ] yield sometimes isn't a keyword and allowed to be used as identifier

### Optimizer progress
- [x] identifiers mangling
- [x] const/let transformation
- [x] destructuring assignment transformation
- [ ] class transformation
- [ ] arrow fn transformation
- [ ] async function transformation
- [ ] generator function transformation
- [ ] object literal extensions transformation
- [ ] unused imports removal
- [ ] dead code elimination

---
### Generator progress
- [ ] variables
- [ ] source map generation

---
##### Parser progress left
- [ ] modern decorators
- [ ] flow declare type/interface/var/function/class
- [ ] flow declare module