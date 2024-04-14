# Morpheus
![coverage](https://img.shields.io/badge/coverage-no_clue-darkred) ![bugs](https://img.shields.io/badge/bugs-a_lot-darkgreen)

Simple programming language with janky constraint based layout system

### Features
#### Strings, Booleans and Integers
```
x = "this is a string";
y = true;
z = 100;
a = z;
```

#### Arithmetic
```
x = 5 + 3;
y = -1 * (1 + 1);
z = x + y;
```

#### Comparisons
```
x = (5 > 2);
z = ("hello" == "hello")
y = (true == false);
```

#### String concatenation
```
x = "hello" ++ " " ++ "world"
```

#### Loops
increasing
```
acc = 0;

for i in (0, 5, 1) {
 acc = acc + i
}
```

decreasing
```
acc = 0;

for i in (0, -10, -1) {
 acc = acc + i
}
```

#### Functions
implicitly returns last expression
```
function add(x, y) {
    x + y
}

add(2, 2);
```


#### Conditionals
```
number = 30;

if (number < 0) {
    print("very small")
} elif (number < 10)  {
    print("small");
} elif (number < 100 ) {
    print("medium");
} else {
    print("big");
}
// outputs medium
```

#### Lists
```
list = [1,2,3];
print(list.len); // 3
list.del(1);
print(list); // [1,3]
list.add(10);
print(list); // [1,3,10]
list.get(2);
```

### Layout

#### Boxes
```
a = Box("box a");
b = Box("box b");
```

#### Groups
```
a = Box("box a");
b = Box("box b");

g = Group([a, b] : []);
```

#### Constraints
```
a = Box("box a");
b = Box("box b");

g = Group([a, b] : [
    *a is below *b,
    *b is right of *a
]);
```


#### Output Layout as HTML
```
a = Box("box a");
b = Box("box b");

g = Group([a, b] : [
    *a is below *b,
    *b is right of *a
]);

g.htmlify("output_file_name"); // save as output_file_name.html
```


