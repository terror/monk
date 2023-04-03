**monk** is a tree-walk interpreter for a small language — following the text
*Writing an Interpreter in Go* by Thorsten Ball.

```
let add = fn(a, b) { return a + b; };

let add = fn(a, b) { a + b; };

add(1, 2);

let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      1
    } else {
     fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};
```
